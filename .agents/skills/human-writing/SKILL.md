---
name: human-writing
description: This skill MUST be used by default for any user-facing writing, including answers to users, emails, documents, formal correspondence, reports, proposals, rules, prompts, skill instructions, project introductions, README files, repository descriptions, landing-page copy, social posts, essays, commentary, explanations, summaries, product copy, internal updates, and other prose intended to be read by people. Apply it whenever the output contains human-readable text. This skill MUST NOT be skipped unless another skill explicitly owns the writing style or the output is purely code, structured data, or a command. Supports Chinese and English.
metadata:
  version: "1.13.0"
---

# Human Writing

## Overview

Use this for prose that a person will read, whether you are writing from scratch or revising a draft. The goal is prose the reader can picture and act on, with no AI tells left in.

Do not apply it to code, JSON, SQL, shell commands, or other structured output. When another skill owns the format or subject matter, use that skill to plan the content and use this one for the prose.

Project introductions, README files, and repository descriptions also use this skill. In mixed technical documents, apply the writing rules to prose while preserving the correctness of commands, code examples, and API identifiers.

This file gives the sequence. The rules are in the linked references. Read the reference named at each step instead of relying on a remembered summary.

## Workflow

1. Choose the mode.
   Read [references/mode-selection-and-language-rules.md](./references/mode-selection-and-language-rules.md).
   Decide between `rewrite` and `generate`, identify the reader and purpose, match the language, and choose `grounded` or `voiced`.

2. Learn the principles and examples.
   Read [references/human-writing-principles.md](./references/human-writing-principles.md), then read [references/human-writing-examples.md](./references/human-writing-examples.md) in full before drafting. Study how the examples turn the rules into natural sentences, including rewrites that still leave too much for the reader to infer.

3. Write the draft.
   When editing existing content, first apply principle 8A to choose the owning passage and reconcile related sections before adding text.
   Apply the principles and what the examples demonstrate about concrete detail, reader context, structure, vocabulary, rhythm, tone, and punctuation.

4. Remove AI patterns.
   Read [references/ai-writing-patterns.md](./references/ai-writing-patterns.md). Remove every listed pattern without changing the meaning, facts, or intent.

5. Check the finished text.
   Read [references/final-checklist.md](./references/final-checklist.md). Test every item against the actual text, fix failures, and check it again.

6. Report ongoing task progress.
   At every turn end for any ongoing task, append a compact progress bar to the conversational reply, not the deliverable. Anchor it to the original goal: verified completed/total agreed items, current work, and what remains before completion. Show pending acceptance separately. Keep the denominator stable unless scope is explicitly changed; when the total is unknown, show milestones without a fabricated percentage or finish date. 100% requires all agreed completion and delivery conditions to be met.
