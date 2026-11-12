import React, { useState, useEffect } from 'react';
import { SessionRecord, StepRecord } from '../types';
import { Flame, AlertCircle, Clock, Filter, Eye } from 'lucide-react';

interface FrictionHeatmapProps {
  sessions: SessionRecord[];
}

interface GazePoint {
  x: number;
  y: number;
  intensity: number;
  duration: number;
  target?: string;
}

export const FrictionHeatmap: React.FC<FrictionHeatmapProps> = ({ sessions }) => {
  const [allSteps, setAllSteps] = useState<{ step: StepRecord; personaName: string }[]>([]);
  const [filterType, setFilterType] = useState<'all' | 'dead_clicks' | 'hesitations' | 'gaze'>('all');
  const [hoveredPoint, setHoveredPoint] = useState<{ step: StepRecord; personaName: string } | null>(null);
  const [gazePoints, setGazePoints] = useState<GazePoint[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (sessions.length === 0) return;
    setLoading(true);

    Promise.all(
      sessions.map((s) =>
        fetch(`/api/sessions/${s.id}/steps`)
          .then((res) => res.json())
          .then((steps: StepRecord[]) =>
            (steps || []).map((step) => ({ step, personaName: s.persona_name }))
          )
      )
    )
      .then((results) => {
        const flattened = results.flat();
        setAllSteps(flattened);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to load steps for heatmap:', err);
        setLoading(false);
      });

    // Fetch simulated gaze map
    if (sessions[0]?.run_id) {
      fetch(`/api/runs/${sessions[0].run_id}/gaze`)
        .then((res) => res.json())
        .then((data: GazePoint[]) => setGazePoints(data || []))
        .catch(() => {});
    }
  }, [sessions]);

  const filteredPoints = allSteps.filter(({ step }) => {
    if (step.x <= 0 && step.y <= 0) return false;
    if (filterType === 'dead_clicks') return step.is_dead_click;
    if (filterType === 'hesitations') return step.hesitation_seconds >= 3.5;
    return true;
  });

  const deadClickCount = allSteps.filter(({ step }) => step.is_dead_click).length;
  const hesitationCount = allSteps.filter(({ step }) => step.hesitation_seconds >= 3.5).length;

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      {/* Header and Filter Controls */}
      <div style={{
        background: '#141e33',
        borderRadius: '14px',
        border: '1px solid #1e293b',
        padding: '1.5rem',
        marginBottom: '2rem',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: '1rem',
      }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', marginBottom: '0.35rem' }}>
            <Flame size={24} color="#f59e0b" />
            <h2 style={{ fontSize: '1.4rem', fontWeight: 800 }}>Friction Coordinates & Gaze Heatmap</h2>
          </div>
          <p style={{ fontSize: '0.85rem', color: '#94a3b8' }}>
            Visual distribution of clicks, misclick traps, hesitation points, and simulated eye-tracking fixations.
          </p>
        </div>

        {/* Filter buttons */}
        <div style={{ display: 'flex', gap: '0.5rem', background: '#090d16', padding: '4px', borderRadius: '10px', border: '1px solid #1e293b' }}>
          <button
            onClick={() => setFilterType('all')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.35rem',
              padding: '6px 12px',
              borderRadius: '6px',
              fontSize: '0.8rem',
              fontWeight: 600,
              border: 'none',
              cursor: 'pointer',
              background: filterType === 'all' ? '#2563eb' : 'transparent',
              color: filterType === 'all' ? 'white' : '#94a3b8',
            }}
          >
            <Filter size={14} /> All Clicks ({allSteps.filter(({ step }) => step.x > 0).length})
          </button>

          <button
            onClick={() => setFilterType('dead_clicks')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.35rem',
              padding: '6px 12px',
              borderRadius: '6px',
              fontSize: '0.8rem',
              fontWeight: 600,
              border: 'none',
              cursor: 'pointer',
              background: filterType === 'dead_clicks' ? '#ef4444' : 'transparent',
              color: filterType === 'dead_clicks' ? 'white' : '#94a3b8',
            }}
          >
            <AlertCircle size={14} /> Dead Clicks ({deadClickCount})
          </button>

          <button
            onClick={() => setFilterType('hesitations')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.35rem',
              padding: '6px 12px',
              borderRadius: '6px',
              fontSize: '0.8rem',
              fontWeight: 600,
              border: 'none',
              cursor: 'pointer',
              background: filterType === 'hesitations' ? '#f59e0b' : 'transparent',
              color: filterType === 'hesitations' ? 'white' : '#94a3b8',
            }}
          >
            <Clock size={14} /> Hesitation Hotspots ({hesitationCount})
          </button>

          <button
            onClick={() => setFilterType('gaze')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.35rem',
              padding: '6px 12px',
              borderRadius: '6px',
              fontSize: '0.8rem',
              fontWeight: 600,
              border: 'none',
              cursor: 'pointer',
              background: filterType === 'gaze' ? '#8b5cf6' : 'transparent',
              color: filterType === 'gaze' ? 'white' : '#94a3b8',
            }}
          >
            <Eye size={14} /> Eye-Tracking Gaze
          </button>
        </div>
      </div>

      {loading ? (
        <div style={{ padding: '4rem', textAlign: 'center', color: '#94a3b8' }}>
          Loading interaction coordinates...
        </div>
      ) : (
        <div style={{ position: 'relative' }}>
          {/* Heatmap Canvas Container */}
          <div style={{
            background: '#090d16',
            borderRadius: '14px',
            border: '2px solid #1e293b',
            height: '650px',
            position: 'relative',
            overflow: 'hidden',
            boxShadow: 'inset 0 0 40px rgba(0,0,0,0.8)',
          }}>
            {/* Grid Lines */}
            <div style={{
              position: 'absolute',
              inset: 0,
              backgroundImage: 'radial-gradient(#1e293b 1px, transparent 1px)',
              backgroundSize: '24px 24px',
              opacity: 0.6,
            }} />

            {/* Simulated Eye-Tracking Gaze Heatmap Overlay */}
            {filterType === 'gaze' && (
              <div style={{ position: 'absolute', inset: 0, pointerEvents: 'none' }}>
                {gazePoints.map((gp, idx) => (
                  <div
                    key={idx}
                    style={{
                      position: 'absolute',
                      left: `${Math.min(1050, Math.max(30, gp.x))}px`,
                      top: `${Math.min(600, Math.max(30, gp.y))}px`,
                      width: `${Math.round(gp.intensity * 90)}px`,
                      height: `${Math.round(gp.intensity * 90)}px`,
                      borderRadius: '50%',
                      background: `radial-gradient(circle, rgba(239, 68, 68, ${gp.intensity * 0.75}) 0%, rgba(245, 158, 11, ${gp.intensity * 0.45}) 40%, rgba(59, 130, 246, ${gp.intensity * 0.15}) 70%, transparent 100%)`,
                      transform: 'translate(-50%, -50%)',
                      filter: 'blur(10px)',
                    }}
                  />
                ))}
              </div>
            )}

            {/* Coordinates markers */}
            {filterType !== 'gaze' && filteredPoints.map(({ step, personaName }, idx) => {
              const isDead = step.is_dead_click;
              const isHesitation = step.hesitation_seconds >= 3.5;

              let color = '#3b82f6';
              if (isDead) color = '#ef4444';
              else if (isHesitation) color = '#f59e0b';

              const posX = Math.min(1050, Math.max(30, step.x));
              const posY = Math.min(600, Math.max(30, step.y));

              return (
                <div
                  key={idx}
                  onMouseEnter={() => setHoveredPoint({ step, personaName })}
                  onMouseLeave={() => setHoveredPoint(null)}
                  style={{
                    position: 'absolute',
                    left: `${posX}px`,
                    top: `${posY}px`,
                    width: isDead ? '20px' : '14px',
                    height: isDead ? '20px' : '14px',
                    borderRadius: '50%',
                    backgroundColor: color,
                    border: '2px solid rgba(255,255,255,0.8)',
                    transform: 'translate(-50%, -50%)',
                    cursor: 'pointer',
                    boxShadow: `0 0 16px ${color}`,
                    zIndex: isDead ? 10 : 5,
                    transition: 'transform 0.15s ease',
                  }}
                />
              );
            })}

            {/* Hover tooltip card */}
            {hoveredPoint && filterType !== 'gaze' && (
              <div style={{
                position: 'absolute',
                left: `${Math.min(900, Math.max(40, hoveredPoint.step.x + 20))}px`,
                top: `${Math.min(520, Math.max(30, hoveredPoint.step.y - 10))}px`,
                background: '#141e33',
                border: '1px solid #3b82f6',
                borderRadius: '8px',
                padding: '10px 14px',
                zIndex: 30,
                pointerEvents: 'none',
                boxShadow: '0 10px 15px -3px rgba(0,0,0,0.5)',
                minWidth: '220px',
              }}>
                <div style={{ fontWeight: 700, fontSize: '0.85rem', color: '#f8fafc' }}>
                  {hoveredPoint.personaName}
                </div>
                <div style={{ fontSize: '0.75rem', color: '#cbd5e1', marginTop: '3px' }}>
                  Target: {hoveredPoint.step.target_description || 'Coordinates'}
                </div>
                <div style={{ fontSize: '0.75rem', color: '#94a3b8', marginTop: '2px' }}>
                  Coords: ({hoveredPoint.step.x}, {hoveredPoint.step.y})
                </div>
                {hoveredPoint.step.hesitation_seconds > 0 && (
                  <div style={{ fontSize: '0.75rem', color: '#f59e0b', marginTop: '2px', fontWeight: 600 }}>
                    Hesitation: {hoveredPoint.step.hesitation_seconds.toFixed(1)}s
                  </div>
                )}
                {hoveredPoint.step.is_dead_click && (
                  <div style={{ fontSize: '0.75rem', color: '#f87171', marginTop: '2px', fontWeight: 700 }}>
                    &#x26A0; Dead Click (Non-interactive)
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Legend */}
          <div style={{
            display: 'flex',
            gap: '1.5rem',
            justifyContent: 'center',
            marginTop: '1.25rem',
            fontSize: '0.85rem',
            color: '#94a3b8',
          }}>
            {filterType === 'gaze' ? (
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#ef4444' }} />
                <span>Simulated F-Pattern & Visual Saliency Fixations (Tobii Eye-Tracking Model)</span>
              </div>
            ) : (
              <>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <span style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#ef4444' }} />
                  <span>Dead Click (Misclick Trap)</span>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <span style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#f59e0b' }} />
                  <span>Hesitation Hotspot (&gt; 3.5s)</span>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <span style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#3b82f6' }} />
                  <span>Standard Interaction</span>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
