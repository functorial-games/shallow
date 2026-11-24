# Yente transport mock (not the matchmaker)

**Opt-in stateful fixtures.** Configure the existing Idriç executable with `-core` for CLI access. Add `-fixture-tokens` to connect FastChat-shaped HTTPS and simulated SMS requests to the same journal. See [IDRIC-BRIDGE.md](IDRIC-BRIDGE.md) for the fixture contract.

This test gateway accepts SMS-shaped inbound JSON, [FastChat](https://github.com/isomorphisms/fastchat) HTTPS/HTTP2 streaming requests, and command-line `curl`. The default fixture returns canned acknowledgements. Stateful fixture replies summarize the Idriç journal for a fictional route. They contain deterministic text and can include a question selected by the core.

It belongs beside [Shallow PR #11](https://github.com/functorial-games/shallow/pull/11), whose `yente-conversation` executable already owns attributed evidence, permissions, revisions, correction and revocation. This gateway does not implement a rival conversation core. Its small in-memory receipt table only detects test retries and is lost on restart.

## Contract

The five fictional entries in `routes.synthetic.tsv` map route names to fixture sessions, actors and audience roles. The URL names the route. In stateful fixture mode, each route has a distinct bearer token, and the gateway verifies that the configured actor belongs to the Idriç session before sending a command. Shared-session requests are refused. These routes are for fictional participants only.

| Client | POST path | Body / response |
| --- | --- | --- |
| Native FastChat | `/v1/chat/completions/partner-a` | `{"model":"yente-mock","stream":true,"messages":[{"role":"user","content":"..."}]}`; HTTP/2 SSE Chat Completions delta then `[DONE]` |
| FastChat Responses format | `/v1/responses/partner-a` | Same model/stream, `input` array; HTTP/2 SSE `response.output_text.delta`, `response.completed`, `[DONE]` |
| SMS transport simulator | `/v1/sms/partner-a` | JSON `{ "sender":"synthetic-a", "message_id":"sms-1", "text":"..." }`; JSON receipt, never an actual outbound SMS |
| CLI / Shizuku rish | `/v1/cli/partner-a` | `text/plain` body and `X-Yente-Event-ID`; JSON receipt |

Event IDs are scoped to route and transport. In stateful mode both FastChat formats share one namespace. The journal deduplicates retries and rejects a reused identity with different command text. A retry gets a fresh summary under current permissions. The default canned fixture keeps its original process-local SMS and CLI receipts.

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

The remaining service work includes durable operation IDs in the native FastChat client, real participant and SMS provider authentication, an outbound delivery sink, core admission of fresh model replies, shared-session permissions and independent Flexible Pipes/ai-ci acceptance.
