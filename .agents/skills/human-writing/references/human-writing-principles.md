# Human writing principles

Use these rules to make prose feel written by a person. They follow the order in which you make decisions: what to say, who will read it, how to structure it, how it should sound, and finally punctuation.

## Make it concrete

### 1. Prefer specificity over abstraction

- Replace broad claims with concrete points.
- Use nouns and verbs that name the thing instead of circling around it.
- Cut lines that only announce importance without adding information.

### 2. Do not use abstract adjectives and adverbs

Abstract modifiers can feel descriptive to the writer, but they give the reader no picture of what happened or what makes something good. The reader has to guess. Watch for four forms:

- **AI-invented shorthand that no real person says**, such as「不繞」「最穩」「一直在撐」「很頂」「不演」. These sound punchy to the model but are not how people talk, and the reader cannot even tell what quality is being claimed. This form is the strongest AI tell of the four. The same family includes bare potential complements posing as capability claims:「XX接得住」「XX補得完」「讀得完的XX」. The construction is normal Chinese when it has a real object（接得住球）; used bare, the model thinks it sounds literary, but the reader never learns what is being caught, completed, or read. When no repair that keeps the word sounds natural, drop the word and say it with an ordinary verb:「接得住」→「流量翻倍也能承受」;「不演」→「不假裝是當場想到的」.
- **Impression summaries**, such as「很好」「很棒」「非常愉快」. They report the writer's feeling about the thing instead of the thing, so the picture in the writer's head never reaches the reader.
- **Ambiguous abstractions**, such as「彈性」「完善」or "robust", "seamless". Every reader decodes them differently, so the writer thinks the point is made while each reader understands something else.
- **Ornate stacked modifiers**, such as「深刻地」「優雅地」「淋漓盡致」. They perform style without adding content. If deleting the word loses no information, it was empty.

Rewrite the sentence in ordinary words. Say what happened, how much, compared with what, what it looks like in use, or what the reader can do. If you do not know a concrete detail, remove the modifier and keep the plain statement. Never invent a number or detail to fill the gap.

The test is whether the reader can picture the scene or carry out the operation, not whether the sentence has a number.「結帳從五步縮成兩步」still fails because the reader does not know those five steps. Concrete writing lets the reader follow the action:「以前在購物網站下單，每次都要重打一遍地址。改版後登入會員，結帳頁直接帶出上次的地址，按送出就完成了」.

The ban targets modifiers that replace content, not all modifiers. When a sentence is already concrete and unambiguous, a modifier can ride on it:「按送出就結束了，順很多」works because the reader has just seen the operation and knows exactly what "順" refers to. Two limits still apply: do not reuse the same modifier across the text, and do not attach a modifier to every sentence. When most sentences stand plain, the few modifiers that remain actually register.

### 3. Cut back on「V得C的N」slogan phrases

- Phrases like「讀得完的報告」turn a potential complement into an attributive because the model thinks it sounds literary. Used as a heading or selling point, the phrase claims a quality without saying what makes it true, and repeating the pattern turns prose into slogans.
- The construction is normal Chinese when the surrounding text says how the claim holds. Keep at most one, backed by specifics, and never stack several in a row.
- Otherwise rewrite it as a full sentence with a subject and verb: state what the reader can actually do, and what makes that possible（「讀得完的報告」→「報告壓在兩頁，開會前讀得完」）.

### 4. Give every action and object an anchor

- Every verb needs a stated object, and every scene needs a stated setting.「一頁一頁按」— pressing what, on which screen?「登入」— logging into which site or app? If the reader has to guess what an action lands on or where it takes place, the sentence is not concrete yet, no matter how plain it sounds.
- Establish the setting once, early:「在購物網站買東西」. After that, later actions inherit the setting and do not need to re-specify it. Anchoring is done at the scene level, not by re-explaining every sentence.
- Stop condition: the text is concrete enough when the intended reader would no longer ask「哪個？」「按什麼？」about anything that matters to the point. Specificity serves the reader's purpose. Do not drill down forever or pad the text with detail the point does not need.

## Write for a reader with zero context

### 5. Keep the subject and the verb recoverable

