# Fresh Yenta model boundary

The checked Idriç conversation core already implements the actual request and
admission slice. It remains the owner of evidence, permissions, corrections and
state. Fresh inference requires the Python-native Transformers/Torch object
interfaces, for which this Idriç slice has no direct binding. Use the user's
existing checked Ithon foreign-model pattern, already used by Flexible Pipes'
RHS adapters, at this narrow boundary. No replacement Python program is authored.

Types: an immutable scoped request is decoded into named records; a model
produces an untrusted structured plan; an exact typed codec serializes that plan
to the core's existing TSV interface. The adapter may read only the scoped
request supplied to it, not a person's journal or another participant's history.
The driver owns fictional sequences and asks the core for each context before
calling inference. A later independent ai-ci observer checks process results.

Effects: acquire the pinned public Qwen model on a build host, record exact model,
tokenizer, runtime, source, prompt and token identities, generate fresh output,
retain failures, and invoke the existing core admission. No paid provider or
real confidential material is used. The model is a small wiring/behavior control,
not a claim that production matchmaker quality is established.

Pin: Qwen/Qwen2.5-1.5B-Instruct at
989aa7980e4cf806f80c7fef2b1adb7bc71aa306. Published usage:
https://huggingface.co/Qwen/Qwen2.5-1.5B-Instruct

The Ithon parser/type receipt and actual inference outcome are recorded after
execution. Local stdlib controls do not count as model inference.

The first hosted BF16 trial took over fifteen minutes in the inference stage
without a completed result. It is superseded by the established float32 CPU
setting from the existing foreign-model adapter, with an explicit Linux/x86_64
and 12 GB physical-memory precondition. The slow run is not relabeled as PASS;
its retained artifacts distinguish completed calls from missing ones. The
driver now records model-load state and per-case progress as well as failures.
