package main

// This is a transport adapter for the existing Idriç executable, NOT another
// conversation engine. The journal, evidence revisions, consent, questions and
// event idempotence remain exclusively owned by yente-conversation.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const maxCoreOutput = 1 << 20

type CoreBridge struct {
	Binary, StateRoot, Materials string
	mu                           sync.Mutex // serialize script/show across this gateway instance
	run                          func(context.Context, ...string) (string, error)
}

func newCoreBridge(binary, root, materials string) (*CoreBridge, error) {
	for _, path := range []string{binary, root, materials} {
		if !filepath.IsAbs(path) || filepath.Clean(path) != path {
			return nil, errors.New("Idriç executable, journal root and materials must use canonical absolute paths")
		}
	}
	if info, err := os.Stat(binary); err != nil || info.IsDir() || info.Mode()&0111 == 0 {
		return nil, errors.New("Idriç executable is missing or is not executable")
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return nil, errors.New("Idriç journal root must exist and be a trusted directory")
	}
	if _, err := os.Stat(filepath.Join(materials, "challenges.tsv")); err != nil {
		return nil, errors.New("Idriç challenge materials are missing")
	}
	bridge := &CoreBridge{Binary: binary, StateRoot: root, Materials: materials}
	bridge.run = func(ctx context.Context, args ...string) (string, error) {
		cmd := exec.CommandContext(ctx, binary, args...)
		// The executable reads fixture-only paths. No shell interpolation.
		out, err := cmd.CombinedOutput()
		if len(out) > maxCoreOutput {
			return "", errors.New("Idriç output exceeds size limit")
		}
		if err != nil {
			return "", fmt.Errorf("Idriç process failed: %w", err)
		}
		return string(out), nil
	}
	return bridge, nil
}

// write a private, ephemeral command file for the real `script` interface.
func (b *CoreBridge) script(ctx context.Context, command string) (string, error) {
	file, err := os.CreateTemp("", "yente-script-*.tsv")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())
	if _, err := io.WriteString(file, command+"\n"); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return b.run(ctx, "script", b.StateRoot, b.Materials, file.Name())
}

type coreView struct {
	discussion, retention            bool
	accounts, interpretations        int
	latestID, latestKind             string
	activeQuestionID, activeQuestion string
}

// A length-prefixed row is a sequence of Unicode scalar counts followed by
// exactly that many characters. Do NOT split on tabs or decode private text as
// commands. In this bridge the resulting account passages are never returned.
func parseFramedRow(line string) ([]string, error) {
	if line == "" {
		return nil, errors.New("empty row")
	}
	runeLine := []rune(line)
	fields := make([]string, 0, 12)
	for at := 0; at < len(runeLine); {
		start := at
		for at < len(runeLine) && runeLine[at] >= '0' && runeLine[at] <= '9' {
			at++
		}
		if at == start || at >= len(runeLine) || runeLine[at] != ':' || at-start > 7 {
			return nil, errors.New("bad length framing")
		}
		size, err := strconv.Atoi(string(runeLine[start:at]))
		if err != nil || size < 0 || size > maxCoreOutput {
			return nil, errors.New("bad field length")
		}
		at++
		if size > len(runeLine)-at {
			return nil, errors.New("truncated framed value")
		}
		fields = append(fields, string(runeLine[at:at+size]))
		at += size
		if len(fields) > 24 {
			return nil, errors.New("too many fields")
		}
	}
	return fields, nil
}

