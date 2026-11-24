package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

const sampleCoreRoutes = "partner-a\tlio-private\tlio\tpartner\npartner-b\tmira-private\tmira\tpartner\nshared\tcouple-shared\tlio\tshared\n"

func framed(fields ...string) string {
	var builder strings.Builder
	for _, f := range fields {
		builder.WriteString(fmt.Sprintf("%d:%s", len([]rune(f)), f))
	}
	return builder.String() + "\n"
}

// A test double isolates the HTTP adapter contract. Only TestRealCoreE2E
// (when opted into a real executable) establishes Idriç journal execution.
func TestCoreBridgeHTTPAdapterWithScriptDouble(t *testing.T) {
	routes, err := loadRoutes(strings.NewReader(sampleCoreRoutes))
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var seen []string
	var scope = map[string]struct {
		count   int
		consent bool
	}{
		"lio-private":  {count: 0, consent: true},
		"mira-private": {count: 0, consent: true},
	}
	bridge := &CoreBridge{StateRoot: t.TempDir(), Materials: t.TempDir()}
	bridge.run = func(_ context.Context, args ...string) (string, error) {
		if len(args) != 4 || args[0] != "script" {
			return "", fmt.Errorf("unexpected exec call %v", args)
		}
		data, err := os.ReadFile(args[3])
		if err != nil {
			return "", err
		}
		line := strings.TrimSuffix(string(data), "\n")
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, line)
		fields := strings.Split(line, "\t")
		if fields[0] == "show" {
			session := fields[1]
			v, ok := scope[session]
			if !ok {
				return "ERR\tshow\tunknown: session\n", nil
			}
			actor := "lio"
			if session == "mira-private" {
				actor = "mira"
			}
			role := "partner"
			mode := "True"
			if !v.consent {
				mode = "False"
			}
			var out strings.Builder
			out.WriteString("CONTEXT\t" + session + "\n")
			out.WriteString(framed("session", session, role, actor, actor, "fictional"))
			out.WriteString(framed("mode", mode, mode, "", "", ""))
			for i := 1; i <= v.count; i++ {
				out.WriteString(framed("evidence", fmt.Sprintf("%s-e%d", session, i), actor, "report", "", "fictional-time", "1", "current", "fictional account"))
			}
			return out.String(), nil
		}
		if fields[0] != "event" {
			return "ERR\toperation\tinvalid: command\n", nil
		}
		operation := fields[2]
		session := ""
		switch operation {
		case "say":
			session = fields[3]
		case "correct":
			session = fields[3]
		case "consent":
			session = fields[3]
		}
		v, ok := scope[session]
		if !ok {
			return "ERR\toperation\tunknown: session\n", nil
		}
		if operation == "consent" {
			v.consent = false
		} else if !v.consent {
			return "ERR\t" + operation + "\tforbidden: retention consent required\n", nil
		} else if operation == "say" {
			v.count++
		}
		scope[session] = v
		return "OK\t" + operation + "\t" + session + "\n", nil
	}
	gateway := newGateway(routes, "fictional-token")
	gateway.Core = bridge
	server := httptest.NewUnstartedServer(gateway)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()
	client := server.Client()
	call := func(route, id, body string) (int, transportReply, string) {
		t.Helper()
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/cli/"+route, strings.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer fictional-token")
		if id != "" {
			req.Header.Set("X-Yente-Event-ID", id)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(res.Body)
		res.Body.Close()
		var reply transportReply
		if res.StatusCode == 200 {
			if err := json.Unmarshal(data, &reply); err != nil {
				t.Fatal(err)
			}
		}
		return res.StatusCode, reply, string(data)
	}
	if s, _, _ := call("partner-a", "1", "/report I prefer Sunday visits"); s != 200 {
		t.Fatal(s)
	}
	s, first, _ := call("partner-a", "", "/context")
	if s != 200 || !strings.Contains(first.Reply, "1 current attributed account") || strings.Contains(first.Reply, "fictional account") {
		t.Fatal(first)
	}
	s, other, _ := call("partner-b", "", "/context")
	if s != 200 || !strings.Contains(other.Reply, "0 current attributed account") || strings.Contains(other.Reply, "lio") {
		t.Fatal(other)
	}
	if s, _, _ := call("partner-a", "2", "/withdraw-consent"); s != 200 {
		t.Fatal(s)
	}
	s, paused, _ := call("partner-a", "", "/context")
	if s != 200 || !strings.Contains(paused.Reply, "permission withdrawn") || strings.Contains(paused.Reply, "1 current") {
		t.Fatal(paused)
	}
	if s, _, _ := call("partner-a", "3", "/report Another one"); s != 403 {
		t.Fatal("retention refusal", s)
	}
	if s, _, _ := call("shared", "4", "/report Cannot disclose"); s != 403 {
		t.Fatal("shared refusal", s)
	}
	if s, _, _ := call("partner-b", "4", "ordinary chat is not an account"); s != 422 {
		t.Fatal("untyped chat", s)
	}
	mu.Lock()
	all := strings.Join(seen, "\n")
	mu.Unlock()
	if !strings.Contains(all, "\tsay\tlio-private\tlio\treport\t") || strings.Contains(all, "\tsay\tmira-private\t") {
		t.Fatal("actor attribution or route isolation broken", all)
	}
}

