# Yente transport mock (not the matchmaker)

**New opt-in mode:** the CLI route can invoke the *existing* Idriç journal and fictional person model. Read [IDRIC-BRIDGE.md](IDRIC-BRIDGE.md) first. Without `-core`, all behavior remains the original transport-only mock. FastChat and SMS remain transport-only in either mode.

This small Go-standard-library test gateway mocks three ways of reaching **one** Yente service: SMS-shaped inbound JSON, existing native [FastChat](https://github.com/isomorphisms/fastchat) HTTPS/HTTP2 streaming provider requests, and command-line `curl` (which may be invoked from a Shizuku/rish shell). **No production service or participant account exists here.** By default every route returns an explicit canned acknowledgement. With the optional `-core` integration, only `/v1/cli/` invokes the original Idriç journal. This gateway does not invoke a language model, authenticate people, send SMS, evaluate a relationship, or authorize disclosure.

It belongs beside [Shallow PR #11](https://github.com/functorial-games/shallow/pull/11), whose `yente-conversation` executable already owns attributed evidence, permissions, revisions, correction and revocation. This gateway does not implement a rival conversation core. Its small in-memory receipt table only detects test retries and is lost on restart.

## Contract

The five **fictional** route entries in `routes.synthetic.tsv` map route name -> existing fixture session -> fictional actor -> audience role. All three clients choose the route in the URL, not in a spoofable message body. The same `mockReply` function handles all clients. An unverified SMS sender is not an authenticated participant. Shared-session routing is refused: it requires the real Idriç consent/delivery checks, not a canned mock. These routes must never be deployed with real participants.

| Client | POST path | Body / response |
| --- | --- | --- |
| Native FastChat | `/v1/chat/completions/partner-a` | `{"model":"yente-mock","stream":true,"messages":[{"role":"user","content":"..."}]}`; HTTP/2 SSE Chat Completions delta then `[DONE]` |
| FastChat Responses format | `/v1/responses/partner-a` | Same model/stream, `input` array; HTTP/2 SSE `response.output_text.delta`, `response.completed`, `[DONE]` |
| SMS transport simulator | `/v1/sms/partner-a` | JSON `{ "sender":"synthetic-a", "message_id":"sms-1", "text":"..." }`; JSON receipt, never an actual outbound SMS |
| CLI / Shizuku rish | `/v1/cli/partner-a` | `text/plain` body and `X-Yente-Event-ID`; JSON receipt |

Event IDs are scope-bound to route and transport. An identical SMS/CLI retry replays the same acknowledgement; a reused ID with different text is rejected. The mock has no cross-window history. Session-specific context and disclosure will be supplied by the *existing* Idriç core, not by the transport route or an AI guess.

## Test

From this directory: `go test -race ./...`.

Tests use the standard Go TLS HTTP/2 fixture and cover FastChat's two wire formats, streaming terminators, SMS/CLI common routing, collisions, separate fictional sessions, refusal of shared disclosure, malformed messages and authentication refusal. None of these tests counts as a model-quality or private-data security acceptance.

## Running locally

TLS is mandatory, including on loopback, because native FastChat's provider requires HTTPS/HTTP2. Supply a certificate and key with a SAN matching the host. For a temporary localhost-only certificate:

```sh
openssl req -x509 -newkey rsa:2048 -nodes -days 1 \
  -keyout /tmp/yente-mock.key -out /tmp/yente-mock.crt \
  -subj '/CN=localhost' \
  -addext 'subjectAltName=DNS:localhost,IP:127.0.0.1'
export YENTE_MOCK_TOKEN="$(openssl rand -hex 24)"
go run . -listen 127.0.0.1:8443 -routes routes.synthetic.tsv \
  -cert /tmp/yente-mock.crt -key /tmp/yente-mock.key
```

From another terminal, with the **same** token in its environment:

```sh
export YENTE_CA=/tmp/yente-mock.crt
./yente-curl.sh https://127.0.0.1:8443 partner-a command-001 'A fictional message'
```

For native FastChat, configure provider `model = yente-mock`, format `Chat Completions`, origin `https://<server-host>:8443`, endpoint `/v1/chat/completions/partner-a`, a trusted CA matching the server certificate and `Authorization: Bearer <token>` in its private provider settings. Its existing HTTP/2 parser accepts a Chat Completions SSE delta and terminal `[DONE]`. The app currently has a durable single conversation and a configurable provider. Multiple *simultaneous* person-specific windows are a further FastChat UI change; this mock demonstrates distinct routed endpoint contexts, **not** that the multiwindow UI already exists. To use a device on another host, supply a cert with that host/IP SAN, an authenticated reachable TLS listener, and actual CA trust. Loopback on a separate server is not reachable from an Android phone.

The same `curl` request can run inside `rish` **only if** that environment has an executable `curl` and can access the trusted CA and token. A Shizuku shell normally does not get `curl` automatically, and ordinary HTTPS does not require Shizuku privileges. A regular Termux shell with curl is equally valid. Do not disable TLS verification to work around a missing CA.

Non-loopback listeners require `YENTE_MOCK_TOKEN` of at least 24 characters. A shared test token is not person authentication; no real identities or personal messages should be routed into this mock.

## Required before a real service

Replace `mockReply` with one bridge to the existing Idriç `script`/`request`/`accept` boundary and independent Flexible Pipes/ai-ci evidence. Verify participant identity independently of phone number, URL and client-supplied messages. Do **not** silently convert every chat line into a factual `say` report: speech, question, event account and hypothesis are different acts. Respect exact journal event IDs, session permissions, corrections, deferrals, per-recipient disclosure approvals and stale-response rejection. Use real private storage and TLS credentials; add a signed SMS provider adapter and an actual delivery path. No general conversation or model quality is established by a passing transport test.
