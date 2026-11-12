import React, { useState } from 'react';
import { Activity, PlayCircle, BarChart3, Users, Flame, Plus, GitCompare, Menu, X } from 'lucide-react';

interface HeaderProps {
  activeTab: 'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare';
  setActiveTab: (tab: 'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare') => void;
  onOpenNewRun: () => void;
}

export const Header: React.FC<HeaderProps> = ({ activeTab, setActiveTab, onOpenNewRun }) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const tabs: {
    id: 'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare';
    label: string;
    icon: React.ReactNode;
  }[] = [
    { id: 'overview', label: 'Overview', icon: <BarChart3 size={16} /> },
    { id: 'replay', label: 'Session Replay', icon: <PlayCircle size={16} /> },
    { id: 'heatmaps', label: 'Heatmap & Gaze', icon: <Flame size={16} /> },
    { id: 'compare', label: 'A/B Compare', icon: <GitCompare size={16} /> },
    { id: 'heuristics', label: 'Heuristics & Patches', icon: <Activity size={16} /> },
    { id: 'personas', label: 'Swarm Catalog', icon: <Users size={16} /> },
  ];

  const handleSelectTab = (id: typeof activeTab) => {
    setActiveTab(id);
    setMobileMenuOpen(false);
  };

  return (
    <header className="header-wrapper">
      <div className="header-container">
        {/* Brand */}
        <div className="header-brand">
          <div
            style={{
              background: 'linear-gradient(135deg, #2563eb, #8b5cf6)',
              width: '36px',
              height: '36px',
              borderRadius: '9px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: 'white',
              flexShrink: 0,
            }}
          >
            <Activity size={20} />
          </div>
          <div>
            <div className="header-brand-title">
              uxswarm <span className="header-brand-badge">STUDIO</span>
            </div>
            <div className="header-brand-subtitle">
              Synthetic Persona Swarm Usability Profiler
            </div>
          </div>
        </div>

        {/* Desktop Navigation */}
        <nav className="header-nav">
          {tabs.map((tab) => {
            const isActive = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                onClick={() => handleSelectTab(tab.id)}
                className={`nav-tab-btn ${isActive ? 'active' : 'inactive'}`}
              >
                {tab.icon}
                <span>{tab.label}</span>
              </button>
            );
          })}
        </nav>

        {/* Actions & Mobile Toggle */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexShrink: 0 }}>
          <button onClick={onOpenNewRun} className="btn-launch-study">
            <Plus size={16} />
            <span className="btn-label">Launch Study</span>
          </button>

          <button
            className="mobile-menu-toggle"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            aria-label="Toggle Navigation Menu"
          >
            {mobileMenuOpen ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>
      </div>

      {/* Mobile Drawer Dropdown */}
      <div className={`mobile-drawer ${mobileMenuOpen ? 'open' : ''}`}>
        {tabs.map((tab) => {
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => handleSelectTab(tab.id)}
              className={`nav-tab-btn ${isActive ? 'active' : 'inactive'}`}
            >
              {tab.icon}
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Mobile Fast-Scroll Tab Rail (Always Accessible on Touch Devices) */}
      <div className="mobile-nav-rail">
        {tabs.map((tab) => {
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => handleSelectTab(tab.id)}
              className={`nav-tab-btn ${isActive ? 'active' : 'inactive'}`}
              style={{ padding: '6px 12px', fontSize: '0.8rem' }}
            >
              {tab.icon}
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>
    </header>
  );
};
