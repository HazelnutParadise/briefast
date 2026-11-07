# Mode selection and language rules

Read this before you start writing.

## 1. Rewrite vs generate

- Use `rewrite` when the user gives you text, a draft, or a paragraph to improve.
- Use `generate` when the user gives you an intent, bullets, notes, or context and wants new prose.
- When the request includes both notes and a rough draft, use `rewrite` and keep the parts that already work.

## 2. Grounded vs voiced

Choose `grounded` when the text is:

- professional
- factual
- internal
- operational
- explanatory
- trust-sensitive

Choose `voiced` when the text is:

- persuasive
- editorial
- opinionated
- social
- brand-led
- meant to sound personal or distinctive

When the prompt is ambiguous:

- default to `grounded` for business, documentation, updates, or neutral explanation
- default to `voiced` for essays, commentary, posts, and marketing-style requests

## 3. Rewrite mode rules

Use this mode when the user provides source text and wants it to read better.

1. Preserve the original claim, purpose, and factual boundaries.
2. Keep useful structure unless the structure itself causes stiffness.
3. Remove obvious AI tells first, then fix rhythm and voice.
4. Improve precision before adding flair.
5. If the source is already strong, edit lightly.
6. When plainness would cost precision, keep the precision. Preserve a hedge, condition, or complex clause when simplifying it would distort the meaning. Explain these editing choices when the user asks for the reasoning behind the revision.

## 4. Generate mode rules

Use this mode when the user wants fresh prose from an idea, outline, notes, or a vague request.

1. Infer the implied audience and intent from the prompt.
2. Choose the narrowest effective format instead of over-explaining.
3. Write as if a competent person with a point of view sat down to say it.
4. Do not pad thin inputs with generic significance language.
5. When details are missing, stay concrete but avoid invented facts.

## 5. Grounded style

Use for business writing, internal updates, memos, documentation, proposals, factual summaries, and other trust-sensitive text.

- Prefer direct claims and plain syntax: state the conclusion first, and hedge only where the claim needs it.
- Use specific nouns and verbs over abstract framing.
- Keep tone natural without sounding casual for no reason.
- Let credibility come from clear claims and supporting facts.

## 6. Voiced style

Use for commentary, opinion, social posts, persuasive copy, essays, and prose that benefits from character.

- Allow harder rhythm breaks and blunter word choices than grounded writing would take. State an opinion outright when the piece calls for one.
- Add personality through perspective, not gimmicks.
- Use a short aside or a direct opinion only when it helps the point.
- Do not add quirks just to sound personal.

## 7. Language handling

- Match the user's language by default.
- Preserve the source language in rewrite mode unless asked to translate.
- Preserve deliberate code-switching in bilingual input.
- Do not force English sentence logic onto Chinese prose.
- Do not force Chinese rhetorical density onto English prose.

## 8. Output defaults

- Keep meaning stable in `rewrite`.
- Do not add unverified facts when details are missing in `generate`.
- Prefer one strong draft over multiple padded options unless the user asks for variants.
