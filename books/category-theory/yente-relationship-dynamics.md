# Relationship dynamics: categorical sketches for Yente

**Aim.** Study how **the relationship itself** evolves, not merely how
A affects B or B affects A. Changes in either person are another
possible process, and the two processes can interact without being
identical. These are possible mathematical data models, not a theory
of attraction or evidence about any particular people.

Read with [Lawvere–Schanuel](lawvere-schanuel.md) and
[Awodey](awodey.md). Actual Yente requirements live in
[its project notes](../../yente/README.md).

## 1. Two independent state changes and one relational change

Fix a time or observation index t. Distinguish:

- A_t: facts, preferences and circumstances attributed to A at t;
- B_t: likewise for B;
- R_{AB,t}: the **relational state**, including shared history,
  expectations, commitments, actions, unresolved tensions and
  developments that need not belong to A or B alone;
- V_{A→B,t}: what A reports experiencing or inferring about B;
- V_{B→A,t}: the distinct reverse perspective;
- H_{Y;A,B,t}: Yente's fallible hypotheses, carrying sources,
  uncertainty and correction history.

The symbols are names for possibly different **types**, not one
vector space. A_t, V_{A→B,t}, and H_{Y;A,B,t} must not be silently
identified. Yente also holds perspectives involving C, D, and
others, and may have different hypotheses about what each wants.

A transition R_{AB,t} → R_{AB,t+1} is not shorthand for
(A_t → A_{t+1}, B_t → B_{t+1}). A couple might develop an
agreement that changes the relationship without a general
change in either individual's temperament. Conversely, one
person can change jobs without the relational state changing
in the same way, though it may later matter.

More broadly, the relational object should include effects of
institutions, labor, housing, money, shared projects, children,
and other participants where relevant. It cannot be reconstructed
reliably from a static pair of personality vectors.

## 2. A mathematically clean first construction: path category

Take a **synthetic**, directed multigraph G:

    r₀ -- invitation --> r₁
    r₁ -- misunderstood expectation --> r₂
    r₂ -- explicit clarification --> r₃
    r₁ -- no reply --> r₄

The vertices are *states of the relation*, not the people.
The arrows are labeled possible transitions. In the **free
category** Path(G):

- Objects are the vertices r₀, ..., r₄.
- Arrows are all finite directed paths, including empty paths.
- Identity at r is the empty path at r.
- Composition is path concatenation.

Associativity and unit laws follow directly from concatenation.
The composite r₀ → r₁ → r₂ → r₃ is a path from r₀ to r₃.
It is **not** a causal claim that invitations cause clarification,
and it is not equal to an arbitrary direct r₀ → r₃ arrow.

If two paths from r₀ to r₃ are distinct, retain both. Quotienting
them by an equation would erase a difference in history; such
identification needs evidence or a declared abstraction.
This is why a thin temporal preorder, with just one arrow s → t
whenever s ≤ t, can be too coarse for relationship history:
it hides competing paths between the same dates.

The path category is a recording structure, not a probability
model, a dynamical law, or an algorithm for deciding what to do.

## 3. The third object: relationship versus cartesian product

In Set the product A × B contains **all possible ordered pairs**
(a,b). It is a poor representation of what A and B have made
together. A particular relation E ⊆ A × B may select compatible
or observed pairs, but still omits a shared history unless E's
elements are richer episodes.

A less reductive design is a **span** of source-preserving records

                p_A          p_B
            A  <--- E_AB ---> B

where E_AB is a collection of attributed shared episodes,
and p_A,p_B indicate relevant participation or account.
A and B here must denote appropriate *record spaces* — not
person identities whose singleton product would tell us nothing.
The object E_AB can store extra context that neither projection
retains. For a truly emergent relationship state, keep a further
structured object R_AB whose maps to participation records
are optional model choices, not a definition of love.

One can represent **shared objects of attention** too:

    A --> music/work/children/place <-- B

This triangle is not reducible to "A likes B" and "B likes A."
It permits a relationship to be constituted by making, caring,
learning or arguing about something beyond the pair.

## 4. Two viewpoints: opposition does not guarantee reciprocity

Let f : A → B denote an attributed directed interaction, and
g : B → A a response. The cycles g∘f : A → A and f∘g : B → B
exist in a category of directed interaction paths. But generally

    g ∘ f ≠ 1_A;      f ∘ g ≠ 1_B.

Thus mutual responses are not automatically inverses or an
isomorphism. Calling the pair a "natural isomorphism" requires
a much more elaborate specified structure and equations.

The opposite category 𝒞ᵒᵖ reverses arrows **formally**. It does not
supply B's first-person testimony. Independently attributed
perspectives should remain independently stored.

**Synthetic contrast.** A says a remembered detail made them feel
understood. B says that remembering it had been burdensome.
These claims are not logical contradictories: they concern
different experiences of the same episode. Treating them as
equal or canceling them is unjustified.

## 5. A relationship changes; a model of it also changes

Let E be the category of authorized evidence revisions and
historical events, and S a category of summaries. If a functor
F : E → S were defined, it would preserve composition and
identities. A correction e : x → y would entail an update
F(e) : F(x) → F(y).

