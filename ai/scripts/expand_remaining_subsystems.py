#!/usr/bin/env python3
"""
Expand remaining short documents in coordination, governance, and knowledge subsystems.
"""

import os


def expand(path, title, scope, rules, protocols):
    content = f"""# {title}

**Project:** HUNTIQ Universal Marketplace AI Platform  
**Document Classification:** Technical Engineering Specification  
**Status:** Approved Architecture Standard  

---

## 1. Scope & System Objective
{scope}

---

## 2. Detailed Technical & Operational Rules
{rules}

---

## 3. Enforcement & Verification Protocols
{protocols}

---

## 4. Cross-System References
- Sovereign Go Authority: backend/internal/
- Python AI Inference Runtime: ai/app/
- Authoritative Marketplace Database: PostgreSQL 16
"""
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w") as f:
        f.write(content.strip() + "\n")
    print(f"Expanded: {path}")


# Coordination docs
coord = [
    (
        "ai/dataset_coordination/conflict_detection.md",
        "Cross-Dataset Conflict Detection Protocol",
        "Automated continuous integration checks to detect conflicting label enums, contradictory intent classifications, and overlapping slot definitions across the 30 dataset modules.",
        "- Canonical Label Dictionary: All modules must draw category slugs strictly from the PostgreSQL taxonomy mirror.\n- Currency Uniformity: Currency codes must be strictly 'ETB' across all dataset modules.\n- Incompatible Constraint Screening: Detection of schemas that define mutually exclusive validation rules.",
        "Executed via automated schema comparison scripts running during CI validation.",
    ),
    (
        "ai/dataset_coordination/cross_dataset_consistency.md",
        "Cross-Dataset Consistency & Ontology Alignment",
        "Establishes a single unified ontology across all 30 conversational, classification, and extraction datasets.",
        "- Shared Entity Spans: A brand name extracted in Module 07 must match the brand filter compiled in Module 19.\n- Dialogue Act Alignment: Dialogue act tags in Module 21 must seamlessly map to state machine transitions in Module 22.\n- Multilingual Parity: Semantic equivalence guarantees that Amharic, Oromo, Tigrinya, and English splits share identical label spaces.",
        "Verified using automated cross-dataset linting tools.",
    ),
    (
        "ai/dataset_coordination/synchronization_policy.md",
        "Marketplace Schema Synchronization Policy",
        "Defines protocols for updating AI dataset schemas whenever Go backend marketplace categories, pricing policies, or geocoding hierarchies evolve.",
        "- Automated Code Generation: Go struct definitions in internal/models generate matching JSON schemas for Python.\n- Synchronous Version Bump: Any addition of a top-level category triggers a minor version bump across all 30 manifests.\n- Backward Compatibility: Models must support previous category slugs for 60 days following a taxonomy migration.",
        "Managed through Git release hooks and database migration listeners.",
    ),
    (
        "ai/dataset_coordination/shared_category_reference.md",
        "Shared Marketplace Category Reference",
        "Authoritative dictionary linking live PostgreSQL category IDs with AI dataset category enums across all 8 primary verticals.",
        "- 8 Primary Verticals: Products & Electronics, Real Estate & Property, Vehicles, Jobs, Agriculture, Services, Hotels, Event Tickets.\n- 45 Subcategories: Full hierarchical tree mapped to canonical slugs.\n- Synonyms and Aliases: Local Ethiopian naming variations mapped to canonical slugs.",
        "Mirrored directly from backend/internal/models/category.go.",
    ),
    (
        "ai/dataset_coordination/semantic_alignment.md",
        "Cross-Lingual Semantic Alignment Guidelines",
        "Rules for aligning parallel semantic concepts across English, Amharic, Afaan Oromoo, and Tigrinya.",
        "- Aligned Semantic IDs: Parallel examples across language split files share identical root identifiers.\n- Cultural Nuance Preservation: Respect regional trade practices without distorting the underlying intent.\n- Script Independence: Semantic classification models must process Latin and Ethiopic scripts with equal fidelity.",
        "Monitored via cross-lingual cosine similarity benchmarks on sentence embeddings.",
    ),
    (
        "ai/dataset_coordination/duplicate_management.md",
        "Cross-Dataset Duplicate Management & Partitioning",
        "Guarantees that training examples are uniquely allocated and not redundantly duplicated across multiple module datasets.",
        "- Hash Partitioning: Exact utterance hashes are indexed centrally in dataset_registry.json.\n- Multi-Task Exception Rule: If an utterance is intentionally used for both intent and entity training, it must be explicitly tagged as multi_task.\n- Zero Leakage Guarantee: Ensures no evaluation example in one module appears in the training split of another module.",
        "Scanned via automated duplicate detection pipelines in ai/data_pipeline/.",
    ),
    (
        "ai/dataset_coordination/module_interaction_matrix.md",
        "Module Interaction & Signal Routing Matrix",
        "Comprehensive routing table detailing the producer-consumer relationships between all 30 intelligence modules.",
        "- Upstream Producers: Ingestion, normalizer, and intent classifiers emit typed signals.\n- Downstream Consumers: Slot completion, tool compiler, and grounding verifiers consume signals.\n- Fallback Escalation: Defines immediate fallback paths when an intermediate module produces low-confidence output.",
        "Visualized in ai/datasets/dataset_dependency_map.md.",
    ),
]