- Chinese drops subjects naturally when the topic chain carries them:「改好了，等你看過確認」loses nothing. The limit of ellipsis is that the reader can still say who acts and what the action is.
- When one clause supports two readings of who does what（「給你驗」: the writer verifying, or the reader checking and accepting?）, the ellipsis went too deep. Write the actor back in.
- The verb must remain a word you can point at. A directional（進）, an aspect particle, or a bare object cannot carry the action; when the verb has vanished, put it back（一律進例集→一律寫進例集）.
- Clipped verbs (rule 10) and obligations without a subject (rule 13) are this rule caught in two specific forms.

### 6. Assume zero shared context

- Write as if the reader has none of the Agent's prior conversation, planning, reasoning, project background, or unstated domain knowledge.
- Do not use「如上」「這一層」「前述問題」or similar references unless the current text has clearly introduced them.
- Make the reader's understanding depend on what is written, not on what the writer remembers thinking.

### 7. Keep internal reasoning out of the final prose

- Before drafting, distinguish instructions addressed to the AI, background supplied to guide the writing, and information intended for the reader. Apply writing instructions through content selection and expression. Include background only when the reader needs it to understand or act; being present in the prompt is not enough reason to put it in the prose. Communication goals such as encouraging or reassuring someone should shape the whole message; they do not automatically need a sentence naming the desired response or disproportionate emphasis.
- When the user says "do not make the reader think X," convey the intended understanding through concrete facts, purposes, or actions. The instruction alone does not establish that the reader has raised X or needs an explicit denial. Keep a denial when the reader's actual question or task requires it. See the English examples under example 16 in [human-writing-examples.md](./human-writing-examples.md).
- Do not assume the reader knows the categories, layers, sequence, or terminology created during planning.
- Translate internal frameworks into concrete actions, observable facts, order, reasons, or consequences.
- Introduce a concept before using it as a shared reference. If it cannot help the reader understand or act, remove it.

### 8. Do not over-explain

- Explain what the reader needs to understand the point, then stop.
- Make prose concise by removing repetition and unnecessary wording while preserving the details needed to understand or act. A complete verb can still leave its meaning unfinished: say what must be marked, checked, or explained, including the content or criterion when the surrounding text does not make it clear. Restore necessary information even when the sentence becomes longer. See example 18B in [human-writing-examples.md](./human-writing-examples.md).
- Do not add an answer to an objection, accusation, or inner judgment that no one raised.
- If a sentence only defends against an imagined criticism and removing it does not weaken the point, remove it.
- The same reflex appears as a limitation or disclaimer appended to every claim. In genres whose job is to persuade — a proposal, a pitch — stressing limitations no one asked about defends the writer instead of informing the reader. Put risks where the genre expects them, once, in their own section, not on the tail of every sentence.

### 8A. Write the current document for its current reader

- Before adding or changing content in an existing document, rule set, prompt, or skill, read the relevant sections and linked references. Search for the same subject, its exceptions, and examples, then identify which file and section should own the information.
- Update or merge the existing passage when it already covers the subject. Place genuinely new content where the reader needs it in the existing structure, not automatically at the end or in the section most recently edited. Preserve distinct conditions and authorized scope when consolidating.
- Keep the full rule or explanation in one appropriate place. Other sections should link to it or give only the reminder needed at that step. Put operational detail in the existing reference or tool help when that is its established home.
- After editing, reread the affected section with its neighboring paragraphs and relevant references. Remove duplicate or superseded wording, resolve contradictions, and update affected examples, headings, and links. Reorganize only as much as the requested change needs.
- When revising rules or documents, state the adopted approach and its conditions. If A has been replaced by B, and A is absent from the current document and irrelevant to the reader's task, write B on its own.
- Keep a contrast or prohibition when the reader needs it to distinguish live alternatives, avoid a likely consequential mistake, or follow a migration. Revision history belongs in a requested change log or comparison.
- Review each negative qualification against the finished document: does it help the reader act, or merely echo an editing discussion, rejected approach, or the writer's self-defense? Remove the latter. See example 16A in [human-writing-examples.md](./human-writing-examples.md).

### 9. Avoid unexplained niche vocabulary

