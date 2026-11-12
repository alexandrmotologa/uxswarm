# Persona specifications

uxswarm defines synthetic user personas in declarative YAML files. Each persona models realistic human browsing behavior, cognitive limits, and domain expectations.

## Persona schema

A persona entry contains the following attributes:

- `id`: Unique identifier string (lowercase alphanumeric with hyphens).
- `name`: Human-readable persona label and role.
- `age`: Demographic reference integer.
- `archetype`: Behavioral profile classification:
  - `HIGH_IMPATIENCE_LOW_TECH`: Scans only for primary buttons, skims text, abandons on delays or complex forms.
  - `SCRUTINIZER_HIGH_TECH`: Reviews security disclaimers, password requirements, and settings before committing.
  - `MOBILE_THUMB_DRIVEN`: Navigates via mobile viewport, sensitive to small tap targets and modal obstruction.
  - `ACCESSIBILITY_RELIANT`: Relies on clear label associations, keyboard order, and high-contrast hierarchy.
  - `PRICE_SENSITIVE_SHOPPER`: Scrutinizes pricing terms, billing cycles, and cancellation conditions.
- `patience_threshold_seconds`: Maximum acceptable idle or confusion duration before abandoning the flow.
- `technical_literacy`: Integer score from 1 (minimal digital experience) to 10 (software engineer).
- `attention_span`: Attention level classification (`LOW`, `MEDIUM`, `HIGH`).
- `viewport`: Screen dimensions and mobile flag (`width`, `height`, `is_mobile`).
- `traits`: List of explicit behavioral heuristics evaluated during navigation.

## Default personas

### Sarah Chen (CFO)
- Archetype: `HIGH_IMPATIENCE_LOW_TECH`
- Patience threshold: 4.0 seconds
- Technical literacy: 3 / 10
- Attention span: LOW
- Focus: Fast signup, visible social login, clear pricing, immediate path to value.

### Marcus Vance (DevSecOps Architect)
- Archetype: `SCRUTINIZER_HIGH_TECH`
- Patience threshold: 15.0 seconds
- Technical literacy: 9 / 10
- Attention span: HIGH
- Focus: Two-factor authentication, terms of service, compliance disclaimers, permission clarity.

### Elena Rostova (Field Sales Rep)
- Archetype: `MOBILE_THUMB_DRIVEN`
- Viewport: 390 x 844 (Mobile)
- Patience threshold: 5.0 seconds
- Technical literacy: 6 / 10
- Attention span: MEDIUM
- Focus: Thumb reach zone, tap targets above 44 pixels, form input responsiveness.

### David Kim (Accessibility Advocate)
- Archetype: `ACCESSIBILITY_RELIANT`
- Patience threshold: 10.0 seconds
- Technical literacy: 5 / 10
- Attention span: HIGH
- Focus: Form label associations, error messages adjacent to inputs, readable contrast.

### Chloe Martin (Budget Conscious Freelancer)
- Archetype: `PRICE_SENSITIVE_SHOPPER`
- Patience threshold: 6.0 seconds
- Technical literacy: 6 / 10
- Attention span: MEDIUM
- Focus: Free tier visibility, transparent renewal terms, no surprise checkout fees.
