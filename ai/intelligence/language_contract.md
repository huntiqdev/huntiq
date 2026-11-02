# Multilingual Processing Contract

---

## 1. Language Handling Protocols
- Auto-detect primary language from first 5 tokens.
- Supported languages: English (`en`), Amharic (`am`), Afaan Oromoo (`or`), Tigrinya (`ti`).
- Tool arguments emitted to Go are always canonical English slugs (`apartments`, `rent`, `Bole`). Conversational responses to the user are in the detected user language.
