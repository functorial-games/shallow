# Uncited reply after interpretation withdrawal

The independent oracle reproduces this defect against the actual pre-fix Chez
executable. A reply can refer to an interpretation without declaring that
interpretation in its reference fields. Previously, withdrawing the
interpretation removed its record but retained that reply in future model
context. Evidence correction and disclosure revocation already conservatively
cleared replies in the affected session; explicit interpretation withdrawal
did not.

The same fictional `journal.tsv` produces both request files here. The old
executable emits `UNCITED-WITHDRAWN-REPLY` after withdrawal; the corrected one
omits it. The underlying account remains. The journal retains the audit. This
is executable regression evidence using deliberately supplied response text,
not fresh model output.

The new independent checks reject the old executable at `withdrawing a
hypothesis removes uncited replies exposed to it`. All 124 named observations
pass against the repaired executable, including restart, the later actual
model request, uncited disclosure revocation and the unrelated session's reply.
These include process observations, not 124 separate semantic properties.

Executed identities:

| Material | SHA-256 |
| --- | --- |
| Old compiled Chez artifact | `86ee804c2af636a823a26b7bc058b680fd0479465d8e7b2a80a91ba30f1d747c` |
| Repaired compiled Chez artifact | `8096ded2f7a4a30b0cd5c712232a86f540f2b32a8c461fe4d8030c88f3747231` |
| Independent Accept.idric source | `ce59cf52cf469c47f64a8f10b33aaedd4bc678ffbbf94193b074fee57256291e` |
| Independent compiled Chez artifact | `bf29ad5b88188f9385a302d4a7a97ad37ce9c06c98986d36a9bc900955459996` |

Canonical Idriç source is `2bce7677f513d711875e978bd6f5ad28b9e54d58`,
using the existing Chez build path. No model execution follows from this
regression. The hosted live trial reports its own evidence separately.
