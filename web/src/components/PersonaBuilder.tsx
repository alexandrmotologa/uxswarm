import React, { useState, useEffect } from 'react';
import { Persona } from '../types';
import { Users, Smartphone, Monitor, Shield, Zap, Eye, DollarSign } from 'lucide-react';

export const PersonaBuilder: React.FC = () => {
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [selectedPersona, setSelectedPersona] = useState<Persona | null>(null);

  useEffect(() => {
    fetch('/api/personas')
      .then((res) => res.json())
      .then((data: Persona[]) => {
        setPersonas(data || []);
        if (data && data.length > 0) {
          setSelectedPersona(data[0]);
        }
      })
      .catch((err) => console.error('Failed to load personas:', err));
  }, []);

  const getArchetypeIcon = (archetype: string) => {
    switch (archetype) {
      case 'HIGH_IMPATIENCE_LOW_TECH':
        return <Zap size={18} color="#f59e0b" />;
      case 'SCRUTINIZER_HIGH_TECH':
        return <Shield size={18} color="#3b82f6" />;
      case 'MOBILE_THUMB_DRIVEN':
        return <Smartphone size={18} color="#a855f7" />;
      case 'ACCESSIBILITY_RELIANT':
        return <Eye size={18} color="#10b981" />;
      default:
        return <DollarSign size={18} color="#ec4899" />;
    }
  };

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      <div style={{ marginBottom: '2rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', marginBottom: '0.35rem' }}>
          <Users size={24} color="#3b82f6" />
          <h2 style={{ fontSize: '1.4rem', fontWeight: 800 }}>Synthetic Persona Catalog</h2>
        </div>
        <p style={{ fontSize: '0.85rem', color: '#94a3b8' }}>
          Configurable cognitive agents modeling human patience limits, reading velocity, and domain expectations.
        </p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '360px 1fr', gap: '2rem' }}>
        {/* Left List of Personas */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          {personas.map((p) => {
            const isSelected = selectedPersona?.id === p.id;
            return (
              <div
                key={p.id}
                onClick={() => setSelectedPersona(p)}
                style={{
                  background: isSelected ? '#1e293b' : '#141e33',
                  border: `1px solid ${isSelected ? '#3b82f6' : '#1e293b'}`,
                  borderRadius: '12px',
                  padding: '1.25rem',
                  cursor: 'pointer',
                  transition: 'all 0.15s ease',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem' }}>
                  <div style={{
                    width: '36px',
                    height: '36px',
                    borderRadius: '8px',
                    background: '#090d16',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}>
                    {getArchetypeIcon(p.archetype)}
                  </div>
                  <div>
                    <div style={{ fontWeight: 700, fontSize: '0.95rem', color: '#f8fafc' }}>{p.name}</div>
                    <div style={{ fontSize: '0.75rem', color: '#94a3b8' }}>Age {p.age} &bull; {p.archetype}</div>
                  </div>
                </div>

                <div style={{ display: 'flex', gap: '0.75rem', fontSize: '0.75rem', color: '#cbd5e1' }}>
                  <span>Patience: <strong>{p.patience_threshold_seconds}s</strong></span>
                  <span>Tech: <strong>{p.technical_literacy}/10</strong></span>
                  <span>Device: <strong>{p.viewport.is_mobile ? 'Mobile' : 'Desktop'}</strong></span>
                </div>
              </div>
            );
          })}
        </div>

        {/* Right Detail Pane */}
        {selectedPersona && (
          <div style={{
            background: '#141e33',
            border: '1px solid #1e293b',
            borderRadius: '14px',
            padding: '2rem',
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.5rem', borderBottom: '1px solid #1e293b', paddingBottom: '1.25rem' }}>
              <div>
                <h3 style={{ fontSize: '1.3rem', fontWeight: 800 }}>{selectedPersona.name}</h3>
                <div style={{ fontSize: '0.85rem', color: '#94a3b8', marginTop: '0.2rem' }}>
                  Archetype: <span style={{ color: '#60a5fa', fontWeight: 600 }}>{selectedPersona.archetype}</span>
                </div>
              </div>

              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.4rem',
                fontSize: '0.8rem',
                color: '#cbd5e1',
                background: '#090d16',
                padding: '4px 10px',
                borderRadius: '6px',
                border: '1px solid #1e293b',
              }}>
                {selectedPersona.viewport.is_mobile ? <Smartphone size={15} /> : <Monitor size={15} />}
                <span>{selectedPersona.viewport.width} x {selectedPersona.viewport.height} ({selectedPersona.viewport.is_mobile ? 'Mobile' : 'Desktop'})</span>
              </div>
            </div>

            {/* Parameter Gauges */}
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.25rem', marginBottom: '2rem' }}>
              <div style={{ background: '#090d16', border: '1px solid #1e293b', borderRadius: '10px', padding: '1rem' }}>
                <div style={{ fontSize: '0.75rem', color: '#94a3b8', textTransform: 'uppercase' }}>Patience Threshold</div>
                <div style={{ fontSize: '1.4rem', fontWeight: 800, marginTop: '0.25rem', color: selectedPersona.patience_threshold_seconds <= 4 ? '#f87171' : '#34d399' }}>
                  {selectedPersona.patience_threshold_seconds} seconds
                </div>
                <div style={{ fontSize: '0.75rem', color: '#64748b', marginTop: '0.2rem' }}>Abandons flow if friction exceeds this latency.</div>
              </div>

              <div style={{ background: '#090d16', border: '1px solid #1e293b', borderRadius: '10px', padding: '1rem' }}>
                <div style={{ fontSize: '0.75rem', color: '#94a3b8', textTransform: 'uppercase' }}>Technical Literacy</div>
                <div style={{ fontSize: '1.4rem', fontWeight: 800, marginTop: '0.25rem', color: '#38bdf8' }}>
                  {selectedPersona.technical_literacy} / 10
                </div>
                <div style={{ fontSize: '0.75rem', color: '#64748b', marginTop: '0.2rem' }}>Dictates reading velocity and error recovery.</div>
              </div>

              <div style={{ background: '#090d16', border: '1px solid #1e293b', borderRadius: '10px', padding: '1rem' }}>
                <div style={{ fontSize: '0.75rem', color: '#94a3b8', textTransform: 'uppercase' }}>Attention Span</div>
                <div style={{ fontSize: '1.4rem', fontWeight: 800, marginTop: '0.25rem', color: '#c084fc' }}>
                  {selectedPersona.attention_span}
                </div>
                <div style={{ fontSize: '0.75rem', color: '#64748b', marginTop: '0.2rem' }}>Determines likelihood of text skimming.</div>
              </div>
            </div>

            {/* Behavioral Traits */}
            <div>
              <h4 style={{ fontSize: '0.95rem', fontWeight: 700, marginBottom: '0.75rem' }}>Cognitive Behavioral Traits</h4>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                {selectedPersona.traits.map((trait, tIdx) => (
                  <div
                    key={tIdx}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: '0.6rem',
                      background: '#090d16',
                      padding: '8px 12px',
                      borderRadius: '8px',
                      border: '1px solid #1e293b',
                      fontSize: '0.85rem',
                      color: '#e2e8f0',
                    }}
                  >
                    <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: '#3b82f6' }} />
                    {trait}
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
