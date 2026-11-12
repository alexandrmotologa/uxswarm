import React from 'react';
import { RunRecord, SessionRecord } from '../types';
import { CheckCircle2, XCircle, Smartphone, Monitor, AlertTriangle, Play, Flame } from 'lucide-react';

interface RunOverviewProps {
  run: RunRecord | null;
  sessions: SessionRecord[];
  onSelectSession: (sessionId: string) => void;
}

export const RunOverview: React.FC<RunOverviewProps> = ({ run, sessions, onSelectSession }) => {
  if (!run) {
    return (
      <div style={{ padding: '3rem', textAlign: 'center', color: '#94a3b8' }}>
        <h3>No usability run loaded</h3>
        <p style={{ marginTop: '0.5rem', fontSize: '0.9rem' }}>Launch a test run or wait for active tests to complete.</p>
      </div>
    );
  }

  const completedCount = sessions.filter(s => s.goal_completed).length;
  const completionRate = sessions.length > 0 ? Math.round((completedCount / sessions.length) * 100) : 0;

  const getGradeColor = (grade: string) => {
    switch (grade) {
      case 'A': return '#10b981'; // Green
      case 'B': return '#3b82f6'; // Blue
      case 'C': return '#f59e0b'; // Amber
      case 'D': return '#f97316'; // Orange
      default: return '#ef4444';  // Red
    }
  };

  const gradeColor = getGradeColor(run.grade);

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      {/* Run Header */}
      <div style={{
        background: '#141e33',
        borderRadius: '14px',
        border: '1px solid #1e293b',
        padding: '1.75rem',
        marginBottom: '2rem',
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', flexWrap: 'wrap', gap: '1rem' }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem', flexWrap: 'wrap' }}>
              <h2 style={{ fontSize: '1.5rem', fontWeight: 800 }}>{run.name}</h2>
              <span style={{
                fontSize: '0.75rem',
                fontWeight: 700,
                textTransform: 'uppercase',
                padding: '3px 8px',
                borderRadius: '6px',
                background: run.status === 'completed' ? 'rgba(16,185,129,0.15)' : 'rgba(59,130,246,0.15)',
                color: run.status === 'completed' ? '#34d399' : '#60a5fa',
                border: `1px solid ${run.status === 'completed' ? '#059669' : '#2563eb'}`,
                whiteSpace: 'nowrap',
              }}>
                {run.status}
              </span>
            </div>
            <div style={{ fontSize: '0.9rem', color: '#94a3b8', display: 'flex', gap: '1.5rem', flexWrap: 'wrap' }}>
              <span>Target: <strong style={{ color: '#cbd5e1' }}>{run.target_url}</strong></span>
              <span>Goal: <strong style={{ color: '#cbd5e1' }}>{run.goal}</strong></span>
            </div>
          </div>

          {run.mandatory_redesign && (
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              padding: '8px 14px',
              borderRadius: '8px',
              background: 'rgba(239,68,68,0.15)',
              border: '1px solid #ef4444',
              color: '#f87171',
              fontSize: '0.85rem',
              fontWeight: 600,
            }}>
              <AlertTriangle size={18} /> Redesign Recommended (FI &gt; 2.5)
            </div>
          )}
        </div>
      </div>

      {/* Primary KPI Grid */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
        gap: '1.25rem',
        marginBottom: '2.5rem',
      }}>
        {/* SUS Score Card */}
        <div style={{
          background: '#141e33',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '1.5rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1.25rem',
        }}>
          <div style={{
            width: '64px',
            height: '64px',
            borderRadius: '50%',
            border: `4px solid ${gradeColor}`,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'rgba(0,0,0,0.2)',
          }}>
            <span style={{ fontSize: '1.3rem', fontWeight: 800, color: gradeColor }}>{run.grade || '-'}</span>
          </div>
          <div>
            <div style={{ fontSize: '2rem', fontWeight: 800, lineHeight: 1.1 }}>{run.average_sus.toFixed(1)}</div>
            <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
              System Usability Scale
            </div>
          </div>
        </div>

        {/* Friction Index */}
        <div style={{
          background: '#141e33',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '1.5rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1.25rem',
        }}>
          <div style={{
            width: '54px',
            height: '54px',
            borderRadius: '12px',
            background: run.mandatory_redesign ? 'rgba(239,68,68,0.15)' : 'rgba(16,185,129,0.15)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: run.mandatory_redesign ? '#ef4444' : '#10b981',
          }}>
            <Flame size={28} />
          </div>
          <div>
            <div style={{ fontSize: '2rem', fontWeight: 800, lineHeight: 1.1 }}>{run.average_friction.toFixed(2)}</div>
            <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
              Friction Index (FI)
            </div>
          </div>
        </div>

        {/* Conversion Rate */}
        <div style={{
          background: '#141e33',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '1.5rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1.25rem',
        }}>
          <div style={{
            width: '54px',
            height: '54px',
            borderRadius: '12px',
            background: 'rgba(59,130,246,0.15)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#3b82f6',
          }}>
            <CheckCircle2 size={28} />
          </div>
          <div>
            <div style={{ fontSize: '2rem', fontWeight: 800, lineHeight: 1.1 }}>{completionRate}%</div>
            <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
              Funnel Completion
            </div>
          </div>
        </div>

        {/* Total Personas */}
        <div style={{
          background: '#141e33',
          border: '1px solid #1e293b',
          borderRadius: '12px',
          padding: '1.5rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1.25rem',
        }}>
          <div style={{
            width: '54px',
            height: '54px',
            borderRadius: '12px',
            background: 'rgba(139,92,246,0.15)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#8b5cf6',
          }}>
            <Monitor size={28} />
          </div>
          <div>
            <div style={{ fontSize: '2rem', fontWeight: 800, lineHeight: 1.1 }}>{sessions.length}</div>
            <div style={{ fontSize: '0.8rem', color: '#94a3b8', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
              Evaluated Personas
            </div>
          </div>
        </div>
      </div>

      {/* Swarm Persona Sessions Table */}
      <div style={{
        background: '#141e33',
        border: '1px solid #1e293b',
        borderRadius: '14px',
        overflow: 'hidden',
      }}>
        <div style={{
          padding: '1.25rem 1.5rem',
          borderBottom: '1px solid #1e293b',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}>
          <h3 style={{ fontSize: '1.1rem', fontWeight: 700 }}>Synthetic Persona Session Breakdown</h3>
          <span style={{ fontSize: '0.8rem', color: '#94a3b8' }}>Click any row to replay session</span>
        </div>

        <div style={{ overflowX: 'auto', WebkitOverflowScrolling: 'touch' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '0.9rem', minWidth: '700px' }}>
            <thead>
              <tr style={{ background: '#090d16', borderBottom: '1px solid #1e293b', color: '#94a3b8', textTransform: 'uppercase', fontSize: '0.75rem', letterSpacing: '0.05em' }}>
                <th style={{ padding: '0.85rem 1.25rem', whiteSpace: 'nowrap' }}>Persona & Role</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>Device</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>SUS Score</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>Friction (FI)</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>Hesitations</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>Dead Clicks</th>
                <th style={{ padding: '0.85rem 1rem', whiteSpace: 'nowrap' }}>Outcome</th>
                <th style={{ padding: '0.85rem 1.25rem', textAlign: 'right', whiteSpace: 'nowrap' }}>Action</th>
              </tr>
            </thead>
            <tbody>
              {sessions.map((sess) => {
              const color = getGradeColor(sess.sus_grade);
              return (
                <tr
                  key={sess.id}
                  style={{
                    borderBottom: '1px solid #1e293b',
                    transition: 'background 0.15s ease',
                    cursor: 'pointer',
                  }}
                  onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#1a2744')}
                  onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                  onClick={() => onSelectSession(sess.id)}
                >
                  <td style={{ padding: '1rem 1.25rem' }}>
                    <div style={{ fontWeight: 700, color: '#f8fafc' }}>{sess.persona_name}</div>
                    <div style={{ fontSize: '0.75rem', color: '#94a3b8' }}>{sess.archetype}</div>
                  </td>

                  <td style={{ padding: '1rem' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', color: '#cbd5e1' }}>
                      {sess.is_mobile ? <Smartphone size={16} /> : <Monitor size={16} />}
                      <span style={{ fontSize: '0.8rem' }}>{sess.is_mobile ? '390x844' : '1440x900'}</span>
                    </div>
                  </td>

                  <td style={{ padding: '1rem' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                      <span style={{ fontWeight: 800, color }}>{sess.sus_score.toFixed(1)}</span>
                      <span style={{
                        fontSize: '0.7rem',
                        fontWeight: 700,
                        padding: '1px 6px',
                        borderRadius: '4px',
                        background: `${color}20`,
                        color,
                      }}>
                        Grade {sess.sus_grade}
                      </span>
                    </div>
                  </td>

                  <td style={{ padding: '1rem', fontWeight: 600, color: sess.friction_index > 2.5 ? '#f87171' : '#f8fafc' }}>
                    {sess.friction_index.toFixed(2)}
                  </td>

                  <td style={{ padding: '1rem', color: sess.hesitation_count > 0 ? '#f59e0b' : '#94a3b8' }}>
                    {sess.hesitation_count} {sess.hesitation_count === 1 ? 'pause' : 'pauses'}
                  </td>

                  <td style={{ padding: '1rem', color: sess.dead_clicks > 0 ? '#ef4444' : '#94a3b8' }}>
                    {sess.dead_clicks}
                  </td>

                  <td style={{ padding: '1rem' }}>
                    {sess.goal_completed ? (
                      <span style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.35rem',
                        fontSize: '0.75rem',
                        fontWeight: 700,
                        color: '#34d399',
                        background: 'rgba(16,185,129,0.1)',
                        padding: '3px 8px',
                        borderRadius: '6px',
                      }}>
                        <CheckCircle2 size={13} /> Completed
                      </span>
                    ) : (
                      <span style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.35rem',
                        fontSize: '0.75rem',
                        fontWeight: 700,
                        color: '#f87171',
                        background: 'rgba(239,68,68,0.1)',
                        padding: '3px 8px',
                        borderRadius: '6px',
                      }}>
                        <XCircle size={13} /> Abandoned
                      </span>
                    )}
                  </td>

                  <td style={{ padding: '1rem 1.25rem', textAlign: 'right' }}>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onSelectSession(sess.id);
                      }}
                      style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '0.35rem',
                        background: '#2563eb',
                        color: 'white',
                        border: 'none',
                        borderRadius: '6px',
                        padding: '6px 12px',
                        fontSize: '0.8rem',
                        fontWeight: 600,
                        cursor: 'pointer',
                      }}
                    >
                      <Play size={13} /> Replay
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        </div>
      </div>
    </div>
  );
};