// Only the real Idriç `show` operation supplies this projection. Nothing is
// constructed from client history or from another route's journal entries.
func (b *CoreBridge) view(ctx context.Context, route Route) (coreView, error) {
	output, err := b.script(ctx, "show\t"+route.Session)
	if err != nil {
		return coreView{}, err
	}
	if len(output) > maxCoreOutput || !strings.HasPrefix(output, "CONTEXT\t"+route.Session+"\n") {
		return coreView{}, errors.New("Idriç context operation was not accepted")
	}
	var view coreView
	hasSession, hasMode := false, false
	for _, line := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(output, "CONTEXT\t"+route.Session+"\n"), "\n"), "\n") {
		if line == "" {
			continue
		}
		fields, err := parseFramedRow(line)
		if err != nil || len(fields) == 0 {
			return coreView{}, errors.New("malformed Idriç context row")
		}
		switch fields[0] {
		case "session":
			if hasSession || len(fields) != 6 || fields[1] != route.Session || fields[2] != route.Role ||
				!containsActor(fields[4], route.Actor) {
				return coreView{}, errors.New("route disagrees with Idriç session identity")
			}
			hasSession = true
		case "mode":
			if hasMode || len(fields) != 6 || (fields[1] != "True" && fields[1] != "False") || (fields[2] != "True" && fields[2] != "False") {
				return coreView{}, errors.New("malformed Idriç permissions")
			}
			view.discussion = fields[1] == "True"
			view.retention = fields[2] == "True"
			hasMode = true
		case "evidence":
			if len(fields) != 9 {
				return coreView{}, errors.New("malformed Idriç evidence")
			}
			if fields[7] == "current" {
				view.accounts++
				view.latestID, view.latestKind = fields[1], fields[3]
			}
		case "tentative":
			view.interpretations++
		case "previous-question":
			if len(fields) != 5 {
				return coreView{}, errors.New("malformed question")
			}
			if fields[2] != "withdrawn" {
				view.activeQuestion, view.activeQuestionID = fields[3], fields[4]
			}
		}
	}
	if !hasSession || !hasMode {
		return coreView{}, errors.New("missing Idriç session/permissions")
	}
	return view, nil
}
func containsActor(csv, actor string) bool {
	for _, value := range strings.Split(csv, ",") {
		if value == actor {
			return true
		}
	}
	return false
}

// Explicit verbs keep a question, wish and hypothetical distinct from a report.
// A free-form chat utterance is never silently turned into factual evidence.
func coreCommand(route Route, eventID, input string) (string, string, error) {
	if !utf8.ValidString(input) || len(input) > maxMessageBytes ||
		strings.ContainsAny(input, "\t\r\n\x00") {
		return "", "", errors.New("command must be one clean UTF-8 line")
	}
	if input == "/context" {
		return "context", "", nil
	}
	if input == "/withdraw-consent" {
		return "consent", "event\t" + eventID + "\tconsent\t" + route.Session + "\tno\tno", nil
	}
	if input == "/ask" {
		return "ask", "event\t" + eventID + "\task-next\t" + route.Session, nil
	}
	word, tail, ok := strings.Cut(input, " ")
	if !ok || tail == "" {
		return "", "", errors.New("use /report TEXT, /wish TEXT, /ask, /reply QUESTION KIND TEXT, /correct EVIDENCE KIND TEXT, /context or /withdraw-consent")
	}
	kind := ""
	switch word {
	case "/report":
		kind = "report"
	case "/wish":
		kind = "wish"
	case "/hearsay":
		kind = "hearsay"
	case "/hypothetical":
		kind = "hypothetical"
	case "/agreement-report":
		kind = "agreement-report"
	}
	if kind != "" {
		if strings.TrimSpace(tail) == "" {
			return "", "", errors.New("account text cannot be empty")
		}
		return "account", "event\t" + eventID + "\tsay\t" + route.Session + "\t" + route.Actor + "\t" + kind + "\t\tfictional-time\t" + tail, nil
	}
	if word == "/correct" || word == "/reply" {
		target, rest, ok := strings.Cut(tail, " ")
		if !ok || !validID(target) || strings.TrimSpace(rest) == "" {
			return "", "", errors.New("correction/reply requires an identifier, kind, and wording")
		}
		sourceKind, wording, ok := strings.Cut(rest, " ")
		if !ok || strings.TrimSpace(wording) == "" {
			return "", "", errors.New("correction/reply requires explicit kind and wording")
		}
		switch sourceKind {
		case "report", "hearsay", "wish", "hypothetical", "agreement-report":
		default:
			return "", "", errors.New("invalid account kind")
		}
		operation := "correct"
		if word == "/reply" {
			operation = "reply"
		}
		return operation, "event\t" + eventID + "\t" + operation + "\t" + route.Session + "\t" + route.Actor + "\t" + target + "\t" + sourceKind + "\t\tfictional-time\t" + wording, nil
	}
	return "", "", errors.New("unrecognized explicit fictional command")
}

