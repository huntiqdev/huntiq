# Task State & Slot Filling Contract

---

## 1. Slot Completion Rules
- Category and Listing Type are mandatory for property searches.
- Budget and Location are optional refinement slots.
- Inherit valid slots from previous turn if intent is refinement; flush slots if intent is domain switch.
