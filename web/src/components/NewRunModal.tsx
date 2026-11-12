import React, { useState } from 'react';
import { X, Play, Sparkles } from 'lucide-react';
import { RunRecord } from '../types';

interface NewRunModalProps {
  isOpen: boolean;
  onClose: () => void;
  onRunCreated: (run: RunRecord) => void;
}

export const NewRunModal: React.FC<NewRunModalProps> = ({ isOpen, onClose, onRunCreated }) => {
  const [name, setName] = useState<string>('B2B SaaS Onboarding Study');
  const [targetUrl, setTargetUrl] = useState<string>('http://localhost:3000/signup');
  const [goal, setGoal] = useState<string>('Sign up for an account, configure workspace, and complete onboarding');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const res = await fetch('/api/runs', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, target_url: targetUrl, goal }),
      });

      if (!res.ok) {
        throw new Error(`Failed with status ${res.status}`);
      }

      const run: RunRecord = await res.json();
      onRunCreated(run);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to launch study');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      position: 'fixed',
      inset: 0,
      backgroundColor: 'rgba(0,0,0,0.7)',
      backdropFilter: 'blur(4px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 50,
      padding: '1rem',
    }}>
      <div style={{
        background: '#141e33',
        border: '1px solid #1e293b',
        borderRadius: '16px',
        maxWidth: '520px',
        width: '100%',
        padding: '2rem',
        boxShadow: '0 25px 50px -12px rgba(0,0,0,0.5)',
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
          <h3 style={{ fontSize: '1.25rem', fontWeight: 800 }}>Launch Swarm Usability Study</h3>
          <button
            onClick={onClose}
            style={{ background: 'transparent', border: 'none', color: '#94a3b8', cursor: 'pointer' }}
          >
            <X size={20} />
          </button>
        </div>

        {error && (
          <div style={{ background: 'rgba(239,68,68,0.15)', border: '1px solid #ef4444', color: '#f87171', padding: '8px 12px', borderRadius: '8px', fontSize: '0.85rem', marginBottom: '1.25rem' }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '1.25rem' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#cbd5e1', marginBottom: '0.4rem' }}>
              Study Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '10px 12px',
                borderRadius: '8px',
                background: '#090d16',
                border: '1px solid #1e293b',
                color: 'white',
                fontSize: '0.9rem',
                outline: 'none',
              }}
            />
          </div>

          <div style={{ marginBottom: '1.25rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.4rem' }}>
              <label style={{ fontSize: '0.85rem', fontWeight: 600, color: '#cbd5e1' }}>
                Target Web Application URL
              </label>
              <button
                type="button"
                onClick={() => setTargetUrl('http://localhost:3000/signup')}
                style={{ background: 'transparent', border: 'none', color: '#38bdf8', fontSize: '0.75rem', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '3px' }}
              >
                <Sparkles size={12} /> Built-in Demo Flow
              </button>
            </div>
            <input
              type="url"
              value={targetUrl}
              onChange={(e) => setTargetUrl(e.target.value)}
              required
              placeholder="https://staging.myapp.com/signup"
              style={{
                width: '100%',
                padding: '10px 12px',
                borderRadius: '8px',
                background: '#090d16',
                border: '1px solid #1e293b',
                color: 'white',
                fontSize: '0.9rem',
                outline: 'none',
              }}
            />
          </div>

          <div style={{ marginBottom: '1.5rem' }}>
            <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#cbd5e1', marginBottom: '0.4rem' }}>
              User Task & Goal Description
            </label>
            <textarea
              value={goal}
              onChange={(e) => setGoal(e.target.value)}
              rows={3}
              required
              style={{
                width: '100%',
                padding: '10px 12px',
                borderRadius: '8px',
                background: '#090d16',
                border: '1px solid #1e293b',
                color: 'white',
                fontSize: '0.9rem',
                outline: 'none',
                resize: 'none',
                fontFamily: 'inherit',
              }}
            />
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
            <button
              type="button"
              onClick={onClose}
              style={{
                background: '#090d16',
                color: '#94a3b8',
                border: '1px solid #1e293b',
                borderRadius: '8px',
                padding: '8px 16px',
                fontSize: '0.85rem',
                fontWeight: 600,
                cursor: 'pointer',
              }}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.4rem',
                background: '#2563eb',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                padding: '8px 20px',
                fontSize: '0.85rem',
                fontWeight: 600,
                cursor: loading ? 'not-allowed' : 'pointer',
                opacity: loading ? 0.7 : 1,
              }}
            >
              <Play size={16} /> {loading ? 'Initializing Swarm...' : 'Execute Study'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
