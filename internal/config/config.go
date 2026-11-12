package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Archetype constants representing synthetic user behavioral classifications.
const (
	ArchetypeHighImpatienceLowTech = "HIGH_IMPATIENCE_LOW_TECH"
	ArchetypeScrutinizerHighTech   = "SCRUTINIZER_HIGH_TECH"
	ArchetypeMobileThumbDriven     = "MOBILE_THUMB_DRIVEN"
	ArchetypeAccessibilityReliant  = "ACCESSIBILITY_RELIANT"
	ArchetypePriceSensitive        = "PRICE_SENSITIVE_SHOPPER"
)

// Viewport represents screen dimensions and device profile.
type Viewport struct {
	Width    int  `yaml:"width" json:"width"`
	Height   int  `yaml:"height" json:"height"`
	IsMobile bool `yaml:"is_mobile" json:"is_mobile"`
}

// Persona represents an individual synthetic user with cognitive and psychological parameters.
type Persona struct {
	ID                       string   `yaml:"id" json:"id"`
	Name                     string   `yaml:"name" json:"name"`
	Age                      int      `yaml:"age" json:"age"`
	Archetype                string   `yaml:"archetype" json:"archetype"`
	PatienceThresholdSeconds float64  `yaml:"patience_threshold_seconds" json:"patience_threshold_seconds"`
	TechnicalLiteracy        int      `yaml:"technical_literacy" json:"technical_literacy"` // 1 (novice) to 10 (expert)
	AttentionSpan            string   `yaml:"attention_span" json:"attention_span"`         // LOW, MEDIUM, HIGH
	Viewport                 Viewport `yaml:"viewport" json:"viewport"`
	Traits                   []string `yaml:"traits" json:"traits"`
}

// SuiteMeta holds metadata and parameters for an automated usability study.
type SuiteMeta struct {
	Name               string `yaml:"name" json:"name"`
	TargetURL          string `yaml:"target_url" json:"target_url"`
	Goal               string `yaml:"goal" json:"goal"`
	MaxSessionDuration string `yaml:"max_session_duration" json:"max_session_duration"`
}

// Suite represents the complete declarative usability test suite.
type Suite struct {
	Version  string    `yaml:"version" json:"version"`
	Suite    SuiteMeta `yaml:"suite" json:"suite"`
	Personas []Persona `yaml:"personas" json:"personas"`
}

// ParsedMaxDuration returns parsed time.Duration for MaxSessionDuration or a 3-minute default.
func (s *Suite) ParsedMaxDuration() time.Duration {
	if s.Suite.MaxSessionDuration == "" {
		return 3 * time.Minute
	}
	d, err := time.ParseDuration(s.Suite.MaxSessionDuration)
	if err != nil {
		return 3 * time.Minute
	}
	return d
}

// LoadSuite loads and validates a YAML persona suite file.
func LoadSuite(path string) (*Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read suite file %s: %w", path, err)
	}

	var s Suite
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse suite yaml %s: %w", path, err)
	}

	if len(s.Personas) == 0 {
		return nil, fmt.Errorf("suite %s contains no personas", path)
	}

	for i := range s.Personas {
		p := &s.Personas[i]
		if p.ID == "" {
			p.ID = fmt.Sprintf("persona-%d", i+1)
		}
		if p.PatienceThresholdSeconds <= 0 {
			p.PatienceThresholdSeconds = 6.0
		}
		if p.TechnicalLiteracy <= 0 {
			p.TechnicalLiteracy = 5
		}
		if p.TechnicalLiteracy > 10 {
			p.TechnicalLiteracy = 10
		}
		if p.Viewport.Width == 0 || p.Viewport.Height == 0 {
			if p.Archetype == ArchetypeMobileThumbDriven || p.Viewport.IsMobile {
				p.Viewport.Width = 390
				p.Viewport.Height = 844
				p.Viewport.IsMobile = true
			} else {
				p.Viewport.Width = 1440
				p.Viewport.Height = 900
				p.Viewport.IsMobile = false
			}
		}
	}

	return &s, nil
}

// GetDefaultPersonas provides a balanced five-persona swarm covering major archetypes.
func GetDefaultPersonas() []Persona {
	return []Persona{
		{
			ID:                       "impatient-executive",
			Name:                     "Sarah Chen (CFO)",
			Age:                      48,
			Archetype:                ArchetypeHighImpatienceLowTech,
			PatienceThresholdSeconds: 4.0,
			TechnicalLiteracy:        3,
			AttentionSpan:            "LOW",
			Viewport:                 Viewport{Width: 1440, Height: 900, IsMobile: false},
			Traits: []string{
				"Skims text rapidly without reading instructions",
				"Expects high-contrast primary buttons",
				"Abandons on unexpected payment or billing requirements",
			},
		},
		{
			ID:                       "skeptical-security-lead",
			Name:                     "Marcus Vance (DevSecOps Lead)",
			Age:                      36,
			Archetype:                ArchetypeScrutinizerHighTech,
			PatienceThresholdSeconds: 15.0,
			TechnicalLiteracy:        9,
			AttentionSpan:            "HIGH",
			Viewport:                 Viewport{Width: 1440, Height: 900, IsMobile: false},
			Traits: []string{
				"Scrutinizes password requirements and security disclaimers",
				"Checks privacy policies and consent checkboxes",
				"Expects visible SSO or 2FA options",
			},
		},
		{
			ID:                       "distracted-mobile-user",
			Name:                     "Elena Rostova (Field Sales Rep)",
			Age:                      29,
			Archetype:                ArchetypeMobileThumbDriven,
			PatienceThresholdSeconds: 5.0,
			TechnicalLiteracy:        6,
			AttentionSpan:            "MEDIUM",
			Viewport:                 Viewport{Width: 390, Height: 844, IsMobile: true},
			Traits: []string{
				"Navigates single-handed within thumb reach zone",
				"Struggles with tap targets below 44 pixels",
				"Frustrated by modal dialogs obscuring form fields",
			},
		},
		{
			ID:                       "accessibility-advocate",
			Name:                     "David Kim (Accessibility Advocate)",
			Age:                      42,
			Archetype:                ArchetypeAccessibilityReliant,
			PatienceThresholdSeconds: 10.0,
			TechnicalLiteracy:        5,
			AttentionSpan:            "HIGH",
			Viewport:                 Viewport{Width: 1440, Height: 900, IsMobile: false},
			Traits: []string{
				"Relies on clear form label associations and high contrast",
				"Confused by unlabeled icon-only buttons",
				"Expects explicit error messages adjacent to inputs",
			},
		},
		{
			ID:                       "budget-freelancer",
			Name:                     "Chloe Martin (Freelance Consultant)",
			Age:                      31,
			Archetype:                ArchetypePriceSensitive,
			PatienceThresholdSeconds: 6.0,
			TechnicalLiteracy:        6,
			AttentionSpan:            "MEDIUM",
			Viewport:                 Viewport{Width: 1440, Height: 900, IsMobile: false},
			Traits: []string{
				"Searches for transparent pricing and free tier terms",
				"Abandons if billing details are required for trials",
				"Avoids ambiguous recurring subscription checkboxes",
			},
		},
	}
}