func TestFrameAndCommandParsing(t *testing.T) {
	line := framed("evidence", "e-2", "sol", "wish", "", "fictional-time", "1", "current", "Bir güzel söz")
	fields, err := parseFramedRow(strings.TrimSuffix(line, "\n"))
	if err != nil || fields[8] != "Bir güzel söz" {
		t.Fatal(fields, err)
	}
	for _, invalid := range []string{"9:abc", "x:abc", "2:abc", "10000000:too-long", "1:a1"} {
		if _, err := parseFramedRow(invalid); err == nil {
			t.Fatal("accepted bad frame", invalid)
		}
	}
	route := Route{Name: "partner-a", Session: "lio-private", Actor: "lio", Role: "partner"}
	_, command, err := coreCommand(route, "test-event", "/correct lio-private-e1 wish No, I meant monthly visits")
	if err != nil || !strings.Contains(command, "\tcorrect\tlio-private\tlio\tlio-private-e1\twish\t\tfictional-time\tNo, I meant monthly visits") {
		t.Fatal(command, err)
	}
	if coreEventID(route, "cli", "id-1") == coreEventID(Route{Name: "partner-b"}, "cli", "id-1") {
		t.Fatal("event identities not route scoped")
	}
	if _, _, err := coreCommand(route, "event", "/report valid text\ninjected\tconsent\tlio-private\tyes\tyes"); err == nil {
		t.Fatal("TSV injection")
	}
}