- Prefer common words over niche terms, jargon, project shorthand, and author-created labels.
- If a term such as「剪刀差」is necessary, explain it in plain language at first use. Otherwise, directly describe the difference or consequence it refers to.
- Do not treat a familiar term within one field as familiar to every reader.
- Do not use unfamiliar imagery or literal translated terms just because they sound compact or professional. Words such as「基線」and「閉環」make readers stop and work out what the sentence means. Unless the intended reader already uses the term and its technical precision matters, replace it with the people, actions, information, and unresolved consequence it hides.
- In ordinary prose, avoid attaching「邊界」to an abstract noun. Phrases such as「事實邊界」「能力邊界」「認知邊界」「風險邊界」and「責任邊界」usually hide a plainer idea: the applicable cases, a limit, who handles what, or the conditions under which something changes. Write that idea directly. Keep「邊界」only for an actual dividing line or an established technical, legal, or organizational term the intended reader already uses; do not invent a new「XX 邊界」for one passage.
- Use「說法」only when the noun points to words someone actually said, wrote, or publicly expressed. A quotation, paraphrased claim, disputed wording, or stated public position can be called a「說法」. When the sentence describes an action or result instead, name that action or result directly（「避免錯誤範圍成為正式說法」→「錯誤的影響範圍沒有發布出去」）. Adding「正式」「普遍」or「常見」does not make「說法」valid when no act of expression is being discussed.

### 10. Do not clip a two-syllable verb down to one character

- A one-character Chinese verb needs something behind it to stand on: an object（查資料）, a complement（改好了、查清楚）, or an aspect marker（測過了、拆掉）. Left bare at the end of a clause, it stops sounding like anything a person would say.「範例照抄也要驗」「命中的已拆」「表現優異的項目也要查」「改好再給你驗」are all the same failure.
- Clipping the verb usually drops its object with it, so putting the second character back is not enough.「範例照抄也要驗一次」reads more smoothly, but the reader still cannot tell what to verify. Write「就算照抄範例，也要自己檢查每一條有沒有放錯欄位」. This is rule 4 from another angle: a verb worn down to one character is usually the sign that its object went missing.
- Do not use bare「收」as a vague editing or completion action when the sentence does not say what is being gathered, shortened, revised, or finished. Adding「一次」does not clarify the action. Name the intended work in ordinary Chinese; choose the wording from context rather than prescribing one replacement. Keep established uses with a clear object, such as「收衣服」or「收回意見」.
- Status columns in a table are the exception. The row subject already supplies the object, so it does not need repeating, but the verb still has to be whole: write「驗過」, not「已驗」.
- Column width is not a reason to shave characters. When something has to be shorter, choose a shorter complete wording instead.
- The same applies to nouns: do not coin an abbreviation on the spot. Use only ones already in circulation（健保、台大、API）, and drop any that collides with an existing word —「產品經理」cut down to「產經」reads as 產業經濟.
- In English the matching habit is inventing an acronym mid-paragraph, or using a verb as a noun ("a solve", "the ask", "learnings"). Restore the ordinary wording.

### 10A. Keep the ordinary full form of Chinese

- A sentence can be grammatically complete and still sound over-compressed. In ordinary prose, prefer the form people normally say（列出事實、混淆初稿與定稿、會不會寫文章）instead of shortening every action（列事實、把兩種狀態混在一起、會不會寫）. Short forms can stay when the object is already unmistakable and the shorter wording is conventional or creates a deliberate rhythm（列三點、他會寫，也會畫）.
- Check neighboring clauses together. Two bare or one-character verbs in close succession often make prose sound clipped. Restore the missing object, complement, or full verb unless the abrupt cadence serves a clear purpose in that passage. Also remove accidental verb echoes whose functions overlap（「主管需要要求逐項查證」→「主管應要求逐項查證」）.
- Do not delete function words merely to make a sentence shorter.「的」can clarify the relation between nouns and give the second half of a contrast enough weight（「付出的卻是文件的可信度」）;「會」can mark an expected or conditional result（「查證納入流程後，錯誤才會下降」）. Keep the shorter form when it is the established compound, a direct state, or an observed event（文件可信度評估、有憑證才算完成、改版後錯誤才下降）.
- Do not use「把 A 說成 B」as a generic shortcut for every distortion of scope or meaning. It is natural when someone literally says or mistakes one thing for another（把週三說成週四）. For an abstract consequence, state what the draft claims or what the reader may misunderstand（「成稿可能把有限情境說成全面故障」→「讀者可能會以為所有情況都會故障」）.
- Do not make an abstract noun perform a physical return.「讓討論回到正題」is an established relation;「讓每句話回到來源」is not. Name the checkable relation instead（「每句話都能對照原始資料」）.

### 10B. State effects and conditions in natural Chinese

- Use「仍」when the context conveys that something continues or holds despite a change or contrary expectation, carrying the sense of「還是」. Ordinary requirements and factual statements need no implied contrast. Choose wording from the actual relationship rather than adding「仍」for emphasis or formality.