# Governance docs
gov = [
    (
        "ai/governance/auditability.md",
        "AI Auditability, Provenance & Traceability Framework",
        "Maintains an immutable, legally defensible audit trail of all dataset acquisitions, annotations, and model updates.",
        "- Git Commit Lineage: Every dataset release is permanently tagged with a cryptographic Git commit hash.\n- Cryptographic Hash Chains: Raw batches must have SHA-256 hashes recorded in sources.json.\n- Human Decision Logging: Model approval and rejection decisions are logged with reviewer credentials.",
        "Audited quarterly by internal engineering and legal compliance teams.",
    ),
    (
        "ai/governance/external_dataset_licensing.md",
        "External Dataset Licensing & Intellectual Property Policy",
        "Governs the acquisition, validation, and commercial use of third-party datasets.",
        "- Prohibited Licenses: Non-Commercial (NC) and ShareAlike (SA) licenses that risk proprietary code infectivity.\n- Permitted Licenses: CC-BY 4.0, MIT, Apache 2.0, Open Data Commons, and proprietary commissioned datasets.\n- Attribution Manifest: All third-party sources are cataloged in ai/datasets/dataset_source_catalog.json.",
        "Legal review sign-off required before merging any external dataset.",
    ),
    (
        "ai/governance/prompt_injection_policy.md",
        "Prompt Injection & Adversarial Jailbreak Defense Policy",
        "Engineering standards for defending against direct and indirect prompt injections in marketplace conversations.",
        "- Strict XML Delimiters: Untrusted user input and database listing content are enclosed in passive XML tags.\n- Passive Data Invariant: Database content is strictly treated as passive data, never executable instructions.\n- Dual-Model Verification: High-risk requests pass through a secondary safety classification filter.",
        "Benchmarked using 1,000 adversarial injection cases in Module 18.",
    ),
    (
        "ai/governance/model_security.md",
        "Model Supply Chain Security & Artifact Integrity",
        "Security protocols for downloading, verifying, and storing model weights and runtime dependencies.",
        "- Safe Serialization Formats: Only Safetensors and GGUF formats are permitted; raw Python pickle files are strictly prohibited.\n- Private Network Isolation: Production AI service runs on 127.0.0.1:8000 with zero public internet ingress.\n- Secret Authentication: Header X-HUNTIQ-AI-SECRET required on all incoming HTTP requests.",
        "Automated vulnerability scanning with Bandit and pip-audit.",
    ),
    (
        "ai/governance/incident_response.md",
        "Data & AI Security Incident Response Plan",
        "Rapid containment and remediation protocols for data leaks, hallucination emergencies, or adversarial exploits.",
        "- Severity 1 (Immediate Containment): Data leak or safety filter bypass triggers automatic Go circuit breaker.\n- Severity 2 (Urgent Patch): High hallucination rate on prices triggers model rollback within 1 hour.\n- Post-Mortem Requirement: Comprehensive root-cause analysis completed within 48 hours.",
        "Drills conducted bi-annually by systems and security teams.",
    ),
    (
        "ai/governance/anonymization.md",
        "Data Anonymization & Cryptographic Masking Protocol",
        "Protocols for identifying and stripping personally identifiable information (PII) from marketplace training records.",
        "- Phone Number Pseudonymization: Ethiopian mobile numbers (+2519... / 09...) replaced with deterministic synthetic placeholders.\n- National ID Masking: Fayda numbers, Kebele IDs, and passport details stripped unconditionally.\n- Location Sanitization: Specific house and door numbers generalized to subcity neighborhood landmarks.",
        "Automated PII scrubbing scripts executed during normalization workflow.",
    ),
    (
        "ai/governance/commercial_use_approval.md",
        "Commercial Use & Monetization Clearance Policy",
        "Ensures that all AI training datasets, model weights, and inference pipelines are 100% cleared for commercial enterprise operations.",
        "- Commercial Freedom: Guarantees HUNTIQ can charge for Gold merchant subscriptions, featured listing promotions, and merchant APIs without licensing encumbrance.\n- Indemnification Verification: External data providers must warrant that data was legally gathered without scraper terms violations.",
        "Formal sign-off documented in ai/datasets/dataset_license_registry.json.",
    ),
    (
        "ai/governance/dataset_access_control.md",
        "Role-Based Access Control (RBAC) for AI Data Platform",
        "Restricts access to raw training data, model checkpoints, and annotation pipelines based on verified job role.",
        "- Annotators: Restricted to labeling interfaces without direct production database access.\n- Data Engineers: Read/write access to ai/data_pipeline/ and ai/datasets/.\n- Production Services: Read-only access to approved model weights on 127.0.0.1:8000.",
        "Enforced via Linux filesystem permissions and Git branch protection rules.",
    ),
    (
        "ai/governance/data_retention.md",
        "Data Retention, Archival & Destruction Policy",
        "Defines mandatory retention schedules and secure data destruction methods.",
        "- Guest User History: Purged from volatile memory within 30 minutes of session inactivity.\n- Anonymized Training Splits: Retained for the active production lifespan of the corresponding model version.\n- Deprecated Datasets: Archived in cold storage for 12 months, then securely wiped.",
        "Automated cron jobs in ai/data_pipeline/ enforce scheduled purges.",
    ),
    (
        "ai/governance/privacy_by_design.md",
        "Privacy-by-Design Architecture & Data Minimization",
        "Technical architecture ensuring user privacy is preserved by default at every layer of the AI system.",
        "- Data Minimization: The AI service never logs complete conversation transcripts containing PII.\n- Ephemeral Session Caching: In-memory context cache stores only active slot values, not conversational logs.\n- Local Processing: All AI inference runs on local hardware in Ethiopia without transmitting data to overseas public LLM APIs.",
        "Audited against Ethiopian Data Protection Proclamation guidelines.",
    ),
    (
        "ai/governance/data_poisoning_policy.md",
        "Data Poisoning & Supply Chain Integrity Defense",
        "Defends the dataset ecosystem against intentional corruption, malicious label flipping, and backdoor attacks.",
        "- Hash Validation: Every batch imported must match its registered SHA-256 hash.\n- Statistical Outlier Detection: Scans incoming records for anomalous token distributions or label concentrations.\n- Quarantined Staging: Raw data cannot enter training splits without passing quality control scripts.",
        "Enforced in ai/data_pipeline/quality_control_workflow.md.",
    ),
    (
        "ai/governance/customer_data_policy.md",
        "Customer Data Handling & Platform Consent Policy",
        "Outlines the terms under which registered buyer and seller data is processed for marketplace optimization.",
        "- Transparent Terms: Users are clearly informed that natural language search queries are processed by the HUNTIQ AI assistant.\n- Opt-Out Mechanism: Users may request their search history not be utilized in future training corpora.\n- No Sensitive Data Collection: Biometric, religious, or political data is never collected or inferred.",
        "Documented in official user-facing Terms of Service.",
    ),
    (
        "ai/governance/consent_and_legal_basis.md",
        "Consent & Legal Basis Architecture",
        "Establishes legal grounds for processing marketplace data under Ethiopian civil and commercial law.",
        "- Contractual Necessity: Processing search queries is required to fulfill the marketplace search service.\n- Legitimate Interest: Fraud prevention and platform safety monitoring.\n- Explicit Consent: Merchant opt-in for advanced AI inventory analytics.",
        "Reviewed by HUNTIQ legal counsel.",
    ),
    (
        "ai/governance/personal_data_policy.md",
        "Personal Data Protection & User Rights Architecture",
        "Technical implementation of user data rights including access, rectification, and erasure.",
        "- Right to Access: Users can download an export of their listings and profile.\n- Right to Erasure: Account deletion purges all active listings and personal data across PostgreSQL and AI caches within 24 hours.\n- Zero Re-identification: Anonymized dataset records cannot be re-identified.",
        "Automated API endpoints in Go backend handle erasure requests.",
    ),
    (
        "ai/governance/guest_conversation_policy.md",
        "Guest Conversation Handling & Ephemeral Privacy Policy",
        "Defines privacy safeguards for unauthenticated visitors browsing HUNTIQ.",
        "- Zero Account Requirement: Guests can search and converse with the AI freely without logging in.\n- Session Isolation: Each guest receives a random ephemeral session ID stored only in session cookies.\n- Memory Flushed: Guest context is completely erased from Redis/memory after 30 minutes.",
        "Enforced in ai/app/core/context_manager.py.",
    ),
]

