// Package main provides synthetic transport fixtures, not a matchmaker or
// participant authentication. Opt-in fixture routes invoke the original Idriç journal.
package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	maxRequestBytes = 8192 // FastChat's bounded provider request
	maxMessageBytes = 4096
	mockModel       = "yente-mock"
)

type Route struct {
	Name    string
	Session string
	Actor   string
	Role    string
}

// A route chooses a synthetic session. It is NOT proof of the sender's identity.
func loadRoutes(reader io.Reader) (map[string]Route, error) {
	bytes, err := io.ReadAll(io.LimitReader(reader, 32769))
	if err != nil || len(bytes) > 32768 {
		return nil, errors.New("route file unreadable or too large")
	}
	routes := make(map[string]Route)
	for lineNumber, line := range strings.Split(string(bytes), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 4 {
			return nil, fmt.Errorf("route line %d: expected four tab-separated fields", lineNumber+1)
		}
		for _, field := range fields[:3] {
			if !validID(field) {
				return nil, fmt.Errorf("route line %d: invalid route/session/actor", lineNumber+1)
			}
		}
		role := fields[3]
		if role != "partner" && role != "parent" && role != "friend" && role != "shared" {
			return nil, fmt.Errorf("route line %d: invalid role", lineNumber+1)
		}
		if _, exists := routes[fields[0]]; exists {
			return nil, fmt.Errorf("route line %d: duplicate route", lineNumber+1)
		}
		routes[fields[0]] = Route{Name: fields[0], Session: fields[1], Actor: fields[2], Role: role}
	}
	if len(routes) == 0 {
		return nil, errors.New("route file is empty")
	}
	return routes, nil
}

func loadRouteTokens(reader io.Reader, routes map[string]Route) (map[string]string, error) {
	body, err := io.ReadAll(io.LimitReader(reader, 32769))
	if err != nil || len(body) > 32768 {
		return nil, errors.New("fixture credential file unreadable or too large")
	}
	tokens := make(map[string]string)
	used := make(map[string]bool)
	for i, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, token, ok := strings.Cut(line, "\t")
		if _, exists := routes[name]; !ok || !exists || !validID(token) || len(token) < 24 || tokens[name] != "" || used[token] {
			return nil, fmt.Errorf("fixture credential line %d: expected a known route and a unique bearer token of at least 24 characters", i+1)
		}
		tokens[name], used[token] = token, true
	}
	if len(tokens) == 0 {
		return nil, errors.New("fixture credential file is empty")
	}
	return tokens, nil
}

func validID(s string) bool {
	if len(s) < 1 || len(s) > 96 || s == "." || s == ".." {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

func validMessage(s string) bool {
	return len(s) > 0 && len(s) <= maxMessageBytes && utf8.ValidString(s) &&
		!strings.ContainsRune(s, 0)
}

type receiptKey struct {
	Route, Channel, ID string
}

type receipt struct {
	Digest [32]byte
	Reply  string
}

type Gateway struct {
	Routes      map[string]Route
	Token       string
	Core        *CoreBridge
	RouteTokens map[string]string // opt-in per-route fixture credentials
	mutex       sync.Mutex
	receipts    map[receiptKey]receipt // volatile fake receipts; NOT Yente's journal
}

func newGateway(routes map[string]Route, token string) *Gateway {
	return &Gateway{Routes: routes, Token: token, receipts: make(map[receiptKey]receipt)}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.URL.Path == "/healthz" && r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "yente-transport-mock\n")
		return
	}
	if g.RouteTokens == nil && g.Token != "" && !hasBearer(r, g.Token) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	prefixes := []struct{ prefix, channel string }{
		{"/v1/chat/completions/", "fastchat"},
		{"/v1/responses/", "responses"},
		{"/v1/sms/", "sms"},
		{"/v1/cli/", "cli"},
	}
	for _, candidate := range prefixes {
		if !strings.HasPrefix(r.URL.Path, candidate.prefix) {
			continue
		}
		routeName := strings.TrimPrefix(r.URL.Path, candidate.prefix)
		route, ok := g.Routes[routeName]
		if !ok || !validID(routeName) || strings.Contains(routeName, "/") || r.URL.RawQuery != "" {
			http.NotFound(w, r)
			return
		}
		if g.RouteTokens != nil && !hasBearer(r, g.RouteTokens[routeName]) {
			http.Error(w, "unauthorized for fixture route", http.StatusUnauthorized)
			return
		}
		if route.Role == "shared" {
			// Only the Idriç core may authorize a shared discussion/retelling.
			http.Error(w, "shared sessions require explicit core authorization", http.StatusForbidden)
			return
		}
		switch candidate.channel {
		case "fastchat", "responses":
			g.chat(w, r, route, candidate.channel)
		case "sms":
			g.sms(w, r, route)
		case "cli":
			g.cli(w, r, route)
		}
		return
	}
	http.NotFound(w, r)
}