- Avoid repeatedly using「代表」or「代表……不代表……」as a generic frame for explaining relationships. State the actual meaning, effect, inference, or contrast in wording that fits the context. Keep the word where it is natural and precise; swapping it for another repeated formula does not solve the problem.
- Make the relationship between neighboring clauses clear, including whether a result is observed, possible, or intended. Choose the sentence structure from the meaning and context rather than applying a fixed connector or replacement phrase.
- Write conditions and qualifications in idiomatic Chinese rather than copying English clause structure. Sentence-initial「只是」often sounds translated when used to introduce a situation, including wording modeled on English “just because.” Reconsider the whole sentence and its connection to the preceding text. It can remain where the context makes it natural, such as a qualification of the preceding statement.
- Check ordinary word combinations as well as grammar. Avoid awkward repetition caused by mixing literal nouns with idioms. Read the whole clause in context and choose a natural combination that preserves the intended action.

### 11. Remove unexplained figurative language

- Keep a metaphor only when the reader can tell what it refers to and what it means in the situation.
- Replace vague images such as「漏過去」with the concrete missing data, action, location, or consequence.

## Structure the argument straight

### 12. Lead with the point in informational writing

- When the text informs — announcements, emails, reports, documentation, status updates — put the decision, change, or required action in the first sentence, and the reasons and background after it. Readers of informational text scan for what changed and what to do; making them read the build-up first is author-centric ordering.
-「鑑於近期內部作業流程有所調整，經多方評估與審慎考量後，本公司決定自下月一日起，將週會調整至週三下午辦理」buries the only fact the reader needs. Write「下個月起，週會改到週三下午開（原本是週一上午）」, or in a formal register,「本年度健康檢查 9 月 1 日至 15 日開放報名，名額增至三百名」with the reasons following. Leading with the point changes the order, not the register.
- Headings state content, not category.「5/20 凌晨 2–4 點系統暫停服務」tells the reader what is happening;「系統維護公告」only names a drawer to look in. A heading should answer, not tease.
- Voiced writing may earn a slower opening, but even there, never bury the fact the reader came for.

### 13. Write instructions as doable actions

- When the reader must do something, the sentence needs an actor, an action with its object, and a deadline.「相關資料應於期限內補齊，逾期恕不受理」names no one, no document, and no date; write「申請人請於 8 月 29 日前，將身分證影本上傳至報名系統」.
- Every obligation needs a subject.「本案應予處理」— by whom?
- Put the condition before the action（「如果 X，就做 Y」）so readers it excludes can stop reading early.
- When the text expects the reader to act, end with the next step.「若有疑問歡迎隨時聯繫」closes nothing; write「請回信告訴我選方案一還是方案二，我週五前要回覆廠商」.

### 14. Fixed-format genres keep their conventions

- Some document types carry a fixed format and conventional phrases: 台灣公文 with 主旨／說明 and 期望語（請查照、請核示）, contracts and legal clauses with their set formulas. In these genres the conventions are how a person writes. Stripping them to sound plainer makes the text less human, not more.
- Apply only the structural rules there — the conclusion goes into 主旨, details move down into 說明, actors and deadlines stay explicit — and leave the register, formulaic phrases, and required vocabulary alone.「主旨：請轉知所屬同仁報名本校 9 月 10 日及 17 日辦理之資訊安全教育訓練，請查照。」keeps every convention and still leads with the point.
- Vocabulary cleanup works the same way.「係」「惟」「俾」are padding in an email and normal register in a contract. Judge by the genre the reader expects, not by a universal word list.

### 15. Do not default to setup-and-reveal structures

- 先抑後揚 and 先破後立 start by building up a wrong view, then reveal the real point. Models fall into this pattern easily. When every section starts that way, readers learn to skip the first half of each paragraph.
- Most points are stronger stated directly: lead with the claim, then support it.
- Keep the structure only when the intended reader genuinely holds the misconception and correcting it is the point of the passage. Never invent a strawman just to have something to knock down. Even when justified, use it once in a piece, not as the skeleton of every section.

### 16. Make the point clear

- The reader should be able to tell what the writer thinks or wants them to understand.
- In grounded writing, state the judgment clearly and keep the framing consistent.
- In voiced writing, use opinion, tension, or a light edge when it helps the point.

