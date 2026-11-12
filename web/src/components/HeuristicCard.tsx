import React, { useState } from 'react';
import { HeuristicFinding } from '../types';
import { AlertOctagon, AlertTriangle, Info, Check, Sparkles, Code2, Copy, CheckCheck } from 'lucide-react';

interface HeuristicCardProps {
  findings: HeuristicFinding[];
}

export const HeuristicCard: React.FC<HeuristicCardProps> = ({ findings }) => {
  const [selectedSeverity, setSelectedSeverity] = useState<'ALL' | 'CRITICAL' | 'MAJOR' | 'MINOR'>('ALL');
  const [expandedPatch, setExpandedPatch] = useState<number | null>(null);
  const [copiedId, setCopiedId] = useState<number | null>(null);

  const filtered = findings.filter((f) => {
    if (selectedSeverity === 'ALL') return true;
    return f.severity === selectedSeverity;
  });

  const getSeverityBadge = (sev: string) => {
    switch (sev) {
      case 'CRITICAL':
        return {
          icon: <AlertOctagon size={14} />,
          bg: 'rgba(239,68,68,0.15)',
          border: '#ef4444',
          color: '#f87171',
        };
      case 'MAJOR':
        return {
          icon: <AlertTriangle size={14} />,
          bg: 'rgba(245,158,11,0.15)',
          border: '#f59e0b',
          color: '#fbbf24',
        };
      default:
        return {
          icon: <Info size={14} />,
          bg: 'rgba(59,130,246,0.15)',
          border: '#3b82f6',
          color: '#60a5fa',
        };
    }
  };

  const getHeuristicPatch = (id: number) => {
    switch (id) {
      case 1:
        return {
          framework: 'CSS & React',
          code: `<button 
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
        };
      case 4:
        return {
          framework: 'CSS & Tailwind',
          code: `/* Fix dead clicks by disabling cursor pointer on non-interactive badge */
.promo-banner {
  cursor: default !important;
  pointer-events: none;
  user-select: none;
}

/* Or convert into an active semantic link */
<a href="/pricing" className="inline-flex items-center text-blue-600 hover:underline">
  Free 14-Day Enterprise Trial &rarr;
</a>`,
        };
      case 7:
        return {
          framework: 'Tailwind CSS',
          code: `/* Enforce Apple HIG minimum 44px tap targets for mobile */
@media (max-width: 640px) {
  button, input[type="submit"], .btn-submit {
    min-height: 44px;
    min-width: 44px;
    padding: 12px 20px;
    touch-action: manipulation;
  }
}

/* Tailwind utility */
<button className="min-h-[44px] min-w-[44px] py-3 px-5 touch-manipulation w-full ...">`,
        };
      case 8:
        return {
          framework: 'HTML & React',
          code: `<!-- Micro-conversion: only request email & password upfront -->
<form onSubmit={handleInitialSignup}>
  <input name="email" type="email" placeholder="Work email" required />
  <input name="password" type="password" required />
  <button type="submit">Get Started</button>
</form>
<!-- Defer workspace configuration and teammate invites until post-signup -->`,
        };
      default:
        return null;
    }
  };

  const copyToClipboard = (text: string, id: number) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  return (
    <div style={{ maxWidth: '1100px', margin: '0 auto', padding: '2rem' }}>
      {/* Header and Filter */}
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '2rem',
        flexWrap: 'wrap',
        gap: '1rem',
      }}>
        <div>
          <h2 style={{ fontSize: '1.4rem', fontWeight: 800 }}>Nielsen Norman Heuristic Audit</h2>
          <p style={{ fontSize: '0.85rem', color: '#94a3b8', marginTop: '0.25rem' }}>
            Usability violations synthesized from synthetic persona hesitation, misclicks, and drop-offs.
          </p>
        </div>

        {/* Severity filter buttons */}
        <div style={{ display: 'flex', gap: '0.5rem', background: '#141e33', padding: '4px', borderRadius: '10px', border: '1px solid #1e293b' }}>
          {(['ALL', 'CRITICAL', 'MAJOR', 'MINOR'] as const).map((sev) => (
            <button
              key={sev}
              onClick={() => setSelectedSeverity(sev)}
              style={{
                padding: '6px 12px',
                borderRadius: '6px',
                fontSize: '0.75rem',
                fontWeight: 700,
                border: 'none',
                cursor: 'pointer',
                background: selectedSeverity === sev ? '#2563eb' : 'transparent',
                color: selectedSeverity === sev ? 'white' : '#94a3b8',
              }}
            >
              {sev}
            </button>
          ))}
        </div>
      </div>

      {filtered.length === 0 ? (
        <div style={{
          background: '#141e33',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '3rem',
          textAlign: 'center',
          color: '#94a3b8',
        }}>
          <Check size={32} color="#10b981" style={{ margin: '0 auto 0.5rem' }} />
          <div style={{ fontWeight: 700, fontSize: '1.1rem', color: '#f8fafc' }}>
            No matching heuristic findings
          </div>
          <p style={{ fontSize: '0.85rem', marginTop: '0.35rem' }}>
            No usability violations recorded for this filter criteria.
          </p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          {filtered.map((h, i) => {
            const badge = getSeverityBadge(h.severity);
            const patchInfo = getHeuristicPatch(h.heuristic_id);
            const isPatchOpen = expandedPatch === h.heuristic_id;

            return (
              <div
                key={i}
                style={{
                  background: '#141e33',
                  border: '1px solid #1e293b',
                  borderRadius: '14px',
                  padding: '1.5rem',
                  boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1)',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '0.75rem' }}>
                  <div>
                    <span style={{ fontSize: '0.75rem', color: '#94a3b8', textTransform: 'uppercase', fontWeight: 700 }}>
                      Heuristic {h.heuristic_id}: {h.heuristic_name}
                    </span>
                    <h3 style={{ fontSize: '1.15rem', fontWeight: 700, color: '#f8fafc', marginTop: '0.2rem' }}>
                      {h.title}
                    </h3>
                  </div>

                  <span style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: '0.35rem',
                    fontSize: '0.75rem',
                    fontWeight: 700,
                    padding: '3px 10px',
                    borderRadius: '999px',
                    background: badge.bg,
                    border: `1px solid ${badge.border}`,
                    color: badge.color,
                  }}>
                    {badge.icon} {h.severity}
                  </span>
                </div>

                <div style={{ fontSize: '0.9rem', color: '#cbd5e1', lineHeight: 1.5, marginBottom: '0.75rem' }}>
                  {h.evidence}
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1rem' }}>
                  <span style={{ fontSize: '0.75rem', color: '#94a3b8' }}>Affected Personas:</span>
                  <div style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap' }}>
                    {h.affected_personas.map((name, pIdx) => (
                      <span
                        key={pIdx}
                        style={{
                          fontSize: '0.75rem',
                          background: '#090d16',
                          color: '#cbd5e1',
                          padding: '2px 8px',
                          borderRadius: '6px',
                          border: '1px solid #1e293b',
                        }}
                      >
                        {name}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Concrete Redesign Specification Box */}
                <div style={{
                  background: 'rgba(37,99,235,0.08)',
                  borderLeft: '4px solid #2563eb',
                  borderRadius: '0 8px 8px 0',
                  padding: '0.85rem 1.25rem',
                  marginBottom: '1rem',
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#93c5fd', fontSize: '0.8rem', fontWeight: 700, marginBottom: '0.25rem' }}>
                    <Sparkles size={14} /> Redesign Specification
                  </div>
                  <div style={{ fontSize: '0.85rem', color: '#f1f5f9', lineHeight: 1.4 }}>
                    {h.recommendation}
                  </div>
                </div>

                {/* Engineering Code Patch Accordion */}
                {patchInfo && (
                  <div style={{ borderTop: '1px solid #1e293b', paddingTop: '0.75rem' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <button
                        onClick={() => setExpandedPatch(isPatchOpen ? null : h.heuristic_id)}
                        style={{
                          background: 'transparent',
                          border: 'none',
                          color: '#38bdf8',
                          fontSize: '0.8rem',
                          fontWeight: 600,
                          cursor: 'pointer',
                          display: 'flex',
                          alignItems: 'center',
                          gap: '0.4rem',
                          padding: 0,
                        }}
                      >
                        <Code2 size={15} /> {isPatchOpen ? 'Hide Engineering Code Patch' : 'View Code Patch (' + patchInfo.framework + ')'}
                      </button>

                      {isPatchOpen && (
                        <button
                          onClick={() => copyToClipboard(patchInfo.code, h.heuristic_id)}
                          style={{
                            background: '#090d16',
                            border: '1px solid #1e293b',
                            borderRadius: '6px',
                            color: copiedId === h.heuristic_id ? '#34d399' : '#94a3b8',
                            fontSize: '0.75rem',
                            padding: '3px 8px',
                            cursor: 'pointer',
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.3rem',
                          }}
                        >
                          {copiedId === h.heuristic_id ? <CheckCheck size={13} /> : <Copy size={13} />}
                          {copiedId === h.heuristic_id ? 'Copied' : 'Copy Patch'}
                        </button>
                      )}
                    </div>

                    {isPatchOpen && (
                      <div style={{ marginTop: '0.75rem', background: '#090d16', borderRadius: '8px', border: '1px solid #1e293b', padding: '1rem', overflowX: 'auto' }}>
                        <pre style={{ fontSize: '0.8rem', color: '#cbd5e1', fontFamily: 'JetBrains Mono, monospace', margin: 0, lineHeight: 1.4 }}>
                          {patchInfo.code}
                        </pre>
                      </div>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};
