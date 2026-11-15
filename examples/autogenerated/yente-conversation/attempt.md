# Yenta conversational engine attempt

## Purpose and type sketch

Implement a fictional conversational prototype that preserves attributed reports,
revisable interpretations, contextual audiences and exact approved retellings.
The core uses indexed identifiers for people, sessions, evidence, hypotheses,
plans, proposals and agreements. A single checked state owns pure transitions;
effects remain at the command/store/model boundary. Records are not predictions
or psychological measures.

Central operations: project a session's permitted context; accept a grounded
turn plan; withdraw a hypothesis and transitive dependents; approve and deliver
an exact proposal revision; confirm an agreement as a named participant.
Errors leave state unchanged. Retrying a successful operation must not repeat
its visible effect. Evidence remains attributed and corrections never erase
another speaker's testimony. A private change must leave other sessions' model
requests byte-identical.

Effects: UTF-8 command and model-request files, atomic journal snapshots,
external Flexible Pipes command/model stages. Only synthetic data may be used.
Runtime and serialization validation cannot establish arbitrary model entailment.

## Implementation and verification

Language: canonical Idriç, Text/Data.Text, Number, Unicode arrows and equality.
Backend and exact compiler invocation are recorded after actual execution below.
The inherited replay is repaired independently before the reusable core is
qualified. There is no generated-C, RefC, Java or Gradle fallback.

## Observed execution

The compiler used was the existing canonical material executable reporting
`2bce7677f513d711875e978bd6f5ad28b9e54d58`, with Chez backend. Its compiler
`.so` SHA-256 is `2e1fd91e9e3c254428c8cf135e56b4d964114bbea53e10a70fa15d3c81051e72`.
A canonical source checkout at that revision was read for syntax and style;
this work did not rebuild that compiler material from source.

The package builds through `idris2 --build yente-conversation.ipkg`. The actual
Linux build used prelude/base/contrib TTC libraries, Chez, and the compiler's
support and support/c directories through `IDRIS2_PATH`, `CHEZ`, `IDRIS2_DATA`
and `IDRIS2_LIBS`. The executable was then run, not inferred from compiler exit
status. The precise build recipe is also in the hosted qualification workflow.

The original eight-case replay and its five hostile controls passed after
flattening a layout-sensitive lambda and declaring the file-IO entrypoint
`covering`. The new independently owned ai-ci process oracle passed against
the compiled conversation executable. It exercises novel identifiers, source
kinds, transitive corrections and restart, authorization, exact agreement
confirmations, disclosure editing/revocation, two-world model-request equality,
recorded-response admission, stale request refusal, corrupt recovery refusal,
and explicit event retry identities.

A later replay review exposed a cached-wording defect: retrying an event-wrapped
reply could repeat text after its supporting evidence was withdrawn. Reply
retries now emit only an opaque receipt, and read operations cannot be cached as
events. Independent cases exercise both correction and revocation before a
retry, check no duplicate history, and require wrapped reads to leave the journal
unchanged. The earlier implementation is deliberately rejected by these cases.

The Store boundary additionally executed thirteen hostile/roundtrip controls,
including strict UTF-8 rejection, checked flush/close, lock refusal and failed
replacement preserving committed bytes. These supplemental probes are not a
production storage qualification.

## First exposed language boundary

No missing Idriç capability required a fallback for the core. Linux journal
replacement uses a narrow declared libc `rename`/`fflush`/`fclose` boundary.
Two compiler failure modes were observed: layout errors can occur in multiline
`do`/`if` expressions, and a missing support directory can print an internal
build error while the compiler returns zero. The app source was corrected and
an existing compiler library path supplied. A future compiler change should
make a missing support file return nonzero and prove artifact creation before
reporting a successful build. This is a compiler-shaped acceptance requirement,
not an app fallback.

Fresh model inference, automated natural-language interpretation, subjective
voice quality, production private storage, deployment and physical-device
behavior remain unqualified. Recorded responses are explicitly authored
fictional material, not observations of a real provider call.
