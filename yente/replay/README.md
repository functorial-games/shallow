# Bounded checked-Idriç replay of the fictional Yente motif fixture

The [eight-case TSV](../tests/motif-revision.tsv) remains the input contract.
`Main.idric` is an explicit typed adapter for **those eight exact synthetic
turns**, with accumulated state for `M(P)`, `W(P)`, `V(P,Q)`,
`L(P,Q,t)`, `R(P,Q)` and `H`. It is not a general natural-language parser.

The replay retains attributed observations independently from interpretations,
a reported preference independently from another person's guess, directed
perspectives, event and report times, a jointly reported agreement, separate
paths with the same endpoints, and exact recipient-specific authorizations.

T04 seeds a deliberately *incorrect assistant interpretation* before replaying
the counterexample. Transitive dependent interpretations **and** their question
must be withdrawn; the underlying observations stay. A one-node-only
withdrawal is a known-bad fixture and must fail. T06's collapsed-history
mutant must fail. T08 tests default denial, an approved exact sentence,
recipient/holder changes, summary/paraphrase attempts, and fresh requests after
revocation. The earlier delivery is not retroactively erased.

## Execution

The workflow `.github/workflows/yente-replay.yml` pins canonical Idriç
compiler commit `0e8cfdd1106ebdeffd01b6c519880ea3afc7dc35`, bootstraps
it, typechecks `Main.idric`, builds the executable, and runs:

```sh
python3 yente/tests/verify-replay.py "$PWD/build/exec/yente-replay"
```

The independent process-level checker requires a successful executable run,
eight ordered case receipts, and deliberate rejection of broken fixture input
and tampered receipts. A structural check or Python-only dry run does **not**
count as Idriç execution.

## Limits

All data are fictional and transient. The test driver recognizes eight fixed
utterances, not arbitrary conversation; it does not implement an LLM, private
storage, user-authorized ingestion, authentication, deletion, secure data
isolation, open-ended dependency provenance, or privacy enforcement across a
deployed system. A passing replay is a narrow deterministic claim, not privacy
certification. Independent cross-repository ai-ci acceptance and the existing
Flexible Pipes delivery/receipt process remain separate requirements before
qualifying an integrated conversational Yente runtime.

No matching algorithm, compatibility rating, or astrological predictor exists
in this implementation.
