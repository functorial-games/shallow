# Aggregate data model

**This is an output schema, not a person-level archive.** The published dataset has no divorce cases, court docket numbers, document snippets, names, addresses, ZIP histories, respondent IDs, or participant-linked event sequences.

## Aggregate relationship events

One row in a results table describes **a cohort and a counted condition**, not a couple.

```text
dataset_id
source_publication
sample_frame
cohort_years
jurisdiction_or_population          # broad enough to avoid identification
relationship_stage                  # dating/cohabiting/married/separated
years_together_band                 # e.g. 0-4, 5-9, 10-19, 20+
comparison_group
event_or_theme_code
event_kind                          # claim / self-report / externally measured / court finding
observation_window_years
outcome                              # divorce / still married / reconciled etc.
denominator_definition
unweighted_numerator
unweighted_denominator
weighted_estimate                    # only if source supports valid weighting
confidence_interval
source_missingness
temporal_order_known
causal_design                        # descriptive / adjusted association / experimental etc.
privacy_release_status
method_and_limitations
```

Examples of eligible *questions*, not results:

- In a given longitudinal cohort, what fraction of still-married couples at year five separated in the next five years after a household employment shock, compared with those without a shock?
- How common is the theme “I did more than my share” in interviews after divorce, compared with interviews of couples still married?
- For marriages of similar duration, how often do speakers attribute the relationship's breakdown primarily to their own behavior, the partner's behavior, or a shared pattern?

**The denominator matters.** “Fraction of ever-married people divorced by age 55,” “fraction of marriages ending in divorce,” and “chance of divorcing in the next ten years, having stayed married for ten” are three different quantities.

## Aggregate language table

```text
dataset_id
speaker_population                 # e.g. individual divorced respondents
interview_stage                   # still married / newly separated / after divorce
years_together_band
speech_act_code                   # complaint, request, blame, apology, etc.
theme_code                        # money, domestic work, affection, safety, etc.
attribution_target                # self / partner / both / external / unclear
n_speakers_answering
n_speakers_coded_theme
weighted_prevalence_if_supported
n_paired_couples_if_applicable
paired_agreement_rate_if_reported
coder_agreement_or_validation
sampling_and_recall_limitations
privacy_release_status
```

A theme like “threatened with a knife” must not be coded as a neutral disagreement. Distinguish **alleged intimidation**, admitted conduct, and adjudicated or corroborated events at the aggregate level. Such distinctions should never be lost for the sake of symmetrical A/B storytelling.

## Language coding rubric

Classify utterances at source analysis time, then keep **counts only**:

- **Speech act:** accusation; description; demand; concession; plea; apology; promise; boundary; threat; offer of repair; expression of appreciation.
- **Attribution:** own actions; partner's actions; reciprocal interaction; outside event; unspecified.
- **Temporal stance:** before commitment; early marriage; transition (birth, job loss, move); prolonged disillusionment; separation; retrospective explanation.
- **Domestic obligation:** paid work, unpaid labor, parenting, caregiving, sex/intimacy, extended family, finances, housing, leisure, safety.
- **Epistemic status:** speaker-reported experience vs objectively verified circumstance. A language model must not turn an attribution into an established fact.

When paired accounts exist, report the rate of agreement/disagreement on a specified issue, including the **number of pairs**, not an identifying A/B transcript. A/B orientation can be randomized before coding; demographic grouping is an optional separate statistical analysis, not an assumed explanatory variable.

## Privacy release gate

Prior to publication:
1. Confirm the source's access terms, lawful use and statistical disclosure limitations.
2. Avoid person-level rows even if an underlying dataset is public-use.
3. Coarsen time, locations, and economic variables. Do not publish unique conjunctions of life events.
4. Suppress or combine sparse cells and also check whether totals and neighboring cells allow suppressed figures to be reconstructed. Consider privacy noise for cumulative or interactive queries.
5. Prevent repeated filters from drilling back to a small group. Keep a query budget or preapproved aggregate cube rather than unconstrained free-form SQL.
6. Record sample size, sampling/weighting, uncertainty and selection bias.
7. Use only synthetic composite speech examples, plainly marked as such. Do not offer reconstruction of original narratives.

A numerical minimum cell count by itself is **not a proof of anonymity**. Release criteria require consideration of linkage and differencing risk.

## Quality checks

- Missing is neither false nor zero.
- Employment loss can precede *or follow* relationship deterioration: do not automatically present an association as causation.
- A selected sample of divorces answers “what divorced people report,” not “what predicts divorce” without a comparison group and follow-up.
- Distinguish duration of relationship, duration of legal marriage, and age of participants.
- Keep gender-blind analyses available; gender-aware comparisons must be grounded in observed metadata and sufficient sample sizes.
- Separate prevalence of conflict from prevalence of abuse; neither is a simple unitary conflict score.