### 17. Write body text as body text, not as slogans

- Body sentences need a subject, a verb, and a connection to the sentences before and after. A paragraph built from headline-style declarations, such as「效率，從此不同」「一個平台，解決所有問題」, reads like presentation titles pasted together, and the reader cannot find the line of reasoning.
- Headline diction belongs in headings. If a sentence would work unchanged as a slide title, it is probably not doing body-text work.
- State the claim as a full sentence with its reason or consequence attached, then continue the thread into the next sentence.

### 18. One sentence takes one step

- The model plans the whole passage at once, so it tries to fit the claim, its parallel aspects, conditions, and exceptions into one sentence, strung together with「同時」「並」「也」「兼顧」. The reader can only walk one path at a time; a sentence that covers everything reads mechanical and cold.
- Let a sentence be temporarily incomplete: state this step plainly and end it. The caveat and the exception are the next footprint. What comes next should be seen on arrival, not announced at the start.
- When one sentence strings two or more separate aspects together with these connectors, split it. Real enumerations may stay lists; this rule governs the pacing of an argument, not items.
- This does not overturn leading with the point: in informational writing the conclusion still comes first. Walking governs how every sentence after it advances.

### 19. Break sterile symmetry

- Do not force every idea into neat triples.
- Vary paragraph length, and let sentence length change with the point.
- When every paragraph follows the same arc, break the arc.

## Control rhythm and tone

### 20. Vary rhythm on purpose

- Mix short sentences with longer ones.
- Avoid paragraphs where every sentence has the same shape.
- Let emphasis come from contrast, not from constant intensity.

### 21. Choose stronger verbs before adding adjectives

- Fix weak verb choice first.
- Prefer "cut", "argue", "miss", "ship", "delay", "expose", or "earn" over soft general verbs.
- Use adjectives only when they add meaning, not sheen.

### 22. Carry the feeling of the moment in the form

- A person thinks while writing, so every sentence carries how the writer feels at this point. The feeling rarely shows up as adjectives; it soaks into the form: the choice of function words and frames（「我以為」already carries an expectation falling through）, the long or short variant of the same word（但／但是、沒辦法／無法）, how tight the punctuation is, and where the words pile up: linger where you care, move fast through routine. Readers cannot always name these signals, but they read them.
- Take the feeling from the content itself: a genuinely messy problem may be written with its weight, a genuinely clean fix may be written light. Do not fabricate personal experience, do not stage spontaneity, and do not hang an emotion word on every sentence. Emotion spread evenly reads as no emotion at all.
- In grounded writing this stays an undertone; voiced writing may let it surface.

### 23. Keep it restrained

- The best human-like line is often simpler than the first draft.
- Do not show off range when a plain sentence will land better.
- Leave room for the reader to trust the voice.

### 24. Do not overcorrect into telegraph style

- Removing AI patterns is a means; the goal is text a person would actually send. Do not strip the courtesy, warmth, or connective tissue the genre normally carries. The fix for an empty polite letter is adding the missing information, not deleting the manners.
- When a cleanup pass leaves nothing but clipped statements of fact, the text has become machine-like from the opposite direction. Reread the result as its genre: a reply, a post, an announcement. If no person would send it that way, put the human parts back.

### 25. Match social context

- Business writing earns trust by stating facts and next steps plainly, without template phrases.
- In appreciation or adoption messages, describe how the recipient's work helps you. Keep the relationship between writer and reader intact; a declaration that the work meets your standards can turn thanks into an unsolicited evaluation. Use evaluative language when assessment is the actual purpose. See example 27A in [human-writing-examples.md](./human-writing-examples.md).
- Public-facing copy should read as if someone chose each word for this product and this reader.
- In commentary, take a position and give its reason. A survey of safe opinions reads as filler.
- When replying to a message, mirror the other side's paralinguistic level: if they write with「～」, exclamation marks, or emoji, do not answer in bare full-stop sentences, which read as cold or even displeased; if they write formally, do not sprinkle symbols. A friendly signal that gets no echo erodes trust.
- Chinese and English should each sound native to their own rhetorical habits.

## Punctuation

### 26. Treat semicolons as an exception

- In ordinary Chinese prose, do not use semicolons by default.
- During the final pass, replace each semicolon with a full stop, comma, colon, or a clearer sentence boundary whenever possible.
- Retain one only for formal quotations, required syntax, or a genuinely complex list where other punctuation would be less clear.
