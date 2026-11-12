package persona

import (
	"fmt"
	"math"

	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// FrustrationTracker tracks accumulated cognitive frustration for a persona session.
type FrustrationTracker struct {
	Persona     config.Persona
	Frustration float64 // 0.0 (content) to 1.0+ (abandonment)
	DropOff     bool
	Reason      string
}

// NewFrustrationTracker constructs a tracker initialized to zero frustration.
func NewFrustrationTracker(persona config.Persona) *FrustrationTracker {
	return &FrustrationTracker{
		Persona: persona,
	}
}

// Update evaluates an action step and increments frustration based on friction indicators.
func (ft *FrustrationTracker) Update(hesitationSeconds float64, isDeadClick bool, isBacktrack bool) (float64, bool) {
	threshold := ft.Persona.PatienceThresholdSeconds
	if threshold <= 0 {
		threshold = 5.0
	}

	// 1. Hesitation impact: if hesitation approaches or exceeds threshold
	if hesitationSeconds > 1.5 {
		ratio := hesitationSeconds / threshold
		delta := ratio * 0.35
		// Archetypes with low patience get frustrated faster
		if ft.Persona.Archetype == config.ArchetypeHighImpatienceLowTech {
			delta *= 1.4
		}
		ft.Frustration += delta
	}

	// 2. Dead click impact
	if isDeadClick {
		deadPenalty := 0.35
		if ft.Persona.Archetype == config.ArchetypeMobileThumbDriven {
			deadPenalty = 0.40 // Missed mobile tap targets are especially annoying
		}
		ft.Frustration += deadPenalty
	}

	// 3. Backtrack loop impact
	if isBacktrack {
		ft.Frustration += 0.25
	}

	// Natural slight recovery if interaction was smooth (0 hesitation and valid action)
	if hesitationSeconds < 1.0 && !isDeadClick && !isBacktrack && ft.Frustration > 0.1 {
		ft.Frustration -= 0.05
	}

	ft.Frustration = math.Round(ft.Frustration*100) / 100

	if ft.Frustration >= 1.0 {
		ft.DropOff = true
		if isDeadClick {
			ft.Reason = fmt.Sprintf("Repeated unresponsive clicks exceeded %s's patience limit (%.1fs threshold).", ft.Persona.Name, threshold)
		} else if isBacktrack {
			ft.Reason = fmt.Sprintf("Circular navigation loop caused %s to abandon the session.", ft.Persona.Name)
		} else {
			ft.Reason = fmt.Sprintf("Excessive cognitive friction and hesitation (%.1fs) exceeded %s's patience threshold.", hesitationSeconds, ft.Persona.Name)
		}
	}

	return ft.Frustration, ft.DropOff
}
