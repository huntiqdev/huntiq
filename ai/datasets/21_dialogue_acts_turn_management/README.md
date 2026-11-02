# HUNTIQ Dataset Module 21 — Dialogue Acts and Turn Management

## 1. Module Overview
- **Module ID:** `21_dialogue_acts_turn_management`
- **Domain:** `conversation_flow`
- **Primary Languages:** English (`en`), Amharic (`am`), Afaan Oromoo (`or`), Tigrinya (`ti`)
- **Format:** JSON Lines (`.jsonl`)
- **Schema:** [`schema.json`](./schema.json)
- **Manifest:** [`manifest.json`](./manifest.json)
- **Sources & Provenance:** [`sources.json`](./sources.json)
- **Data Dictionary:** [`data_dictionary.md`](./data_dictionary.md)

## 2. Purpose & Scope
Classify conversational dialogue acts (greeting, question, request, confirmation, rejection, clarification, thanks, farewell) to maintain conversational coherence.

All examples in this dataset module relate directly to the HUNTIQ marketplace ecosystem in Ethiopia.

## 3. Directory Layout
```text
21_dialogue_acts_turn_management/
├── README.md
├── schema.json
├── manifest.json
├── sources.json
├── data_dictionary.md
├── raw/
│   ├── README.md
│   ├── external/
│   ├── kaggle/
│   ├── huggingface/
│   └── original/
├── processed/
│   ├── normalized.jsonl
│   ├── deduplicated.jsonl
│   └── quality_checked.jsonl
├── training/
│   ├── train.jsonl
│   ├── validation.jsonl
│   └── test.jsonl
├── languages/
│   ├── en.jsonl
│   ├── am.jsonl
│   ├── or.jsonl
│   └── ti.jsonl
├── evaluation/
│   ├── benchmark.jsonl
│   ├── expected_results.jsonl
│   └── edge_cases.jsonl
├── statistics/
│   ├── statistics.json
│   └── quality_report.json
└── examples/
    └── examples.jsonl
```

## 4. Quality & Governance Invariants
- Zero fabricated records.
- Zero customer PII or private database credentials.
- Multi-language split consistency across EN, AM, OR, and TI.
- Strict schema validation enforced before moving from `raw/` to `processed/`.
