# Nielsen Norman usability heuristics mapping

uxswarm correlates detected user friction events (hesitations, dead clicks, route loops, and abandonment) with the 10 Jakob Nielsen usability heuristics.

## Heuristics reference

1. Visibility of system status
   - Symptoms: Lack of loading spinners, silent form submission delays, missing progress steps.
   - Persona reaction: Repeated clicks on submit buttons, wondering if the application froze.
   - Detection: Hesitation on submission or multiple rapid clicks on the same button.

2. Match between system and the real world
   - Symptoms: Jargon, internal database IDs, technical error codes, unfamiliar naming.
   - Persona reaction: Confusion over button labels or form field purposes.
   - Detection: Dead clicks on explanatory text and prolonged pauses before dropdown selections.

3. User control and freedom
   - Symptoms: Inability to edit prior steps, missing back buttons, trap modals.
   - Persona reaction: Rapid navigation attempts to previous pages or browser back button usage.
   - Detection: Backtrack loops (A -> B -> A) and cancelled modal sequences.

4. Consistency and standards
   - Symptoms: Non-standard button colors, primary actions placed on the left instead of right, unconventional icons.
   - Persona reaction: Misclicks on decorative elements mistaken for buttons.
   - Detection: Clicking on non-clickable headings or badges.

5. Error prevention
   - Symptoms: Unchecked inputs that fail only on final submit, missing confirmation on destructive steps.
   - Persona reaction: Frustration upon seeing a full-page reload with erased form data.
   - Detection: Form rejection events requiring re-entry of credentials or fields.

6. Recognition rather than recall
   - Symptoms: Hidden navigation menus, options buried behind obscure hover states.
   - Persona reaction: Hovering over multiple unrelated elements to discover available actions.
   - Detection: Cursor wander latency without element engagement.

7. Flexibility and efficiency of use
   - Symptoms: Forcing experienced users through lengthy step-by-step wizards without bulk or keyboard shortcuts.
   - Persona reaction: Low satisfaction ratings from technical personas (e.g. DevSecOps Architect).
   - Detection: Session duration far exceeding target completion time.

8. Aesthetic and minimalist design
   - Symptoms: Visual clutter, competing calls to action, dense paragraphs on conversion landing pages.
   - Persona reaction: Attention drift, misclicking secondary links instead of primary flow CTA.
   - Detection: Drop in reading velocity and high ratio of secondary to primary element focus.

9. Help users recognize, diagnose, and recover from errors
   - Symptoms: Generic error notices ("An error occurred") without specific field guidance.
   - Persona reaction: Hesitation around error banner, inability to identify which field failed.
   - Detection: Long pause (> 5s) following an error alert with no subsequent input change.

10. Help and documentation
    - Symptoms: Unclear tooltips, missing inline examples for complex inputs like phone numbers or tax IDs.
    - Persona reaction: Hesitation and abandonment during configuration steps.
    - Detection: Stalling on optional configuration fields with no help hints.