func hasBearer(r *http.Request, token string) bool {
	supplied := r.Header.Values("Authorization")
	return token != "" && len(supplied) == 1 &&
		subtle.ConstantTimeCompare([]byte(supplied[0]), []byte("Bearer "+token)) == 1
}

func (g *Gateway) statefulFixtures() bool {
	return g.Core != nil && g.RouteTokens != nil
}

// The single dispatch point for all frontends. There is deliberately no model
// generation, journal mutation, attribution, inference or disclosure here.
func mockReply(route Route, channel, text string) string {
	_ = channel
	_ = text
	return "Yente mock (" + route.Name + "): transport received the message. No model inference or disclosure occurred."
}

func strictJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "application/json required", http.StatusUnsupportedMediaType)
		return false
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "request unreadable", http.StatusBadRequest)
		}
		return false
	}
	if !utf8.Valid(body) {
		http.Error(w, "invalid UTF-8 JSON request", http.StatusBadRequest)
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return false
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		http.Error(w, "invalid trailing JSON", http.StatusBadRequest)
		return false
	}
	return true
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
	Messages []message `json:"messages"`
}

type responsesRequest struct {
	Model  string    `json:"model"`
	Stream bool      `json:"stream"`
	Input  []message `json:"input"`
}

func (g *Gateway) chat(w http.ResponseWriter, r *http.Request, route Route, channel string) {
	var model string
	var stream bool
	var messages []message
	if channel == "fastchat" {
		var payload chatRequest
		if !strictJSON(w, r, &payload) {
			return
		}
		model, stream, messages = payload.Model, payload.Stream, payload.Messages
	} else {
		var payload responsesRequest
		if !strictJSON(w, r, &payload) {
			return
		}
		model, stream, messages = payload.Model, payload.Stream, payload.Input
	}
	stateful := g.statefulFixtures()
	if model != mockModel || !stream || len(messages) == 0 || (!stateful && len(messages) > 32) {
		http.Error(w, "unsupported mock model, stream or message history", http.StatusUnprocessableEntity)
		return
	}
	last := messages[len(messages)-1]
	if last.Role != "user" {
		http.Error(w, "last message must be user", http.StatusUnprocessableEntity)
		return
	}
	var answer string
	if stateful {
		// FastChat resends earlier messages, including local status records.
		// The journal owns history. Only the final submitted command goes to it.
		for _, m := range messages[:len(messages)-1] {
			if m.Role != "user" && m.Role != "assistant" && m.Role != "system" {
				http.Error(w, "invalid earlier message role", http.StatusUnprocessableEntity)
				return
			}
		}
		reply, failure := g.executeCore(r.Context(), route, "fastchat", requestEventID(r), last.Content)
		if failure != nil {
			http.Error(w, failure.message, failure.status)
			return
		}
		answer = reply.Reply
		w.Header().Set("X-Yente-Replayed", strconv.FormatBool(reply.Replayed))
	} else {
		for _, m := range messages {
			if (m.Role != "user" && m.Role != "assistant") || !validMessage(m.Content) {
				http.Error(w, "invalid message", http.StatusUnprocessableEntity)
				return
			}
		}
		answer = mockReply(route, channel, last.Content)
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusOK)
	if channel == "fastchat" {
		sse(w, map[string]any{"choices": []any{map[string]any{"index": 0, "delta": map[string]string{"content": answer}, "finish_reason": nil}}})
		sse(w, map[string]any{"choices": []any{map[string]any{"index": 0, "delta": map[string]string{}, "finish_reason": "stop"}}})
	} else {
		sse(w, map[string]any{"type": "response.output_text.delta", "delta": answer})
		sse(w, map[string]any{"type": "response.completed", "response": map[string]string{"status": "completed"}})
	}
	_, _ = io.WriteString(w, "data: [DONE]\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func sse(w io.Writer, payload any) {
	bytes, _ := json.Marshal(payload) // all inputs here are server-authored
	_, _ = fmt.Fprintf(w, "data: %s\n\n", bytes)
}

type smsRequest struct {
	Sender    string `json:"sender"` // an unverified synthetic fixture label
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

type transportReply struct {
	Route    string `json:"route"`
	Channel  string `json:"channel"`
	Session  string `json:"session"`
	Reply    string `json:"reply"`
	Replayed bool   `json:"replayed"`
}

func (g *Gateway) respondOnce(w http.ResponseWriter, route Route, channel, id, sender, text string) {
	if !validID(id) || !validMessage(text) {
		http.Error(w, "invalid event identity or message", http.StatusUnprocessableEntity)
		return
	}
	key := receiptKey{Route: route.Name, Channel: channel, ID: id}
	digest := sha256.Sum256([]byte(sender + "\x00" + text))
	g.mutex.Lock()
	previous, exists := g.receipts[key]
	if exists && previous.Digest != digest {
		g.mutex.Unlock()
		http.Error(w, "event ID reused for different input", http.StatusConflict)
		return
	}
	answer := mockReply(route, channel, text)
	if exists {
		answer = previous.Reply
	} else if len(g.receipts) >= 4096 {
		g.mutex.Unlock()
		http.Error(w, "mock receipts exhausted", http.StatusServiceUnavailable)
		return
	} else {
		g.receipts[key] = receipt{Digest: digest, Reply: answer}
	}
	g.mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transportReply{Route: route.Name, Channel: channel, Session: route.Session,
		Reply: answer, Replayed: exists})
}

func (g *Gateway) sms(w http.ResponseWriter, r *http.Request, route Route) {
	var payload smsRequest
	if !strictJSON(w, r, &payload) {
		return
	}
	if !validID(payload.Sender) {
		http.Error(w, "invalid synthetic sender label", http.StatusUnprocessableEntity)
		return
	}
	if g.statefulFixtures() {
		if payload.Sender != route.Actor {
			http.Error(w, "fixture sender disagrees with route", http.StatusForbidden)
			return
		}
		g.respondCoreJSON(w, r, route, "sms", payload.MessageID, payload.Text)
		return
	}
	// sender is deliberately NOT used as authentication or actor attribution.
	g.respondOnce(w, route, "sms", payload.MessageID, payload.Sender, payload.Text)
}

func (g *Gateway) cli(w http.ResponseWriter, r *http.Request, route Route) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "text/plain required", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	if err != nil {
		http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
		return
	}
	if g.Core != nil {
		g.respondCoreJSON(w, r, route, "cli", requestEventID(r), string(body))
		return
	}
	g.respondOnce(w, route, "cli", r.Header.Get("X-Yente-Event-ID"), "", string(body))
}

func requestEventID(r *http.Request) string {
	// A comma makes any duplicate-header value invalid under validID.
	return strings.Join(r.Header.Values("X-Yente-Event-ID"), ",")
}
