# HUNTIQ Dataset Dependency Map

**Project:** HUNTIQ Marketplace AI Platform  
**Architecture:** 30-Module Coordinated Conversational Intelligence Pipeline  
**Version:** 1.5.0 (Audited Architecture)  

---

## 1. Pipeline Execution Dependency Graph

The 30 HUNTIQ dataset modules operate as a single coordinated intelligence pipeline, organized into 6 distinct execution tiers. Upstream modules produce structured signals consumed by downstream reasoning and tool execution layers.

```text
[User Utterance]
      │
      ▼
Tier 1: Ingestion & Safety (M01 Greetings, M02 Identity, M14 Multilingual, M15 Normalizer, M18 Safety)
      │
      ▼
Tier 2: Routing & Classification (M04 Domain, M05 Intent, M06 Category, M11 Multi-Intent, M21 Dialogue Acts)
      │
      ▼
Tier 3: Extraction & Constraints (M07 Entities, M08 Search Query, M09 Price, M10 Location, M17 Math, M23 Constraints)
      │
      ▼
Tier 4: Dialogue State & Tracking (M12 Disambiguation, M13 Followups, M22 Task State, M24 Prefs, M25 Repair, M26 Reset, M27 Filtering)
      │
      ▼
Tier 5: Go Tool Execution (M03 Platform FAQ, M16 Ranking, M19 Tool Calling, M28 Uncertainty, M29 Error Recovery)
      │
      ▼
Tier 6: Grounded Response NLG (M20 Grounding Quality, M30 Natural Language Generation)
```

---

## 2. Explicit Module Dependency Table

| Module ID | Module Name | Upstream Dependencies | Downstream Consumers | Runtime Requirement |
| :--- | :--- | :--- | :--- | :--- |
| **01** | General Greetings | 14, 15 | 04, 21, 30 | Python Fast Path |
| **02** | Identity & Behavior | 01, 14 | 04, 30 | Python Fast Path |
| **03** | Platform Knowledge | 04, 05 | 20, 30 | Retrieval / Go Authority |
| **04** | Domain Classification | 14, 15, 18 | 05, 06, 12, 19 | Python Routing |
| **05** | Intent Classification | 04, 14, 21 | 06, 07, 08, 11, 19 | Python Routing |
| **06** | Category Understanding | 04, 05, 14 | 07, 08, 12, 19, 22 | PostgreSQL Taxonomy Sync |
| **07** | Entity & Attribute Extraction | 05, 06, 14, 15 | 08, 19, 22, 23 | Python NER Extraction |
| **08** | Search Query Understanding | 05, 06, 07, 09, 10 | 19, 22, 27 | Query Parser |
| **09** | Price & Budget Understanding | 07, 14, 15 | 08, 17, 23, 27 | Deterministic Numeric Parser |
| **10** | Location Understanding | 07, 14, 15 | 08, 19, 22, 27 | Geo Authority (Ethiopia DB) |
| **11** | Multi-Intent Understanding | 05, 21 | 08, 19, 22 | Multi-turn Intent Queue |
| **12** | Clarification & Ambiguity | 06, 07, 22, 23 | 21, 30 | Disambiguation Generator |
| **13** | Conversation Follow-ups | 21, 22, 26 | 08, 22, 27 | Session Context Manager |
| **14** | Multilingual Understanding | None (Base Layer) | 01–13, 15–30 | Normalizer / Tokenizer |
| **15** | Spelling & Noisy Input | None (Base Layer) | 01–14 | Ge'ez & Latin Normalizer |
| **16** | Recommendation & Ranking | 08, 19, 24 | 20, 30 | PostgreSQL Result Ranker |
| **17** | Comparison & Calculations | 07, 09 | 20, 30 | Deterministic Python Math |
| **18** | Safety & Platform Policy | 14, 15 | 04, 19, 20, 30 | Security Guardrail |
| **19** | Tool Calling & Structured Outputs | 05, 06, 07, 08, 09, 10, 18 | Go Backend Client | Go Authorized Gateway |
| **20** | Response Grounding & Quality | 19, Go Search Results | 30 | Anti-Hallucination Verifier |
| **21** | Dialogue Acts & Turn Management | 01, 05, 13 | 05, 12, 22, 25 | Conversation State Machine |
| **22** | Task State & Slot Completion | 06, 07, 08, 09, 10, 13 | 12, 19, 23, 27 | Session Slot Matrix |
| **23** | Constraint Reasoning | 08, 09, 10, 22 | 12, 19, 27 | Constraint Satisfier |
| **24** | User Preference Understanding | 06, 07, 10, 22 | 16, 27 | Ephemeral / Profile Cache |
| **25** | Negative Feedback Recovery | 21 | 12, 22, 26, 30 | Repair & Apology Handler |
| **26** | Topic Switching & Reset | 05, 21, 22 | 22, 27 | Context Slot Flusher |
| **27** | Search Refinement & Filtering | 08, 13, 22, 23 | 19, 22 | Progressive Query Builder |
| **28** | Confidence & Uncertainty | 04, 05, 07, 19 | 12, 29, 30 | Calibrated Uncertainty Gate |
| **29** | Service Failure & Recovery | Backend Circuit Breaker | 30 | Go Degraded Fallback Mode |
| **30** | Conversation Quality & NLG | 01, 02, 03, 12, 20, 25, 29 | User Client | Localized Response Engine |
