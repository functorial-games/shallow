# Steve Awodey — *Category Theory*

Steve Awodey, *Category Theory*, 2nd ed., Oxford University Press,
2010, Oxford Logic Guides 52.
[Second-edition record and table of contents](https://old.maa.org/press/maa-reviews/category-theory).
[Publisher description of the first edition](https://academic.oup.com/book/7134).
Requested as a foundational source for Shallow/Yente; this is not a
claim the reader owns this edition.

These are **original study notes and mathematical derivations** keyed
to the 2nd-edition topic order, not copied text or a record of completed
page-by-page reading. Yente applications are separate from claims
about Awodey's own subject matter.

## Chapters 1–4: categories, structure, duality, groups

A category consists of objects, sets (or classes) of arrows,
composable arrow operations, identities, associativity and unit
laws. An *opposite category* 𝒞ᵒᵖ reverses all arrows:
Hom_𝒞ᵒᵖ(A,B) = Hom_𝒞(B,A), with reversed composition. Duality turns
a theorem about products into one about coproducts, about initial
objects into one about terminal objects, and about monomorphisms
into one about epimorphisms, where the hypotheses genuinely dualize.

**A → B vs B → A.** Reversing the arrow is a formal operation, but
it does not turn "A's impression of B" into "B's impression of A."
Those are separately sourced data. The *opposite category* is a
category built from one category; the second person's account is
another empirical input. The distinction prevents an elegant
piece of notation from fabricating reciprocity.

A one-object category is a monoid; a group is a one-object category
with every morphism invertible. A relationship's repeated acts are
not automatically a group, because repair, departure, conversation,
and memory may lack inverses. A directed *path category* is
usually a better starting point than a groupoid. Groupoids may be
appropriate for explicitly reversible changes of representation,
not for undoing history.

## Chapter 5: limits and colimits

For a diagram D : 𝒥 → 𝒞, a limit is a universal cone to D;
a colimit is a universal cocone from D. In Set, familiar special
cases include products, equalizers, pullbacks, coproducts,
coequalizers and pushouts.

**Equalizer.** Given f,g : A → B,

    Eq(f,g) = { a ∈ A | f(a) = g(a) }.

Use this on two *comparable encodings* of a recorded fact, if equality
has a specified operational meaning. Never require two speakers'
experiences to "equalize" as a condition of being listened to.

**Pullback.** For f : A → S and g : B → S,

    A ×_S B = { (a,b) | f(a)=g(b) }.

This is useful for precise shared constraints: two sets of activity
options, each mapped to permissible date/time/venue combinations.
"Compatible" here is entirely relative to the chosen common
codomain S; it is not the property "a good couple."

**Pushout.** Given E → A and E → B, the pushout A ⊔_E B identifies
images of E in the disjoint union, in Set. We might merge two public
records when a particular common event is independently identified.
But do not merge distinct attributions or interpretations just
because they refer to the same episode. A valid data-modeling
pushout needs explicit identification rules, and may still lose
important provenance if applied indiscriminately.

## Chapter 6: exponentials

In a cartesian closed category, B^A represents arrows from A to B:
there are natural bijections

    Hom(X × A, B) ≅ Hom(X, B^A).

This is currying, not a romantic pairing operation. A category of
structured relationship-state descriptions *may not be cartesian
closed*; checking the required universal property is a research
obligation. In particular, a private person's possible responses
are not presumed to form a known complete function space.

## Chapter 7: functors and naturality

A functor F : 𝒞 → 𝒟 preserves identities and composition:

    F(1_A) = 1_FA
    F(g ∘ f) = Fg ∘ Ff.

A natural transformation η : F ⇒ G consists of components
η_A : F(A) → G(A) such that for f : A → B,

    G(f) ∘ η_A = η_B ∘ F(f).

Possible **implementation** use: interpret an authorized event
history at two levels, detailed and summarized. If both are
honest functorial representations, a compatible summarization
can form a natural transformation. A failure of the square to
commute warns that summarizing before updating produces something
different from updating before summarizing.

But real summarization can be lossy, path-dependent, and
uncertain. Do not assert functoriality before testing it, and do
not mistake algebraic coherence for consent or factual truth.

## Chapter 8: categories of diagrams

A category of diagrams 𝒞^𝒥 has functors 𝒥 → 𝒞 as objects and
natural transformations as arrows. Different stages of a
relationship could be represented by structured diagrams rather
than one "compatibility number": people, shared practices,
constraints, participation, interpretations, permissions, and
the relations among them.

The **shape of the diagram** is a modelling decision. If it
contains only the two individuals and mutual attraction, the
model excludes children, time, class constraints, shared pursuits,
and people who influence the pair. Yente's public research can
study diagram shapes without putting private people in Git.

## Chapter 9: adjunctions

F ⊣ G means there is a natural family of bijections

    Hom_𝒟(FX, Y) ≅ Hom_𝒞(X, GY).

Adjunctions arise from universal solutions, not simply
"complementary personalities." If a data model sends raw records
to a free structure, then forgets structure, the free/forgetful
pair may be an adjunction. This is a legitimate example of
how *record-construction algorithms* might be organized.
Calling two partners "adjoint" supplies no such proof.

## Chapter 10: monads and algebras; monoidal structures

A monad (T,η,μ) on 𝒞 has T : 𝒞 → 𝒞, unit η : 1 ⇒ T, and
multiplication μ : T² ⇒ T satisfying unit and associativity laws.
This structure can organize effects or repeated computational
operations. It is not a law that people are reducible to state
update functions; a conversational history has uncertain and
sometimes conflicting reports.

The second edition also adds treatment of **monoidal categories**.
A tensor ⊗ can represent composition alongside other activity,
if suitable associators, a unit and coherence are supplied.
The cartesian product is one monoidal example, but a tensor
need not mean "put these two people in a relationship."
Monoidal structure is promising for composing independent
*processes* — e.g. scheduling, messaging and disclosure checks —
provided they actually satisfy the laws.

## Yoneda: an especially instructive warning

For locally small 𝒞, Yoneda embeds 𝒞 fully faithfully into
presheaves by

    A ↦ Hom_𝒞(−, A).

An object is determined up to isomorphism by the **whole
representable functor** and its natural action, not by a handful
of observed interactions. Yente sees incomplete, permitted,
context-dependent testimony, in a category we have not even
specified fully. Invoking Yoneda to infer someone's essence
from listening history would be an abuse of the lemma.
A more defensible use is to ask: *what interactions would
distinguish two models, and what remains unobserved?*

## Precise exercises

1. Prove that an ordered path of episodes admits associative
   concatenation with empty-path identities. Show that two paths
   with the same endpoints need not be equal.
2. Form the opposite category of that path category. Explain why
   a reversed formal arrow is not evidence of a reversed action.
3. Give two explicitly different shared-constraint codomains S
   for a pullback of compatible meeting options.
4. Find a square of detailed/summarized histories that does not
   commute after one piece of evidence is corrected.
5. State what the objects, arrows, hom-sets and natural bijection
   would have to be *before* saying that any proposed process
   is an adjunction.

### Reading companions

[Lawvere–Schanuel](lawvere-schanuel.md) makes maps and universal
properties tangible; [Leinster](https://webhomes.maths.ed.ac.uk/~tl/bct/)
is especially good for comparing adjoints, representables and limits;
[Riehl](https://math.jhu.edu/~eriehl/context/) is a second rigorous
reference. See [relationship dynamics](yente-relationship-dynamics.md)
for a model that keeps the changing relationship itself in focus.
