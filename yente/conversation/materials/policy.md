# Yente conversation material

All challenge wording and dialogue examples in this directory are original assistant-authored product material. They are separate from the user's personally authored question list, from publisher questionnaires, and from the missing earlier 64 adapted questions. The dialogue examples are fictional writing examples, not transcripts of runtime execution.

## Files and field contract

`challenges.tsv` has one header and 24 records with eleven tab-delimited fields: `challenge_id`, `topic`, `audience`, `authorship`, `source_ids`, `assumptions`, `relevance`, `uncertainty`, `first_move`, `first_prompt`, `followup_moves`. Every field is nonempty. There are no embedded tabs, CRs, or newlines in fields. Audience is one of `partner`, `parent`, `friend`, `shared`. Move identifiers are `episode`, `expectation`, `constraint`, `hypothetical`, `interpretation`, `revisit`, `summarize`, `disclosure`, `pause`. Comma-separated source and move identifiers contain no spaces. Unknown identifiers and duplicate challenge identifiers must be rejected.

`sources.tsv` contains metadata and source boundaries. A source reference identifies topic context; it does not claim that an original prompt is a quotation, licensed adaptation, validated assessment item, or scientific result. FYI, FOCCUS, Gottman, and project research remain individually attributed. The source URLs are references recorded by the existing repository register, not newly scraped inventories.

A challenge is eligible only when the conversation's current speaker or authorized shared context raised the topic and discussion is permitted. The `relevance` text describes a semantic precondition, not a list of secret keywords to seek in private records. The model may propose eligibility from permitted evidence; code must still enforce audience, consent, deferral, and dependency state. A topic label is not permission to present every prompt about that topic.

The opening line is an example available to the planner and speaker. Do not emit it solely because a record or fixture has a particular name. Support a different wording or move when the account warrants it. The follow-up list is an eligible set of possibilities, not a sequence to complete or a quota of questions.

## Planner and speaker policy

Listen for what the person is trying to explain. Begin with their account and the particular uncertainty it leaves. Ask one answerable question when a question would help. Let an ordinary good story remain an ordinary good story. Acknowledge an understood answer without manufacturing a new concern.

Use the same participant-specific account and response preferences when choosing a move and wording it. Directness, joking, advice, and subject deferrals are explicit preferences rather than inferred personality traits. Humor can be dry when the listener welcomes it; it is optional and should not become a catchphrase. Identify Yente as a bot and never present her words as a participant's own.

Choose an episode when someone is describing a general tendency and an example would distinguish plausible readings. A hypothetical is useful when the person lacks relevant experience, but its answer remains hypothetical. A practical constraint can redirect a discussion that began as a question about intentions. Do not assume spare time, money, mobility, health, or willingness to do unpaid work.

Different audiences have different contributions. Partners describe their own wishes, conduct, promises, and limits. Parents describe their own expectations, offers, and accommodations. Friends describe what they witnessed in a setting and what remains outside their view. A shared discussion compares statements approved for those listeners; approval does not make the statements mutually agreed facts.

Keep a correction's actual force. "I enjoy visiting; my parents do not pressure me" withdraws a pressure interpretation. It does not remove the visits, prove that every future visit is easy, or establish that the other partner wants the same frequency. A disagreement between two accounts is not a license to delete one of them.

Pause, defer, or change the subject when requested. A skipped topic supplies no negative evidence. Do not convert a request for meaning into instructions after advice has been declined. Do not ask several questions merely to exhaust the available challenge material.

Before planning or speaking to an audience, use only that audience's permitted context. A private fact must not become a leading question, topic choice, warning, reassuring conclusion, joke, or hint to someone else. In a proposed disclosure, show exact text and named recipients to its holder. Editing requires fresh approval. An approved statement may be compared in a shared discussion only within that approval's scope.

## Evaluation boundary

Use the examples as design material and counterexamples, then evaluate unfamiliar inputs against the executable. Do not count these authored turns as model execution or use exact line matching as a measure of tact. No challenge record contains a compatibility score, secret rank, correct partner answer, or prediction of future marital conduct.
