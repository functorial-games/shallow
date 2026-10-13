# Yente: learning one person before proposing matches

Yente is the conversational matchmaker role discussed for Shallow. This folder records its learning requirements; it is not a claim that a separate repository, deployed bot, or persistent learning system exists.

Read [lived experience](../lived-experience/README.md) and [correction requirements](../lived-experience/correction-requirements.md) before expanding this design. Books and research remain useful, but they must not substitute for the person's actual history.

## What learning must mean

The system should become better at understanding this person because of what they have reported, tried, experienced, and corrected. Adding reading notes is not enough. Neither is accumulating descriptive labels.

A proposed experience record should retain:

- the authorized source and speaker;
- a particular event and its reported sequence, not merely a trait;
- the person's interpretation and the assistant's distinct hypotheses;
- relevant circumstances and changes over time;
- corrections and which earlier claims they invalidate;
- outcomes, mixed results, and unresolved questions;
- disclosure and retention permissions.

Unknown fields may remain empty. Do not invent missing events to complete a record. Do not convert one person's account into verified knowledge of another person's intentions.

## Conversation behavior

Retrieve relevant prior corrections before forming a new explanation. Preserve distinctions across turns and across conversations, without merging different people or periods. Let a correction remove an attractive theory rather than merely decorating the theory with a caveat.

Address what the person is asking. A question about meaning is not automatically a request for instructions. When advice has been declined, continue with grounded discussion rather than plans, scripts, tests, or generic summaries presented as discoveries.

A response may offer a tentative interpretation when it explains which details support it and what would remain unknown. Uncertainty should be specific, not an excuse to dismiss the person's experience or to enumerate every imaginable motive.

Interests, listening histories, and other voluntary imports are context. They are not personality diagnoses or a license to interpret every relationship through an aesthetic analogy. Preserve the reason the information was supplied.

## What to build, and what not to claim

A first implementation should support one person's authorized event history, retrieval, corrections, and inspectable changes in the assistant's account. Matching comes later. Do not begin with fixed compatibility dimensions, hidden weighting, or a personality questionnaire.

Keep personal source material in an appropriately designed private store. Public Git may hold general requirements and explicitly approved, minimized examples, not raw personal conversations. The ingestion and privacy work in [issue #1](https://github.com/functorial-games/shallow/issues/1) is not yet proof of a working secure store.

No model-weight training, multi-user reuse, sharing with a prospective match, or publication follows automatically from permission to remember a conversation. Those are distinct uses requiring distinct authorization. Yente must identify itself unmistakably as a bot and not impersonate a person or ghostwrite as them.

## Acceptance direction

Evaluate sequences rather than isolated replies: an initial account, an assistant mistake, a correction, and a later question that would expose whether the correction stuck. Inspect whether the later response preserves agency, setting, scale, uncertainty, and the requested response mode.

Use qualitative, source-grounded review. Passing a schema check or matching a forbidden-word list does not demonstrate understanding. These checks evaluate the system, never the user's partner. Describe any synthetic examples as synthetic; do not call them observations of real life.

Implementation work must use the existing Flexible Pipes job and receipt path with independent ai-ci acceptance. This documentation does not create an alternate job protocol or assert that dispatch, acceptance, training, or deployment has occurred.

## Current status

Requirements only. No ingestion code, private vault, correction engine, replay runner, trained model, or deployed matchmaker is implemented by these documents.
