package patch

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
)

// CodePatch represents an actionable software engineering patch to fix a usability issue.
type CodePatch struct {
	HeuristicID    int    `json:"heuristic_id"`
	Title          string `json:"title"`
	TargetSelector string `json:"target_selector,omitempty"`
	Framework      string `json:"framework"` // CSS, Tailwind, HTML, React
	BeforeCode     string `json:"before_code"`
	AfterCode      string `json:"after_code"`
	Explanation    string `json:"explanation"`
}

// GeneratePatches produces concrete code modifications for all detected usability findings.
func GeneratePatches(findings []metrics.HeuristicFinding, steps []metrics.StepEvent) []CodePatch {
	patches := make([]CodePatch, 0)

	// Collect any dead click selectors from steps
	deadClickSelector := ".promo-banner, [data-interactive='false']"
	for _, st := range steps {
		if st.IsDeadClick && st.TargetSelector != "" {
			deadClickSelector = st.TargetSelector
			break
		}
	}

	for _, f := range findings {
		switch f.HeuristicID {
		case 1: // Visibility of system status (Hesitation after submit)
			patches = append(patches, CodePatch{
				HeuristicID:    1,
				Title:          "Add Loading State & Button Spinner on Form Submission",
				TargetSelector: "button[type='submit'], .btn-submit",
				Framework:      "CSS & React",
				BeforeCode: `<button type="submit" className="btn-submit">
  Continue to Workspace
</button>`,
				AfterCode: `<button 
  type="submit" 
  disabled={isSubmitting} 
  className="btn-submit disabled:opacity-60 disabled:cursor-wait flex items-center justify-center gap-2"
>
  {isSubmitting ? (
    <>
      <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
      <span>Processing...</span>
    </>
  ) : (
    <span>Continue to Workspace</span>
  )}
</button>`,
				Explanation: "Eliminates hesitation by providing instantaneous visual feedback when users submit forms, preventing duplicate clicks.",
			})

		case 4: // Consistency and standards (Dead Clicks)
			patches = append(patches, CodePatch{
				HeuristicID:    4,
				Title:          "Remove Deceptive Button Styling from Non-Interactive Elements",
				TargetSelector: deadClickSelector,
				Framework:      "CSS & Tailwind",
				BeforeCode: fmt.Sprintf(`/* Deceptive element styling causing dead clicks */
%s {
  cursor: pointer;
  background: #f1f5f9;
  border-radius: 8px;
}`, deadClickSelector),
				AfterCode: fmt.Sprintf(`/* Option A: Make element explicitly non-interactive */
%s {
  cursor: default;
  pointer-events: none;
  user-select: none;
  background: transparent;
  border: 1px solid #e2e8f0;
}

/* Option B: If meant to be clicked, convert to semantic anchor or button */
<a href="/pricing" className="inline-flex items-center text-blue-600 hover:underline">
  Free 14-Day Enterprise Trial &rarr;
</a>`, deadClickSelector),
				Explanation: "Resolves dead clicks by ensuring clickable visual cues (hover effects, cursor pointer, box shadows) exist only on interactive controls.",
			})

		case 7: // Flexibility and efficiency of use (Mobile tap targets < 44px)
			patches = append(patches, CodePatch{
				HeuristicID:    7,
				Title:          "Enforce Apple HIG Minimum 44px Touch Targets on Mobile",
				TargetSelector: "button, input, select, .btn-submit",
				Framework:      "Tailwind & CSS",
				BeforeCode: `/* Sub-optimal touch target height (38px) */
.btn-submit {
  height: 38px;
  padding: 8px 16px;
}`,
				AfterCode: `/* Compliant touch target adhering to WCAG and Apple HIG */
@media (max-width: 640px) {
  .btn-submit,
  button,
  input[type="submit"] {
    min-height: 44px;
    min-width: 44px;
    padding: 12px 20px;
    touch-action: manipulation;
  }
}

/* Tailwind utility equivalent */
<button className="min-h-[44px] min-w-[44px] py-3 px-5 touch-manipulation w-full ...">`,
				Explanation: "Prevents missed taps and reach frustration on mobile viewports by sizing controls to human thumb dimensions.",
			})

		case 8: // Aesthetic and minimalist design (Drop-off / form abandonment)
			patches = append(patches, CodePatch{
				HeuristicID:    8,
				Title:          "Progressive Disclosure: Defer Non-Essential Onboarding Fields",
				TargetSelector: "form.onboarding-form",
				Framework:      "HTML & React",
				BeforeCode: `<form>
  <input name="name" required />
  <input name="email" required />
  <input name="password" required />
  <input name="workspace_name" required />
  <select name="company_size" required />
  <input name="invite_colleague" />
  <button type="submit">Sign Up</button>
</form>`,
				AfterCode: `<!-- Step 1: Micro-conversion (Zero Friction) -->
<form onSubmit={handleInitialSignup}>
  <input name="email" type="email" placeholder="Work email" required />
  <input name="password" type="password" required />
  <button type="submit">Get Started</button>
</form>

<!-- Step 2: Post-signup onboarding wizard inside authenticated dashboard -->
<!-- Workspace name and teammate invitations deferred until user experiences first value -->`,
				Explanation: "Cuts drop-offs by reducing upfront cognitive load. Only request critical authentication credentials during the initial signup conversion.",
			})

		case 3: // User control and freedom (Backtrack loops)
			patches = append(patches, CodePatch{
				HeuristicID:    3,
				Title:          "Add Explicit Emergency Navigation Exit & Breadcrumbs",
				TargetSelector: "nav.step-navigation",
				Framework:      "HTML & React",
				BeforeCode: `<div className="wizard-step">
  <button type="submit">Next</button>
  <!-- No back or cancel control provided -->
</div>`,
				AfterCode: `<div className="flex items-center justify-between mt-4">
  <button 
    type="button" 
    onClick={goToPreviousStep}
    className="text-slate-400 hover:text-slate-200 text-sm font-medium"
  >
    &larr; Back to previous step
  </button>
  <button type="submit" className="btn-primary">
    Continue
  </button>
</div>`,
				Explanation: "Eliminates circular browser back-button loops by providing visible, state-preserving backward navigation.",
			})
		}
	}

	return patches
}

// FormatPatchText prints human-readable patch code to console.
func FormatPatchText(p CodePatch) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("=== Patch for Heuristic %d: %s [%s] ===\n", p.HeuristicID, p.Title, p.Framework))
	if p.TargetSelector != "" {
		b.WriteString(fmt.Sprintf("Selector: %s\n\n", p.TargetSelector))
	}
	b.WriteString("--- BEFORE ---\n")
	b.WriteString(p.BeforeCode + "\n\n")
	b.WriteString("+++ AFTER (FIX) +++\n")
	b.WriteString(p.AfterCode + "\n\n")
	b.WriteString("Why: " + p.Explanation + "\n")
	return b.String()
}