func coreEventID(route Route, channel, userID string) string {
	sum := sha256.Sum256([]byte(route.Name + "\x00" + channel + "\x00" + userID))
	return "transport-" + hex.EncodeToString(sum[:16])
}

func coreStatus(output string) int {
	switch {
	case strings.Contains(output, "\tstale:"):
		return http.StatusConflict
	case strings.Contains(output, "\tforbidden:"):
		return http.StatusForbidden
	case strings.Contains(output, "\tunknown:"), strings.Contains(output, "\tinvalid:"):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusBadGateway
	}
}

func (g *Gateway) respondCore(w http.ResponseWriter, r *http.Request, route Route, input string) {
	if !validMessage(input) {
		http.Error(w, "invalid fictional command", http.StatusUnprocessableEntity)
		return
	}
	id := r.Header.Get("X-Yente-Event-ID")
	if id != "" && !validID(id) {
		http.Error(w, "invalid event identity", http.StatusUnprocessableEntity)
		return
	}
	if input != "/context" && id == "" {
		http.Error(w, "state-changing commands require X-Yente-Event-ID", http.StatusUnprocessableEntity)
		return
	}
	operation, command, err := coreCommand(route, coreEventID(route, "cli", id), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// A single local prototype process owns its trusted journal root.
	g.Core.mu.Lock()
	defer g.Core.mu.Unlock()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	var emitted string
	if operation != "context" {
		emitted, err = g.Core.script(ctx, command)
		if err != nil {
			http.Error(w, "Idriç core unavailable", http.StatusBadGateway)
			return
		}
		if !strings.HasPrefix(emitted, "OK\t") && !strings.HasPrefix(emitted, "YENTA\t") {
			http.Error(w, "Idriç rejected fictional command", coreStatus(emitted))
			return
		}
		if strings.Count(strings.TrimSuffix(emitted, "\n"), "\n") != 0 {
			http.Error(w, "unexpected Idriç output", http.StatusBadGateway)
			return
		}
	}
	view, err := g.Core.view(ctx, route)
	if err != nil {
		http.Error(w, "Idriç scoped context unavailable", http.StatusBadGateway)
		return
	}
	reply := "Yente deterministic fictional mock (" + route.Name + "): "
	if !view.discussion || !view.retention {
		reply += "discussion or retention permission withdrawn; no account history or question is disclosed. The audit journal is not erased."
	} else {
		switch operation {
		case "ask":
			if strings.HasPrefix(emitted, "YENTA\t"+route.Session+"\t") && view.activeQuestionID != "" {
				reply += "question " + view.activeQuestionID + ": " + view.activeQuestion
			} else {
				reply += "a prior question was not replayed; use /context for the current session."
			}
		case "account", "reply", "correct":
			reply += fmt.Sprintf("%s recorded by Idriç; %d current attributed account(s)", operation, view.accounts)
			if view.latestID != "" {
				reply += "; latest " + view.latestID + " (" + view.latestKind + ")"
			}
			reply += ". No conversational inference."
		case "context":
			reply += fmt.Sprintf("%d current attributed account(s), %d active interpretation(s)", view.accounts, view.interpretations)
			if view.latestID != "" {
				reply += "; latest " + view.latestID + " (" + view.latestKind + ")"
			}
			if view.activeQuestionID != "" {
				reply += "; latest question " + view.activeQuestionID + ": " + view.activeQuestion
			}
			reply += ". No cross-recipient data."
		default:
			reply += "Idriç accepted the fictional command."
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transportReply{Route: route.Name, Channel: "cli", Session: route.Session, Reply: reply, Replayed: strings.HasPrefix(emitted, "OK\treplayed\t")})
}
