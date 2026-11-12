import React from 'react';
import { Activity, PlayCircle, BarChart3, Users, Flame, Plus, GitCompare } from 'lucide-react';

interface HeaderProps {
  activeTab: 'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare';
  setActiveTab: (tab: 'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare') => void;
  onOpenNewRun: () => void;
}

export const Header: React.FC<HeaderProps> = ({ activeTab, setActiveTab, onOpenNewRun }) => {
  return (
    <header style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '1rem 2rem',
      backgroundColor: '#0f172a',
      borderBottom: '1px solid #1e293b',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.85rem' }}>
        <div style={{
          background: 'linear-gradient(135deg, #2563eb, #8b5cf6)',
          width: '38px',
          height: '38px',
          borderRadius: '10px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'white',
        }}>
          <Activity size={22} />
        </div>
        <div>
          <div style={{ fontWeight: 800, fontSize: '1.25rem', letterSpacing: '-0.02em', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            uxswarm <span style={{ fontSize: '0.7rem', padding: '2px 6px', background: '#1e293b', color: '#93c5fd', borderRadius: '4px', fontWeight: 600 }}>STUDIO</span>
          </div>
          <div style={{ fontSize: '0.75rem', color: '#94a3b8' }}>Synthetic Persona Swarm Usability Profiler</div>
        </div>
      </div>

      <nav style={{ display: 'flex', gap: '0.5rem', background: '#141e33', padding: '4px', borderRadius: '10px', border: '1px solid #1e293b' }}>
        <button
          onClick={() => setActiveTab('overview')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'overview' ? '#2563eb' : 'transparent',
            color: activeTab === 'overview' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <BarChart3 size={16} /> Overview
        </button>

        <button
          onClick={() => setActiveTab('replay')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'replay' ? '#2563eb' : 'transparent',
            color: activeTab === 'replay' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <PlayCircle size={16} /> Session Replay
        </button>

        <button
          onClick={() => setActiveTab('heatmaps')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'heatmaps' ? '#2563eb' : 'transparent',
            color: activeTab === 'heatmaps' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <Flame size={16} /> Heatmap & Gaze
        </button>

        <button
          onClick={() => setActiveTab('compare')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'compare' ? '#2563eb' : 'transparent',
            color: activeTab === 'compare' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <GitCompare size={16} /> A/B Compare
        </button>

        <button
          onClick={() => setActiveTab('heuristics')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'heuristics' ? '#2563eb' : 'transparent',
            color: activeTab === 'heuristics' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <Activity size={16} /> Heuristics & Patches
        </button>

        <button
          onClick={() => setActiveTab('personas')}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            padding: '6px 14px',
            borderRadius: '8px',
            fontSize: '0.85rem',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            backgroundColor: activeTab === 'personas' ? '#2563eb' : 'transparent',
            color: activeTab === 'personas' ? 'white' : '#94a3b8',
            transition: 'all 0.15s ease',
          }}
        >
          <Users size={16} /> Swarm Catalog
        </button>
      </nav>

      <div>
        <button
          onClick={onOpenNewRun}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.4rem',
            background: 'linear-gradient(135deg, #2563eb, #1d4ed8)',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            padding: '8px 16px',
            fontSize: '0.85rem',
            fontWeight: 600,
            cursor: 'pointer',
            boxShadow: '0 2px 4px rgba(37,99,235,0.2)',
          }}
        >
          <Plus size={16} /> Launch Study
        </button>
      </div>
    </header>
  );
};