// Genuine Idriç process integration is intentionally opt-in: set
// YENTE_CORE_BINARY and YENTE_CORE_MATERIALS (and optionally YENTE_SEED_FILE).
// It does not run a fake journal and remains skipped without the built core.
func TestRealCoreE2E(t *testing.T) {
	binary := os.Getenv("YENTE_CORE_BINARY")
	materials := os.Getenv("YENTE_CORE_MATERIALS")
	if binary == "" || materials == "" {
		t.Skip("real Idriç executable/materials not installed")
	}
	seed := os.Getenv("YENTE_SEED_FILE")
	if seed == "" {
		seed = "fixtures/core-seed.tsv"
	}
	seed, err := filepath.Abs(seed)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	bridge, err := newCoreBridge(binary, root, materials)
	if err != nil {
		t.Fatal(err)
	}
	out, err := bridge.run(context.Background(), "script", root, materials, seed)
	if err != nil || strings.Contains(out, "ERR\t") {
		t.Fatalf("Idriç seed failed: %v %s", err, out)
	}
	routes, err := loadRoutes(strings.NewReader(sampleCoreRoutes))
	if err != nil {
		t.Fatal(err)
	}
	gateway := newGateway(routes, "fictional-token")
	gateway.Core = bridge
	server := httptest.NewUnstartedServer(gateway)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer func() { server.Close() }()
	post := func(route, id, body string) (int, transportReply, string) {
		t.Helper()
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/cli/"+route, strings.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer fictional-token")
		if id != "" {
			req.Header.Set("X-Yente-Event-ID", id)
		}
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(res.Body)
		res.Body.Close()
		var reply transportReply
		_ = json.Unmarshal(data, &reply)
		return res.StatusCode, reply, string(data)
	}
	check := func(route, id, body, needle string) {
		t.Helper()
		status, reply, raw := post(route, id, body)
		if status != 200 || !strings.Contains(reply.Reply, needle) {
			t.Fatalf("%s %s: %d %s", route, body, status, raw)
		}
		t.Logf("%s %s => %s", route, body, reply.Reply)
	}
	check("partner-a", "ask-one", "/ask", "How often would you like to visit your family?")
	check("partner-a", "answer-one", "/reply lio-private-p1 wish I would like to visit weekly", "1 current attributed account")
	check("partner-a", "event-one", "/wish I want to visit weekly", "2 current attributed account")
	check("partner-b", "", "/context", "0 current attributed account")
	check("partner-a", "event-two", "/correct lio-private-e2 wish I meant once a month", "2 current attributed account")
	check("partner-a", "", "/context", "latest lio-private-e3")
	if status, _, _ := post("partner-b", "cross-correct", "/correct lio-private-e3 wish I did not say that"); status != 403 {
		t.Fatalf("another route corrected someone else's evidence: %d", status)
	}
	// A repeated ask must not disclose old Yente wording after its context changed.
	check("partner-a", "ask-one", "/ask", "prior question was not replayed")
	// A second Gateway process uses the same durable core journal.
	server.Close()
	gateway2 := newGateway(routes, "fictional-token")
	bridge2, err := newCoreBridge(binary, root, materials)
	if err != nil {
		t.Fatal(err)
	}
	gateway2.Core = bridge2
	server = httptest.NewUnstartedServer(gateway2)
	server.EnableHTTP2 = true
	server.StartTLS()
	check("partner-a", "", "/context", "latest lio-private-e3")
	// Idriç's persisted event identity must survive this gateway restart.
	check("partner-a", "event-one", "/wish I want to visit weekly", "2 current attributed account")
	if status, _, _ := post("partner-a", "event-one", "/wish A different assertion"); status != 409 {
		t.Fatalf("event ID collision accepted after gateway restart: %d", status)
	}
	check("partner-a", "off", "/withdraw-consent", "permission withdrawn")
	check("partner-a", "", "/context", "permission withdrawn")
	if status, _, _ := post("partner-a", "after-off", "/wish something else"); status != 403 {
		t.Fatalf("post-withdrawal accepted: %d", status)
	}
	check("partner-b", "", "/context", "0 current attributed account")
	if status, _, _ := post("shared", "disclose", "/context"); status != 403 {
		t.Fatalf("shared accepted: %d", status)
	}
	journal, err := os.ReadFile(filepath.Join(root, "journal.tsv"))
	if err != nil || !strings.Contains(string(journal), "transport-") {
		t.Fatal("real Idriç did not persist gateway events", err)
	}
}

