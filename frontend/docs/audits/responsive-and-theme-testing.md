# SellBuy — Responsive Design & Theme System Testing Audit

**Report File:** `docs/audits/responsive-and-theme-testing.md`  
**Date:** 2026-10-22  
**Target Applications:** Next.js 14 / App Router frontend  
**CSS Framework:** Tailwind CSS 3.4 with custom tokens

---

## 1. Viewport Test Matrix & Target Results

| Breakpoint / Device Class | Viewport Range | Core Layout Adaptation | Status |
| :--- | :--- | :--- | :--- |
| **Small Mobile** (iPhone SE, Galaxy A) | 320px – 375px | Single-column listing feed, collapsible search bar, sticky compact header, touch targets ≥ 44px. | **PASS** |
| **Typical Mobile** (iPhone 14/15, Pixel 7) | 390px – 430px | 2-column product grid with compressed card spacing, bottom navigation sheet or full-width search. | **PASS** |
| **Tablet** (iPad Mini, iPad Air) | 768px – 834px | 3-column product grid, integrated category icon carousel, accessible modal dialogs. | **PASS** |
| **Laptop** (MacBook Air, ThinkPad) | 1024px – 1440px | 4-column product grid, persistent left filter sidebar (`FilterSidebar`), top search auto-complete. | **PASS** |
| **Wide Desktop** (24"–27" Monitors) | 1920px+ | 5-column product grid, max container constraint (`max-w-7xl mx-auto`), zero horizontal overflow. | **PASS** |

---

## 2. Navigation & Layout Verification

### A. Header Navigation (`Header.tsx` & `ClientShell.tsx`)
- **Sticky Behavior:** `sticky top-0 z-40 backdrop-blur-md` maintains navigational accessibility without clipping content.
- **Mobile Menu:** Collapses to hamburger on screens `< 768px`; slides out cleanly with background backdrop scrim.
- **Unified Account Dropdown (`UserAccountMenu.tsx`):** On screens `< 640px`, text collapses to user avatar pill with full touch-friendly popover menu.
- **Language Switcher:** Dropdown features outside-click detection (`mousedown` listener with timeout cleanup) preventing stuck menus.

### B. Product Grids & Filter Sidebar (`FilterSidebar.tsx`)
- **Desktop:** Sticky 280px sidebar on desktop (`hidden lg:block w-72 shrink-0`).
- **Mobile Drawer:** Mobile filter toggle button opens full-screen or slide-over drawer on viewports `< 1024px`.
- **Card Resiliency:** Product title clamps at 2 lines (`line-clamp-2`), price formatting handles 8-figure numbers gracefully without clipping or overlapping badges.

---

## 3. Dark Theme & Light Theme Verification

### A. Theme Architecture
- **Root Injection:** In `frontend/src/app/layout.tsx`, inline `<script dangerouslySetInnerHTML={{ __html: themeScript }} />` executes before paint, checking `localStorage.getItem('habesha_theme')` or system `prefers-color-scheme`.
- **Flash of Unstyled Theme (FOUT):** **PREVENTED**. Zero white flashes on dark mode reload.
- **Persistence:** Toggling the theme in `Header.tsx` updates `document.documentElement.classList` and saves `'dark' | 'light'` to `habesha_theme`.

### B. Color Tokens & Contrast Check
- **Dark Mode Background:** `dark:bg-zinc-950` with surface cards at `dark:bg-zinc-900`.
- **Borders & Dividers:** High-contrast subtle borders: `border-zinc-200 dark:border-zinc-800`.
- **Text Readability:** 
  - Primary text: `text-zinc-900 dark:text-zinc-100` (WCAG AAA contrast ratio > 7:1).
  - Secondary/Muted text: `text-zinc-500 dark:text-zinc-400` (WCAG AA contrast ratio > 4.5:1).
  - Accent / Primary CTA: `bg-orange-600 hover:bg-orange-500 text-white` (WCAG AA compliant).
- **Form Inputs:** `bg-zinc-100 dark:bg-white/5 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white` with visible focus rings (`focus:ring-2 focus:ring-orange-500/30`).

---

## 4. Accessibility & Touch Usability Findings

1. **Touch Targets:** All primary action buttons (`Submit`, `Filter`, `Sign In`, `Add Listing`) exceed the 44×44px touch target guideline.
2. **Focus Visibility:** Interactive buttons and links include standard focus rings (`focus:outline-none focus:ring-2`).
3. **Semantic Landmarks:** Pages use `<header>`, `<main>`, `<nav>`, and `<footer>` semantics.
