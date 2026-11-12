import React, { useState, useEffect } from 'react';
import { SessionRecord, StepRecord } from '../types';
import { Play, Pause, SkipBack, SkipForward, Clock, AlertCircle, MousePointer, Smartphone, Monitor } from 'lucide-react';

interface SessionReplayProps {
  sessions: SessionRecord[];
  selectedSessionId: string | null;
  onSelectSession: (id: string) => void;
}

export const SessionReplay: React.FC<SessionReplayProps> = ({
  sessions,
  selectedSessionId,
  onSelectSession,
}) => {
  const currentSession = sessions.find((s) => s.id === selectedSessionId) || sessions[0];
  const [steps, setSteps] = useState<StepRecord[]>([]);
  const [currentStepIndex, setCurrentStepIndex] = useState<number>(0);
  const [isPlaying, setIsPlaying] = useState<boolean>(false);
  const [playbackSpeed, setPlaybackSpeed] = useState<number>(1);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (!currentSession) return;
    setLoading(true);
    fetch(`/api/sessions/${currentSession.id}/steps`)
      .then((res) => res.json())
      .then((data: StepRecord[]) => {
        setSteps(data || []);
        setCurrentStepIndex(0);
        setIsPlaying(false);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to load session steps:', err);
        setLoading(false);
      });
  }, [currentSession?.id]);

  // Autoplay ticker
  useEffect(() => {
    if (!isPlaying || steps.length === 0) return;

    const interval = setInterval(() => {
      setCurrentStepIndex((prev) => {
        if (prev >= steps.length - 1) {
          setIsPlaying(false);
          return prev;
        }
        return prev + 1;
      });
    }, 1500 / playbackSpeed);

    return () => clearInterval(interval);
  }, [isPlaying, steps.length, playbackSpeed]);

  if (!currentSession) {
    return (
      <div style={{ padding: '3rem', textAlign: 'center', color: '#94a3b8' }}>
        No sessions available to replay.
      </div>
    );
  }

  const currentStep = steps[currentStepIndex];

  return (
    <div style={{ maxWidth: '1300px', margin: '0 auto', padding: '1.5rem 2rem' }}>
      {/* Session Tab Selector */}
      <div style={{
        display: 'flex',
        gap: '0.75rem',
        overflowX: 'auto',
        paddingBottom: '1rem',
        marginBottom: '1.5rem',
        borderBottom: '1px solid #1e293b',
      }}>
        {sessions.map((s) => {
          const isSelected = s.id === currentSession.id;
          return (
            <button
              key={s.id}
              onClick={() => onSelectSession(s.id)}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
                padding: '8px 16px',
                borderRadius: '10px',
                background: isSelected ? '#1e293b' : '#141e33',
                border: `1px solid ${isSelected ? '#3b82f6' : '#1e293b'}`,
                color: isSelected ? '#ffffff' : '#94a3b8',
                cursor: 'pointer',
                whiteSpace: 'nowrap',
                transition: 'all 0.15s ease',
              }}
            >
              {s.is_mobile ? <Smartphone size={16} /> : <Monitor size={16} />}
              <span style={{ fontWeight: 600, fontSize: '0.85rem' }}>{s.persona_name}</span>
              <span style={{
                fontSize: '0.75rem',
                padding: '2px 6px',
                borderRadius: '4px',
                background: s.sus_score >= 80 ? 'rgba(16,185,129,0.2)' : 'rgba(239,68,68,0.2)',
                color: s.sus_score >= 80 ? '#34d399' : '#f87171',
                fontWeight: 700,
              }}>
                SUS {s.sus_score.toFixed(0)}
              </span>
            </button>
          );
        })}
      </div>

      {loading ? (
        <div style={{ padding: '4rem', textAlign: 'center', color: '#94a3b8' }}>
          Loading session telemetry and frames...
        </div>
      ) : (
        <div className="session-replay-grid">
          {/* Main Visual Player */}
          <div>
            {/* Virtual Device Bezel */}
            <div style={{
              background: '#090d16',
              borderRadius: '14px',
              border: '2px solid #1e293b',
              overflow: 'hidden',
              boxShadow: '0 20px 25px -5px rgba(0,0,0,0.5)',
              position: 'relative',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
            }}>
              {/* Browser / Device Chrome Header */}
              <div style={{
                width: '100%',
                background: '#141e33',
                padding: '10px 16px',
                display: 'flex',
                alignItems: 'center',
                gap: '0.75rem',
                borderBottom: '1px solid #1e293b',
              }}>
                <div style={{ display: 'flex', gap: '6px' }}>
                  <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: '#ef4444' }} />
                  <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: '#f59e0b' }} />
                  <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: '#10b981' }} />
                </div>
                <div style={{
                  flex: 1,
                  background: '#090d16',
                  borderRadius: '6px',
                  padding: '4px 12px',
                  fontSize: '0.75rem',
                  fontFamily: 'monospace',
                  color: '#94a3b8',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}>
                  {currentStep?.url || currentSession.run_id}
                </div>
              </div>

              {/* Viewport Frame Canvas */}
              <div style={{
                position: 'relative',
                width: '100%',
                minHeight: '440px',
                maxHeight: '620px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: '#020617',
                overflow: 'hidden',
              }}>
                {currentStep?.frame_data_uri ? (
                  <img
                    src={currentStep.frame_data_uri}
                    alt={`Step ${currentStep.step_index}`}
                    style={{
                      maxWidth: '100%',
                      maxHeight: '620px',
                      objectFit: 'contain',
                      display: 'block',
                    }}
                  />
                ) : (
                  <div style={{ padding: '3rem', textAlign: 'center', color: '#64748b' }}>
                    <div style={{ fontSize: '1.1rem', fontWeight: 600 }}>{currentStep?.target_description || 'Active Navigation State'}</div>
                    <div style={{ fontSize: '0.85rem', marginTop: '0.5rem' }}>Action: {currentStep?.action_type}</div>
                  </div>
                )}

                {/* Simulated Cursor and Click Ripple Overlay */}
                {currentStep && currentStep.x > 0 && currentStep.y > 0 && (
                  <div style={{
                    position: 'absolute',
                    left: `${currentStep.x}px`,
                    top: `${currentStep.y}px`,
                    pointerEvents: 'none',
                    transform: 'translate(-4px, -4px)',
                    zIndex: 20,
                  }}>
                    <MousePointer
                      size={24}
                      color={currentStep.is_dead_click ? '#ef4444' : '#3b82f6'}
                      fill={currentStep.is_dead_click ? '#ef4444' : '#3b82f6'}
                      style={{ filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.8))' }}
                    />
                    {currentStep.is_dead_click && (
                      <div
                        className="click-ripple"
                        style={{
                          width: '40px',
                          height: '40px',
                          border: '2px solid #ef4444',
                          background: 'rgba(239,68,68,0.3)',
                        }}
                      />
                    )}
                  </div>
                )}
              </div>

              {/* Persona Inner Monologue Subtitle Box */}
              <div style={{
                width: '100%',
                background: 'linear-gradient(180deg, rgba(20,30,51,0.95), #0f172a)',
                borderTop: '1px solid #1e293b',
                padding: '1.25rem 1.5rem',
              }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <span style={{ fontSize: '0.75rem', fontWeight: 800, color: '#60a5fa', textTransform: 'uppercase' }}>
                      {currentSession.persona_name} &bull; Thought Stream
                    </span>
                    {currentStep?.hesitation_seconds && currentStep.hesitation_seconds >= 3.5 ? (
                      <span style={{
                        fontSize: '0.7rem',
                        fontWeight: 700,
                        color: '#f59e0b',
                        background: 'rgba(245,158,11,0.15)',
                        padding: '2px 6px',
                        borderRadius: '4px',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '3px',
                      }}>
                        <Clock size={11} /> Paused {currentStep.hesitation_seconds}s
                      </span>
                    ) : null}
                    {currentStep?.is_dead_click && (
                      <span style={{
                        fontSize: '0.7rem',
                        fontWeight: 700,
                        color: '#ef4444',
                        background: 'rgba(239,68,68,0.15)',
                        padding: '2px 6px',
                        borderRadius: '4px',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '3px',
                      }}>
                        <AlertCircle size={11} /> Dead Click
                      </span>
                    )}
                  </div>

                  <div style={{ fontSize: '0.75rem', color: '#94a3b8' }}>
                    Step {currentStepIndex + 1} of {steps.length || 1}
                  </div>
                </div>

                <div style={{ fontSize: '0.95rem', color: '#f1f5f9', fontStyle: 'italic', lineHeight: 1.4 }}>
                  &ldquo;{currentStep?.inner_speech || 'Observing page structure...'}&rdquo;
                </div>
              </div>
            </div>

            {/* Playback Controls Bar */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '1rem',
              background: '#141e33',
              border: '1px solid #1e293b',
              borderRadius: '12px',
              padding: '0.85rem 1.25rem',
              marginTop: '1.25rem',
            }}>
              <button
                onClick={() => setCurrentStepIndex((p) => Math.max(0, p - 1))}
                disabled={currentStepIndex === 0}
                style={{ background: 'transparent', border: 'none', color: '#cbd5e1', cursor: 'pointer' }}
              >
                <SkipBack size={18} />
              </button>

              <button
                onClick={() => setIsPlaying(!isPlaying)}
                style={{
                  background: '#2563eb',
                  border: 'none',
                  borderRadius: '50%',
                  width: '36px',
                  height: '36px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: 'white',
                  cursor: 'pointer',
                }}
              >
                {isPlaying ? <Pause size={18} /> : <Play size={18} style={{ marginLeft: '2px' }} />}
              </button>

              <button
                onClick={() => setCurrentStepIndex((p) => Math.min(steps.length - 1, p + 1))}
                disabled={currentStepIndex >= steps.length - 1}
                style={{ background: 'transparent', border: 'none', color: '#cbd5e1', cursor: 'pointer' }}
              >
                <SkipForward size={18} />
              </button>

              {/* Scrubber slider */}
              <input
                type="range"
                min={0}
                max={Math.max(0, steps.length - 1)}
                value={currentStepIndex}
                onChange={(e) => {
                  setCurrentStepIndex(Number(e.target.value));
                  setIsPlaying(false);
                }}
                style={{ flex: 1, cursor: 'pointer', accentColor: '#2563eb' }}
              />

              {/* Speed Selector */}
              <div style={{ display: 'flex', gap: '4px' }}>
                {[1, 2, 4].map((spd) => (
                  <button
                    key={spd}
                    onClick={() => setPlaybackSpeed(spd)}
                    style={{
                      background: playbackSpeed === spd ? '#2563eb' : '#090d16',
                      border: '1px solid #1e293b',
                      borderRadius: '6px',
                      padding: '3px 8px',
                      fontSize: '0.75rem',
                      fontWeight: 600,
                      color: playbackSpeed === spd ? 'white' : '#94a3b8',
                      cursor: 'pointer',
                    }}
                  >
                    {spd}x
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Right Sidebar: Step Details & Friction Log */}
          <div>
            <div style={{
              background: '#141e33',
              borderRadius: '14px',
              border: '1px solid #1e293b',
              padding: '1.25rem',
              height: '100%',
              display: 'flex',
              flexDirection: 'column',
            }}>
              <h4 style={{ fontSize: '0.95rem', fontWeight: 700, marginBottom: '1rem', color: '#f8fafc' }}>
                Step Interaction Telemetry
              </h4>

              {currentStep ? (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem', fontSize: '0.85rem' }}>
                  <div>
                    <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>Action Type</div>
                    <div style={{ fontWeight: 700, color: '#f8fafc', textTransform: 'capitalize' }}>
                      {currentStep.action_type}
                    </div>
                  </div>

                  <div>
                    <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>Target Element</div>
                    <div style={{ fontWeight: 600, color: '#cbd5e1' }}>
                      {currentStep.target_description || 'None / Page'}
                    </div>
                  </div>

                  {currentStep.target_selector && (
                    <div>
                      <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>CSS Selector</div>
                      <code style={{ fontSize: '0.75rem', color: '#93c5fd', background: '#090d16', padding: '3px 6px', borderRadius: '4px', display: 'block', overflowWrap: 'break-word' }}>
                        {currentStep.target_selector}
                      </code>
                    </div>
                  )}

                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' }}>
                    <div>
                      <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>Hesitation</div>
                      <div style={{ fontWeight: 700, color: currentStep.hesitation_seconds >= 3.5 ? '#f59e0b' : '#f8fafc' }}>
                        {currentStep.hesitation_seconds.toFixed(1)}s
                      </div>
                    </div>

                    <div>
                      <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase' }}>Frustration Index</div>
                      <div style={{ fontWeight: 700, color: currentStep.frustration > 0.5 ? '#ef4444' : '#10b981' }}>
                        {(currentStep.frustration * 100).toFixed(0)}%
                      </div>
                    </div>
                  </div>

                  {/* Frustration bar */}
                  <div>
                    <div style={{ color: '#94a3b8', fontSize: '0.75rem', textTransform: 'uppercase', marginBottom: '0.35rem' }}>
                      Frustration Level
                    </div>
                    <div style={{ width: '100%', height: '8px', background: '#090d16', borderRadius: '4px', overflow: 'hidden' }}>
                      <div style={{
                        width: `${Math.min(100, Math.round(currentStep.frustration * 100))}%`,
                        height: '100%',
                        background: currentStep.frustration > 0.7 ? '#ef4444' : currentStep.frustration > 0.4 ? '#f59e0b' : '#3b82f6',
                        transition: 'width 0.2s',
                      }} />
                    </div>
                  </div>
                </div>
              ) : (
                <div style={{ color: '#94a3b8' }}>Select a step to inspect telemetry.</div>
              )}

              {/* Step Sequence Timeline List */}
              <div style={{ marginTop: '1.5rem', flex: 1, overflowY: 'auto', maxHeight: '280px' }}>
                <div style={{ fontSize: '0.75rem', color: '#94a3b8', textTransform: 'uppercase', marginBottom: '0.5rem', fontWeight: 700 }}>
                  Session Steps
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.35rem' }}>
                  {steps.map((st, i) => (
                    <button
                      key={st.id || i}
                      onClick={() => {
                        setCurrentStepIndex(i);
                        setIsPlaying(false);
                      }}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        padding: '6px 10px',
                        borderRadius: '6px',
                        background: i === currentStepIndex ? '#2563eb' : '#090d16',
                        border: '1px solid #1e293b',
                        color: i === currentStepIndex ? 'white' : '#cbd5e1',
                        fontSize: '0.75rem',
                        cursor: 'pointer',
                        textAlign: 'left',
                      }}
                    >
                      <span>#{i + 1} {st.action_type}</span>
                      {st.is_dead_click && <span style={{ color: '#f87171', fontWeight: 700 }}>Dead Click</span>}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
