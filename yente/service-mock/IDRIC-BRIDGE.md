# Yente fictional transport through the Idriç journal

This is the proposed integration of [shallow PR #12](https://github.com/functorial-games/shallow/pull/12) with the **existing** Idriç conversation core from [PR #11](https://github.com/functorial-games/shallow/pull/11), including the `ask-next` and `reply` commands in [PR #13](https://github.com/functorial-games/shallow/pull/13). It is deliberately **not** a new conversation engine and does not implement a matching model.

The `-core`, `-state-root`, `-materials` options enable stateful `/v1/cli/ROUTE` requests. Add `-fixture-tokens` to enable the same journal through both FastChat stream formats and simulated SMS. Without core options, the gateway returns the original transport acknowledgements.

## What goes into the journal

The gateway parses explicit fictional commands and supplies the server-owned route's session and actor to the original `yente-conversation script` process. It does **not** accept a caller-chosen actor, session, or disclosure recipient. Event IDs are derived from route, channel and caller event ID, and passed to the core's existing persistent `event` receipts. After each accepted command, the gateway invokes the existing `show SESSION` projection, parses the Idriç length-framed rows, and reports only a bounded, route-scoped summary. Core `ERR` records are refusals even if the process exit code is zero. The gateway serializes script and show within one instance and fails closed on missing or mismatched scoped output.

The supported single-line message bodies are:

| Message | Meaning |
| --- | --- |
| `/report TEXT`, `/wish TEXT`, `/hearsay TEXT`, `/hypothetical TEXT`, `/agreement-report TEXT` | **Explicitly** attribute the exact account to this route's fictional participant. No type is guessed. |
| `/ask` | Run the Idriç `ask-next` person-model inquiry. Return its active question ID and wording. |
| `/reply QUESTION-ID KIND TEXT` | Attribute an answer to an Idriç question (kind as above). |
| `/correct EVIDENCE-ID KIND TEXT` | Supersede the participant's own account; the core invalidates dependent interpretations/wording. |
| `/context` | Report scoped current account/interpretation counts and the latest active question, never the original account text. |
| `/withdraw-consent` | Set discussion and retention permissions to `no` in the Idriç session. Subsequent accounts are refused; the audit journal remains but the gateway withholds historical metadata. |

CLI and FastChat mutations require `X-Yente-Event-ID`. SMS mutations use `message_id` in the JSON body. Input is validated as a single clean UTF-8 command line. The gateway rejects bare chat text and refuses shared routes. Each accepted operation belongs to the session and actor in the route table.

## Stateful FastChat protocol and SMS fixtures

Pass `-fixture-tokens` a private regular file with `ROUTE<TAB>TOKEN` rows. Use a distinct token of at least 24 letters, digits, underscores, hyphens or periods for each enabled route. Send that token in `Authorization: Bearer TOKEN` on all three interfaces. Keep the file outside the checkout. A configured route without a token is inaccessible in this mode.

FastChat Chat Completions and Responses requests use the last user message as the command. Earlier history, including local system status messages and empty assistant replies, has no effect on the journal. Send the same event ID when retrying the same operation, including across these two FastChat formats. A successful request streams a scoped reply and the format's terminal event after the core operation and projection return. `X-Yente-Replayed: true` identifies a recognized retry. The current native app does not supply this event header. The host acceptance sends the native wire shapes directly.

An SMS JSON fixture uses `sender` equal to the route's configured fictional actor, `message_id` for the durable operation ID, and `text` for the same command used by CLI. The bearer token grants fixture access, and `sender` confirms the fixture mapping. The HTTPS response contains a JSON receipt. No carrier delivery occurs.

## Local fictional run

Build the original Idriç executable as documented in `yente/conversation/README.md`. For its pinned CI toolchain, `./.toolchain/Idric/idris2 --build yente-conversation.ipkg` yields `build/exec/yente-conversation`.

Initialize a *fresh* trusted, absolute journal directory using the supplied original Idriç `script` command:

```sh
repo="$PWD"
state_root="$(mktemp -d)"
"$repo/build/exec/yente-conversation" script "$state_root" \
  "$repo/yente/conversation/materials" \
  "$repo/yente/service-mock/fixtures/core-seed.tsv"
```

The seed includes only fictional Lio, Mira, Nera and Suri, separate private sessions, and two open person-model questions. A `couple-shared` session exists but has no permission and no transport entry point.

Start the mock TLS listener (certificate setup remains as in PR #12):

```sh
cd "$repo/yente/service-mock"
export YENTE_MOCK_TOKEN='some-long-random-test-token'
go run . -listen 127.0.0.1:8443 -routes routes.synthetic.tsv \
  -cert /tmp/yente-mock.crt -key /tmp/yente-mock.key \
  -core "$repo/build/exec/yente-conversation" \
  -state-root "$state_root" \
  -materials "$repo/yente/conversation/materials"
```

The provided `yente-curl.sh` still accepts the same four arguments. Set `YENTE_CA` to the certificate/CA file, then from a second shell:

```sh
./yente-curl.sh https://127.0.0.1:8443 partner-a ask-1 '/ask'
./yente-curl.sh https://127.0.0.1:8443 partner-a reply-1 '/reply lio-private-p1 wish I would like weekly visits.'
./yente-curl.sh https://127.0.0.1:8443 partner-a say-1 '/wish I prefer Sunday afternoons.'
./yente-curl.sh https://127.0.0.1:8443 partner-b context-1 '/context'
./yente-curl.sh https://127.0.0.1:8443 partner-a fix-1 '/correct lio-private-e2 wish I meant monthly afternoons.'
./yente-curl.sh https://127.0.0.1:8443 partner-a context-1 '/context'
./yente-curl.sh https://127.0.0.1:8443 partner-a consent-1 '/withdraw-consent'
./yente-curl.sh https://127.0.0.1:8443 partner-a say-2 '/wish This must be refused.'
```

The `-core` executable is real: no replacement journal interpreter is used at runtime. The expected IDs above assume a new state directory and successful seed. Responses identify themselves as **deterministic fictional mock**, not model-generated prose. Journal event retries survive gateway restarts; old generated question wording is *not* replayed by the core after withdrawal/correction.

The independent recorded-reply fixture can take a transport event through the checked Ithon provider codec and the core's `accept-event` command. The [conversation contract](../conversation/README.md#model-request-and-response-framing) describes admission and `recover-answer`. `TestRealCoreAdmissionRecovery` exercises a committed SMS operation, an authored plan, recovery in a fresh core process and withholding after correction or consent withdrawal. The HTTPS gateway still returns its deterministic summaries. Connecting the provider to HTTP is a separate integration step.

## Verification and limits

Run `go test -race ./...` and `go vet ./...` in this directory. Set `YENTE_CORE_BINARY` and `YENTE_CORE_MATERIALS` to canonical absolute paths to exercise the compiled core. `TestRealCoreE2E` covers CLI questions, replies and corrections. `TestRealCoreAcrossFixtureTransports` writes a FastChat command, corrects it through simulated SMS, checks through CLI, restarts against the same journal and tests collision and withdrawal behavior. Both tests skip if the executable or materials are missing.

This remains fictional-only. The route label and shared bearer token are not authentication. Storage is the existing single-process trusted-path journal with its documented atomic rename and lock limits; it is not encrypted or a multi-process transaction. Never use this with real relationship stories or deploy it as a participant service. Consent withdrawal withholds future gateway context and blocks new retained accounts; it does not erase the original audit journal. No compatibility predictions, LLM intelligence, live SMS or FastChat stateful UI are claimed.