# Knowledge docs
know = [
    (
        "ai/knowledge/retrieval_policy.md",
        "Knowledge Retrieval & Grounded Citation Policy",
        "Specifies how the AI assistant queries, validates, and cites official platform knowledge documents.",
        "- Official Authority: Knowledge citations must reference verified documents in ai/knowledge/.\n- Exact Quote Accuracy: Fee schedules, operating hours, and contact numbers must be quoted verbatim.\n- Fallback on Uncertainty: If an inquiry cannot be answered by official documentation, route to human support.",
        "Verified in ai/datasets/03_huntiq_platform_knowledge/evaluation/.",
    ),
    (
        "ai/knowledge/marketplace_terminology.md",
        "Authoritative Ethiopian Marketplace Terminology Reference",
        "Comprehensive glossary defining commercial marketplace terminology across Amharic, Oromo, Tigrinya, and English.",
        "- Real Estate: ካርታ (Title deed), ፔንሲዮን (Pension/Guest house), ኮንዶሚኒየም (Condominium), ቪላ (Villa), ሰርቪስ ቤት (Service quarter).\n- Automotive: ቀረጥ የተከፈለ (Duty paid), ኮድ 2 (Private plate), ኮድ 3 (Commercial plate), ቦሌ (Manual), አውቶማቲክ (Automatic).\n- Electronics: ቦክስ ያልተከፈተ (Brand new sealed), ከዱባይ የመጣ (Dubai used import), ኦሪጅናል (Genuine).\n- Agriculture: ጤፍ (Teff), ማኛ (Magna teff), ሰርገኛ (Sergegna teff), ቡና (Coffee), ኩንታል (Quintal).",
        "Maintained collaboratively by HUNTIQ domain specialists.",
    ),
    (
        "ai/knowledge/listing_lifecycle.md",
        "Listing Lifecycle & Moderation States",
        "Complete state machine detailing the lifecycle of marketplace listings on HUNTIQ.",
        "- Draft: Created by user, not yet submitted for review.\n- Pending Approval: Automated and manual screening for prohibited items (15-30 min).\n- Active: Searchable and visible to all users (30-day default duration).\n- Promoted: Highlighted or featured via Telebirr/CBE payment.\n- Expired: Automatically hidden after 30 days; renewable for free.\n- Sold / Inactive: Marked as sold by seller or removed upon request.",
        "Implemented in backend/internal/models/listing.go.",
    ),
    (
        "ai/knowledge/knowledge_versioning.md",
        "Knowledge Document Versioning & Change Management",
        "Governs changes to official fee schedules, terms of service, and platform feature descriptions.",
        "- Semantic Versioning: All documents in ai/knowledge/ maintain an effective date and version header.\n- Change Approval: Modifications require dual authorization from Product Operations and Legal.\n- Immediate Cache Invalidation: Updates trigger automated invalidation of retrieval chunk caches.",
        "Version history tracked in Git log.",
    ),
    (
        "ai/knowledge/marketplace_category_reference.md",
        "Marketplace Category Taxonomy & Hierarchy Reference",
        "Complete taxonomic reference of all 8 primary marketplace categories and 45 subcategories.",
        "- Category 1: Products & Electronics (Phones, Computers, TVs, Audio, Cameras).\n- Category 2: Real Estate & Property (Apartments, Houses, Condos, Commercial, Land).\n- Category 3: Vehicles & Automotive (Cars, Motorcycles, Trucks, Machinery, Parts).\n- Category 4: Jobs & Employment (IT, Accounting, Drivers, Skilled Trades, Domestic).\n- Category 5: Agriculture & Commodities (Cereals, Coffee, Livestock, Farm Tools).\n- Category 6: Professional Services (Plumbers, Electricians, Tutors, Tailors).\n- Category 7: Hotels & Accommodation (Hotels, Resorts, Guest Houses, Homestays).\n- Category 8: Event Tickets & Announcements (Concerts, Expos, Auctions, Notices).",
        "Directly synchronized with PostgreSQL categories table.",
    ),
    (
        "ai/knowledge/membership_reference.md",
        "Membership Tiers, Privileges & Commercial Features Reference",
        "Detailed specification of HUNTIQ user account tiers.",
        "- Free Guest: Unlimited browsing, search, and direct seller chat. Up to 5 free listings.\n- Gold Merchant: Verified Gold Badge, up to 100 active listings, priority search ranking, store profile page (500 ETB/month).\n- Enterprise Seller: Unlimited listings, dedicated account manager, API integration (Custom quote).",
        "Configured in backend/internal/service/membership_service.go.",
    ),
    (
        "ai/knowledge/payment_reference.md",
        "Payment Channels, Fee Schedules & Financial Settlement Reference",
        "Comprehensive reference on platform financial operations in Ethiopia.",
        "- Telebirr: Native integration for mobile wallet payments via USSD/SuperApp.\n- Commercial Bank of Ethiopia (CBE): CBE Birr and direct bank transfer reconciliation.\n- Chapa: Visa, Mastercard, and international diaspora card processing.\n- VAT Compliance: All platform fees include 15% Ethiopian Value Added Tax (VAT).\n- Zero Commission on Sales: Peer-to-peer sales carry 0% transaction commission.",
        "Maintained by Finance and Payments engineering team.",
    ),
    (
        "ai/knowledge/location_reference.md",
        "Ethiopian Geographical Hierarchy & Landmark Reference",
        "Authoritative reference on Ethiopian administrative boundaries and micro-neighborhoods.",
        "- Administrative Tiers: Country (Ethiopia) -> Regions & Chartered Cities -> Subcities (Kifle Ketema) -> Woredas -> Kebeles.\n- Addis Ababa Subcities: Bole, Yeka, Kirkos, Arada, Gullele, Lideta, Kolfe Keranio, Nifas Silk-Lafto, Akaki-Kality, Addis Ketema, Lemi Kura.\n- Major Regional Cities: Adama, Bahir Dar, Hawassa, Mekelle, Dire Dawa, Gondar, Jimma, Dessie, Bishoftu, Harar.",
        "Mirrored from backend/internal/service/geocoding/.",
    ),
    (
        "ai/knowledge/platform_policy_reference.md",
        "Platform Terms of Service & User Safety Policy",
        "Official behavioral rules for buyers, sellers, landlords, and service providers.",
        "- Honest Listings: Sellers must accurately describe condition, defects, and price.\n- Prohibited Contraband: Immediate ban for listing weapons, illegal drugs, or counterfeit goods.\n- Anti-Harassment: Respectful communication required across all in-app chat channels.",
        "Enforced by Go backend moderation service.",
    ),
    (
        "ai/knowledge/support_reference.md",
        "Customer Support Channels, Escalation Paths & Operating Hours",
        "Official support channels and response time commitments for HUNTIQ users.",
        "- Channels: In-app live chat, email (support@huntiq.com), and official phone hotline.\n- Operating Hours: Monday through Saturday, 8:30 AM to 6:00 PM East Africa Time (EAT).\n- First Response SLA: Within 2 hours for standard tickets; within 15 minutes for payment issues.",
        "Monitored by Customer Experience Operations team.",
    ),
]

for p, t, s, r, pr in coord + gov + know:
    expand(p, t, s, r, pr)

print("Expanded all coordination, governance, and knowledge documents successfully.")