func TestRealCoreAcrossFixtureTransports(t *testing.T) {
	binary, materials := os.Getenv("YENTE_CORE_BINARY"), os.Getenv("YENTE_CORE_MATERIALS")
	if binary == "" || materials == "" {
		t.Skip("real Idriç executable/materials not installed")
	}
	root := t.TempDir()
	seed, err := filepath.Abs("fixtures/core-seed.tsv")
	if err != nil {
		t.Fatal(err)
	}
	routes, err := loadRoutes(strings.NewReader(sampleCoreRoutes))
	if err != nil {
		t.Fatal(err)
	}
	const firstToken = "fictional-partner-a-token"
	const secondToken = "fictional-partner-b-token"
	tokens, err := loadRouteTokens(strings.NewReader("partner-a\t"+firstToken+"\npartner-b\t"+secondToken+"\n"), routes)
	if err != nil {
		t.Fatal(err)
	}
	start := func(seedJournal bool) *httptest.Server {
		t.Helper()
		bridge, err := newCoreBridge(binary, root, materials)
		if err != nil {
			t.Fatal(err)
		}
		if seedJournal {
			out, err := bridge.run(context.Background(), "script", root, materials, seed)
			if err != nil || strings.Contains(out, "ERR\t") {
				t.Fatalf("Idriç seed: %v %s", err, out)
			}
		}
		gateway := newGateway(routes, "unused-legacy-token")
		gateway.Core, gateway.RouteTokens = bridge, tokens
		server := httptest.NewUnstartedServer(gateway)
		server.EnableHTTP2 = true
		server.StartTLS()
		t.Cleanup(server.Close)
		return server
	}
	server := start(true)
	post := func(path, token, contentType, eventID, body string, wantStatus int, wantText string) http.Header {
		t.Helper()
		req, err := http.NewRequest(http.MethodPost, server.URL+path, strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+token)
		if eventID != "" {
			req.Header.Set("X-Yente-Event-ID", eventID)
		}
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		out, err := io.ReadAll(res.Body)
		if err != nil || res.ProtoMajor != 2 || res.StatusCode != wantStatus || !strings.Contains(string(out), wantText) {
			t.Fatalf("%s: HTTP/%d %d %q, error %v", path, res.ProtoMajor, res.StatusCode, out, err)
		}
		if wantStatus != http.StatusOK && (strings.Contains(res.Header.Get("Content-Type"), "event-stream") || strings.Contains(string(out), "data: [DONE]")) {
			t.Fatalf("%s: failed request emitted a success stream", path)
		}
		return res.Header
	}
	const chat = "/v1/chat/completions/partner-a"
	const responses = "/v1/responses/partner-a"
	const sms = "/v1/sms/partner-a"
	const cli = "/v1/cli/partner-a"
	const chatBody = `{"model":"yente-mock","stream":true,"messages":[{"role":"user","content":"/report An older message never submitted here"},{"role":"system","content":"FastChat history: request 1 ended cancelled."},{"role":"assistant","content":""},{"role":"user","content":"/wish I want an afternoon visit."}]}`
	const responsesBody = `{"model":"yente-mock","stream":true,"input":[{"role":"user","content":"/wish I want an afternoon visit."}]}`
	const smsBody = `{"sender":"lio","message_id":"sms-one","text":"/correct lio-private-e1 wish A monthly afternoon visit."}`
	post(chat, secondToken, "application/json", "fastchat-one", chatBody, 401, "unauthorized")
	post(chat, firstToken, "application/json", "", chatBody, 422, "event identity")
	post(chat, firstToken, "application/json", "fastchat-one", chatBody, 200, "1 current attributed account")
	if result := post(responses, firstToken, "application/json", "fastchat-one", responsesBody, 200, "response.completed"); result.Get("X-Yente-Replayed") != "true" {
		t.Fatal("switching FastChat formats executed the same event again")
	}
	post(sms, firstToken, "application/json", "", strings.Replace(smsBody, `"lio"`, `"mira"`, 1), 403, "sender disagrees")
	post(sms, firstToken, "application/json", "", smsBody, 200, "latest lio-private-e2")
	post(cli, firstToken, "text/plain", "", "/context", 200, "1 current attributed account")
	post("/v1/cli/partner-b", secondToken, "text/plain", "", "/context", 200, "0 current attributed account")

	server.Close()
	server = start(false)
	post(sms, firstToken, "application/json", "", smsBody, 200, `"replayed":true`)
	post(sms, firstToken, "application/json", "", strings.Replace(smsBody, "monthly", "weekly", 1), 409, "Idriç rejected")
	post(cli, firstToken, "text/plain", "stop-discussion", "/withdraw-consent", 200, "permission withdrawn")
	post(chat, firstToken, "application/json", "fastchat-one", chatBody, 200, "permission withdrawn")
	post(chat, firstToken, "application/json", "fastchat-two", chatBody, 403, "Idriç rejected")

	journal, err := os.ReadFile(filepath.Join(root, "journal.tsv"))
	if err != nil || strings.Count(string(journal), "transport-event\ttransport-") != 3 || strings.Contains(string(journal), "older message never submitted") {
		t.Fatalf("expected one FastChat event, one SMS correction and one CLI withdrawal in Idriç journal: %v", err)
	}
}

