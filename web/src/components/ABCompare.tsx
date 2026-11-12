import React, { useState, useEffect } from 'react';
import { RunRecord } from '../types';
import { GitCompare, TrendingUp, TrendingDown, CheckCircle2, ArrowRight } from 'lucide-react';

interface ABCompareProps {
  runs: RunRecord[];
}

interface ComparisonData {
  winner: string;
  recommendation: string;
  sus_a: number;
  sus_b: number;
  delta_sus: number;
  grade_a: string;
  grade_b: string;
  friction_a: number;
  friction_b: number;
  delta_friction: number;
  drop_offs_a: number;
  drop_offs_b: number;
  delta_drop_offs: number;
  persona_deltas: {
    persona_name: string;
    archetype: string;
    sus_a: number;
    sus_b: number;
    delta_sus: number;
    friction_a: number;
    friction_b: number;
    outcome_a: string;
    outcome_b: string;
  }[];
  resolved_issues: string[];
}

export const ABCompare: React.FC<ABCompareProps> = ({ runs }) => {
  const [runAId, setRunAId] = useState<string>(runs[1]?.id || runs[0]?.id || '');
  const [runBId, setRunBId] = useState<string>(runs[0]?.id || '');
  const [cmpData, setCmpData] = useState<ComparisonData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (!runAId || !runBId) return;
    setLoading(true);

    fetch(`/api/compare?run_a=${runAId}&run_b=${runBId}`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch comparison');
        return res.json();
      })
      .then((data: ComparisonData) => {
        setCmpData(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setLoading(false);
      });
  }, [runAId, runBId]);

  if (runs.length < 2) {
    return (
      <div style={{ maxWidth: '900px', margin: '3rem auto', textAlign: 'center', color: '#94a3b8' }}>
        <GitCompare size={48} style={{ margin: '0 auto 1rem', opacity: 0.5 }} />
        <h3 style={{ fontSize: '1.25rem', color: '#f8fafc' }}>Need at least 2 usability studies</h3>
        <p style={{ marginTop: '0.5rem', fontSize: '0.9rem' }}>
          Execute another study or test a redesign variant to unlock A/B usability differential analysis.
        </p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      {/* Selector Controls */}
      <div style={{
        background: '#141e33',
        borderRadius: '14px',
        border: '1px solid #1e293b',
        padding: '1.5rem',
        marginBottom: '2rem',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        flexWrap: 'wrap',
        gap: '1rem',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <GitCompare size={24} color="#3b82f6" />
          <div>
            <h2 style={{ fontSize: '1.3rem', fontWeight: 800 }}>A/B Usability Differential Suite</h2>
            <div style={{ fontSize: '0.8rem', color: '#94a3b8' }}>Compare baseline and redesign variants</div>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <div>
            <div style={{ fontSize: '0.75rem', color: '#94a3b8', marginBottom: '0.25rem' }}>Baseline (Variant A)</div>
            <select
              value={runAId}
              onChange={(e) => setRunAId(e.target.value)}
              style={{
                background: '#090d16',
                border: '1px solid #1e293b',
                color: 'white',
                padding: '6px 12px',
                borderRadius: '8px',
                fontSize: '0.85rem',
                outline: 'none',
              }}
            >
              {runs.map((r) => (
                <option key={r.id} value={r.id}>{r.name} (SUS: {r.average_sus.toFixed(1)})</option>
              ))}
            </select>
          </div>

          <ArrowRight size={18} color="#64748b" style={{ marginTop: '1.25rem' }} />

          <div>
            <div style={{ fontSize: '0.75rem', color: '#94a3b8', marginBottom: '0.25rem' }}>Candidate (Variant B)</div>
            <select
              value={runBId}
              onChange={(e) => setRunBId(e.target.value)}
              style={{
                background: '#090d16',
                border: '1px solid #1e293b',
                color: 'white',
                padding: '6px 12px',
                borderRadius: '8px',
                fontSize: '0.85rem',
                outline: 'none',
              }}
            >
              {runs.map((r) => (
                <option key={r.id} value={r.id}>{r.name} (SUS: {r.average_sus.toFixed(1)})</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {loading ? (
        <div style={{ padding: '4rem', textAlign: 'center', color: '#94a3b8' }}>Computing differential metrics...</div>
      ) : cmpData ? (
        <div>
          {/* Winner Banner */}
          <div style={{
            background: 'linear-gradient(135deg, rgba(37,99,235,0.15), rgba(139,92,246,0.1))',
            border: '1px solid #3b82f6',
            borderRadius: '14px',
            padding: '1.5rem',
            marginBottom: '2rem',
          }}>
            <div style={{ fontSize: '0.75rem', color: '#93c5fd', fontWeight: 800, textTransform: 'uppercase' }}>
              Evaluation Recommendation
            </div>
            <div style={{ fontSize: '1.5rem', fontWeight: 800, color: '#f8fafc', margin: '0.35rem 0' }}>
              {cmpData.winner}
            </div>
            <div style={{ fontSize: '0.9rem', color: '#cbd5e1' }}>
              {cmpData.recommendation}
            </div>
          </div>

          {/* Delta KPI Grid */}
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: '1.25rem', marginBottom: '2.5rem' }}>
            <div style={{ background: '#141e33', border: '1px solid #1e293b', borderRadius: '12px', padding: '1.5rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: cmpData.delta_sus >= 0 ? '#10b981' : '#ef4444' }}>
                {cmpData.delta_sus >= 0 ? <TrendingUp size={24} /> : <TrendingDown size={24} />}
                <span style={{ fontSize: '2rem', fontWeight: 800 }}>
                  {cmpData.delta_sus >= 0 ? `+${cmpData.delta_sus.toFixed(1)}` : cmpData.delta_sus.toFixed(1)} pts
                </span>
              </div>
              <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', marginTop: '0.35rem' }}>
                SUS Benchmark Delta (A: {cmpData.sus_a.toFixed(1)} &rarr; B: {cmpData.sus_b.toFixed(1)})
              </div>
            </div>

            <div style={{ background: '#141e33', border: '1px solid #1e293b', borderRadius: '12px', padding: '1.5rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: cmpData.delta_friction <= 0 ? '#10b981' : '#ef4444' }}>
                <span style={{ fontSize: '2rem', fontWeight: 800 }}>
                  {cmpData.delta_friction.toFixed(2)}
                </span>
              </div>
              <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', marginTop: '0.35rem' }}>
                Friction Index Delta (A: {cmpData.friction_a.toFixed(2)} &rarr; B: {cmpData.friction_b.toFixed(2)})
              </div>
            </div>

            <div style={{ background: '#141e33', border: '1px solid #1e293b', borderRadius: '12px', padding: '1.5rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: cmpData.delta_drop_offs <= 0 ? '#10b981' : '#ef4444' }}>
                <span style={{ fontSize: '2rem', fontWeight: 800 }}>
                  {cmpData.delta_drop_offs > 0 ? `+${cmpData.delta_drop_offs}` : cmpData.delta_drop_offs}
                </span>
              </div>
              <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', marginTop: '0.35rem' }}>
                Session Abandonment Delta (A: {cmpData.drop_offs_a} &rarr; B: {cmpData.drop_offs_b})
              </div>
            </div>
          </div>

          {/* Persona Breakdown Table */}
          <div style={{ background: '#141e33', border: '1px solid #1e293b', borderRadius: '14px', overflow: 'hidden', marginBottom: '2rem' }}>
            <div style={{ padding: '1.25rem', borderBottom: '1px solid #1e293b', fontWeight: 700 }}>
              Persona Level Differential
            </div>
            <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '0.9rem' }}>
              <thead>
                <tr style={{ background: '#090d16', color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>
                  <th style={{ padding: '0.85rem 1.25rem' }}>Persona</th>
                  <th style={{ padding: '0.85rem 1rem' }}>Variant A (SUS)</th>
                  <th style={{ padding: '0.85rem 1rem' }}>Variant B (SUS)</th>
                  <th style={{ padding: '0.85rem 1rem' }}>Delta SUS</th>
                  <th style={{ padding: '0.85rem 1rem' }}>Outcome</th>
                </tr>
              </thead>
              <tbody>
                {cmpData.persona_deltas?.map((p, idx) => (
                  <tr key={idx} style={{ borderBottom: '1px solid #1e293b' }}>
                    <td style={{ padding: '1rem 1.25rem' }}>
                      <div style={{ fontWeight: 700, color: '#f8fafc' }}>{p.persona_name}</div>
                      <div style={{ fontSize: '0.75rem', color: '#94a3b8' }}>{p.archetype}</div>
                    </td>
                    <td style={{ padding: '1rem' }}>{p.sus_a.toFixed(1)}</td>
                    <td style={{ padding: '1rem', fontWeight: 700 }}>{p.sus_b.toFixed(1)}</td>
                    <td style={{ padding: '1rem', fontWeight: 700, color: p.delta_sus >= 0 ? '#10b981' : '#ef4444' }}>
                      {p.delta_sus >= 0 ? `+${p.delta_sus.toFixed(1)}` : p.delta_sus.toFixed(1)}
                    </td>
                    <td style={{ padding: '1rem' }}>
                      <span style={{ fontSize: '0.8rem', color: '#cbd5e1' }}>{p.outcome_a} &rarr; <strong>{p.outcome_b}</strong></span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Resolved Friction Points */}
          {cmpData.resolved_issues && cmpData.resolved_issues.length > 0 && (
            <div style={{ background: 'rgba(16,185,129,0.08)', border: '1px solid #10b981', borderRadius: '12px', padding: '1.25rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: '#34d399', fontWeight: 700, marginBottom: '0.5rem' }}>
                <CheckCircle2 size={18} /> Usability Barriers Resolved in Candidate Variant
              </div>
              <ul style={{ paddingLeft: '1.25rem', fontSize: '0.85rem', color: '#cbd5e1' }}>
                {cmpData.resolved_issues.map((res, idx) => (
                  <li key={idx} style={{ marginBottom: '0.25rem' }}>{res}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      ) : null}
    </div>
  );
};
