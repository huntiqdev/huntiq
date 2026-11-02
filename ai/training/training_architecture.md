# HUNTIQ Model Training Architecture

---

## 1. Modular Model Strategy
HUNTIQ does not train a monolithic foundation model from scratch. We employ a modular, task-specific architecture:
1. **Fast-Path Intent Classifier:** Lightweight distilled transformer (e.g. XLM-RoBERTa-base or MiniLM) fine-tuned on Modules 01, 04, 05, 21.
2. **Named Entity Recognizer:** Token classification model fine-tuned on Modules 06, 07, 09, 10.
3. **Conversational Assistant:** Small multilingual instruction-tuned LLM (e.g. LLaMA-3-8B / Qwen-2.5-7B) optimized via QLoRA for dialogue grounding and tool calling.