func TestRealCoreAdmissionRecovery(t *testing.T) {
	binary, materials := os.Getenv("YENTE_CORE_BINARY"), os.Getenv("YENTE_CORE_MATERIALS")
	if binary == "" || materials == "" {
		t.Skip("real Idriç executable/materials not installed")
	}
	root := t.TempDir()
	fixtureDir := filepath.Join(filepath.Dir(materials), "fixtures")
	seed, err := os.ReadFile(filepath.Join(fixtureDir, "family.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	beforeCorrection, _, found := strings.Cut(string(seed), "event\tfamily-correction-lio\t")
	if !found {
		t.Fatal("family fixture correction missing")
	}
	script := filepath.Join(root, "trial.tsv")
	run := func(ok bool, needle string, args ...string) string {
		t.Helper()
		out, err := exec.Command(binary, args...).CombinedOutput()
		if (err == nil) != ok || !strings.Contains(string(out), needle) || (!ok && strings.Contains(string(out), "ANSWER\t")) {
			t.Fatalf("%s: %v, %q", args[0], err, out)
		}
		return string(out)
	}
	write := func(path, content string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write(script, beforeCorrection)
	run(true, "OK\tsay\tlio-private-e1", "script", root, materials, script)
	routes, err := loadRoutes(strings.NewReader(sampleCoreRoutes))
	if err != nil {
		t.Fatal(err)
	}
	bridge, err := newCoreBridge(binary, root, materials)
	if err != nil {
		t.Fatal(err)
	}
	const token = "fixture-partner-a-admission"
	gateway := newGateway(routes, "")
	gateway.Core = bridge
	gateway.RouteTokens = map[string]string{"partner-a": token}
	server := httptest.NewTLSServer(gateway)
	defer server.Close()
	post := func(id, text string) {
		t.Helper()
		payload, _ := json.Marshal(smsRequest{Sender: "lio", MessageID: id, Text: text})
		req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/sms/partner-a", strings.NewReader(string(payload)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			t.Fatalf("SMS %s returned %d %s", id, res.StatusCode, body)
		}
	}
	const say = "/correct lio-private-e1 wish I want to visit most weeks myself; Mira is welcome, and I do not expect her every week."
	post("source", say)
	id := coreEventID(routes["partner-a"], "sms", "source")
	const session = "lio-private"
	request := filepath.Join(root, "request.txt")
	plan := filepath.Join(fixtureDir, "family-response.tsv")
	run(false, "unknown: answer receipt", "recover-answer", root, materials, session, id)
	run(true, "OK\trequest", "request", root, materials, session, request)
	run(false, "unknown: conversation input event", "accept-event", root, materials, session, "foreign", request, plan)
	answer := run(true, "ANSWER\tlio-private\t"+id+"\tlio-private-p2\texpectation", "accept-event", root, materials, session, id, request, plan)
	if recovered := run(true, "ANSWER\t", "recover-answer", root, materials, session, id); recovered != answer {
		t.Fatalf("new process recovered a different answer: %q", recovered)
	}
	if recovered := run(true, "ANSWER\t", "accept-event", root, materials, session, id, request, plan); recovered != answer {
		t.Fatalf("retry after commit recovered a different answer: %q", recovered)
	}
	planBytes, err := os.ReadFile(plan)
	if err != nil {
		t.Fatal(err)
	}
	changedPlan := filepath.Join(root, "different.tsv")
	write(changedPlan, strings.Replace(string(planBytes), "deciding which visits", "deciding what time", 1))
	run(false, "stale: answer event reused with different plan", "accept-event", root, materials, session, id, request, changedPlan)
	run(false, "unknown: conversation input event", "recover-answer", root, materials, "mira-private", id)
	post("later", "/wish I may visit tomorrow.")
	later := coreEventID(routes["partner-a"], "sms", "later")
	run(false, "stale: model request changed", "accept-event", root, materials, session, later, request, plan)
	post("correct-again", "/correct lio-private-e2 wish I may visit every month.")
	run(false, "withheld", "recover-answer", root, materials, session, id)
	request = filepath.Join(root, "fresh-request.txt")
	write(changedPlan, strings.Replace(string(planBytes), "lio-private-e2", "lio-private-e4", 1))
	run(true, "OK\trequest", "request", root, materials, session, request)
	run(true, "ANSWER\t", "accept-event", root, materials, session, later, request, changedPlan)
	post("withdraw", "/withdraw-consent")
	run(false, "withheld", "recover-answer", root, materials, session, later)
	write(script, "consent\tlio-private\tyes\tyes\n")
	run(true, "OK\tconsent", "script", root, materials, script)
	run(false, "withheld", "recover-answer", root, materials, session, later)
	journal, err := os.ReadFile(filepath.Join(root, "journal.tsv"))
	if err != nil || strings.Count(string(journal), "\nadmit-answer\t") != 2 {
		t.Fatalf("expected two committed answers despite retries: %v", err)
	}
}
