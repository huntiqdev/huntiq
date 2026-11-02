# HUNTIQ Linguistic Processing Specification

**Languages:** English (`en`), Amharic (`am`), Afaan Oromoo (`or`), Tigrinya (`ti`).

---

## 1. Ethiopic (Ge'ez) Script Normalization Rules
1. **Homophone Standardization:**
   - Unified to first-order Fidel: ሀ, ሐ, ኀ $\to$ ሀ
   - ሰ, ሠ $\to$ ሰ
   - ጸ, ፀ $\to$ ጸ
   - አ, ዐ $\to$ አ
2. **Punctuation Normalization:**
   - Replace Ethiopic full stop (፠, ፡, ፣, ፤, ፥, ፦, ፧, ።) with canonical punctuation tokens.
   - Traditional word separator colon (፡) normalized to whitespace.
3. **Number Systems:**
   - Traditional Ethiopic numerals (፩, ፪, ፫, ፬, ፭, ፮, ፯, ፰, ፱, ፲, ፳, ፴, ፵, ፶, ፷, ፸, ፹, ፺, ፻, ፼) mapped to Arabic integer values.

---

## 2. Afaan Oromoo (Qubee) Normalization Rules
1. **Glottal Stop (Hudhaa) Handling:**
   - Standardize apostrophes (' or ’) to standard ASCII apostrophe `'`.
2. **Vowel and Consonant Lengthening:**
   - Preserve double vowels (aa, ee, ii, oo, uu) and double consonants (bb, dd, gg) essential for semantic distinction.