In practice there are two processes:

(a) a real or reported event changes the **relationship**;
(b) new testimony changes **Yente's account** of the relationship.

They can occur at different times and must not be conflated.
For example, two people may have reconciled last month, while
Yente learns this today. Updating H_Y does not retroactively
cause R_AB to change today. Conversely, Yente's earlier
classification may have been false from the outset.

Corrections should identify dependent claims; withdrawing
one hypothesis should not erase a historical episode. In
categories of evidence states, a correction arrow records
what was revised, not an invertible "undo" of lived history.

## 6. Precise place for universal constructions

**Pullback of shared constraints.** Let A be A's feasible meeting
options, B B's, and S shared time-slot labels. For maps
f : A → S and g : B → S,

    A ×_S B = { (a,b) | f(a)=g(b) }.

This construction precisely answers: "which two options have
the same slot?" It does not answer whether either participant
wants the meeting or whether the pairing is healthy. To
include accessibility, location and costs, choose a richer
codomain S and check whether the resulting matching options
change.

**Coproduct and provenance.** If reports P_A and P_B come from
different speakers, P_A ⊔ P_B preserves the origin tag. This
prevents a convenient common story from obliterating who said
what. An evidence-based identification or a pushout may
sometimes be justified, but mere mention of the same event
is not permission to merge interpretations into one fact.

**Limits, agreement, disagreement.** A commutative diagram is
a property of the chosen maps. If following two interpretive
paths gives different results, record the noncommuting square
rather than forcing consensus.

## 7. Monogamy: a combinatorial layer, not a category-theoretic law

For a set P of people at time t, a strict pairwise monogamous
matching can be represented by a relation M_t ⊆ P × P with:

- no self-pairs: (p,p) ∉ M_t;
- symmetry: (p,q) ∈ M_t ⇔ (q,p) ∈ M_t;
- at most one partner: for every p, there is at most one q
  with (p,q) ∈ M_t.

Equivalently, on the set of currently paired people one can
use a **fixed-point-free involution** m_t with m_t² = identity;
on all P it is naturally a partial involution.

This representation tracks exclusivity **at a time**, not love,
desire, duration, consent, cohabitation, parenting, or the
quality of the bond. Distinguish whether the match is a
self-description, mutually agreed commitment, legal status,
or observer inference. These are not interchangeable.

As t changes, M_t can change independently of many properties
of R_{AB,t}. A pair can remain formally exclusive while the
relationship changes radically. A relationship can continue
to have obligations or consequences after the matching edge
is removed. No mathematics here favors monogamy over other
arrangements.

## 8. More than two people; permissions

Yente must model A, B, C, D, ..., including:

- A's description of B and B's description of A;
- what A thinks B wants, and what Yente thinks A wants;
- friendship, kinship, shared endeavors, historical influence,
  and competing interactions, not only proposed matches;
- permissions about what may be remembered and what may be
  revealed, with viewer- and subject-specific rules.

A perspective is not a public fact. The operation "show C the
permitted portion of Yente's model of B" needs an explicit
authorization and disclosure check. Neither a product
projection nor a natural transformation by itself proves
privacy safety. This public repository must contain neither
real private case histories nor reconstructed disclosures.

## 9. Comparison table: mathematical claim versus overreach

| Concept | Mathematical statement | Useful possibility | Invalid shortcut |
| --- | --- | --- | --- |
| Morphism | Directed arrow with type | Directed act or state transition | Treat every feeling as a total function |
| Composition | Associative paths | History of sequential events | Pretend a composite proves causation |
| Isomorphism | Inverses satisfying two identities | Distinguish representation changes | Equate reciprocal affection with inversion |
| Product | Universal paired object | Paired reports/options | Assert the product *is* the bond |
| Pullback | Fiber product over S | Match explicit shared constraints | Infer holistic compatibility |
| Functor | Preserves identity/composition | Consistent event abstraction | Say every LLM summary is functorial |
| Naturality | Coherent components of functor maps | Check summary/update order | Pretend disagreement commutes |
| Monoidal tensor | Parallel composition with coherence | Compose independent processes | Claim ⊗ means romance |
| Yoneda | Representables determine objects up to iso | Think carefully about observations | Infer a person from partial behavior |

## 10. Tests before any implementation

1. Use only explicitly labeled **synthetic** cases in a public fixture.
2. Two histories with same endpoints must remain distinguishable.
3. A new account about a *past* relationship change updates the
   model, without inventing a new relationship event at report time.
4. A and B's reports of the same event retain attribution and
   can disagree or concern different aspects.
5. A correction withdraws claims that depended on a false
   premise; it does not erase the supported observations.
6. No inference of B's consent from A's report about B.
7. An authorized display of B's information to C must be
   explicitly permitted; a mathematical projection isn't enough.
8. No single compatibility number or predefined personality
   vector is required to store a relationship's evolution.
9. A three-person scenario must be representable without
   forcing the third person into one of the pair slots.
10. A mathematical analogy must list its assumptions and a
    concrete failure case before being adopted.

**Status:** research notes only. No categorical data model,
privacy mechanism, evaluation suite, or deployed Yente system
is implemented here.
