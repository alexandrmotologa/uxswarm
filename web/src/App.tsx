import { useState, useEffect } from 'react';
import { Header } from './components/Header.tsx';
import { RunOverview } from './components/RunOverview.tsx';
import { SessionReplay } from './components/SessionReplay.tsx';
import { FrictionHeatmap } from './components/FrictionHeatmap.tsx';
import { HeuristicCard } from './components/HeuristicCard.tsx';
import { PersonaBuilder } from './components/PersonaBuilder.tsx';
import { ABCompare } from './components/ABCompare.tsx';
import { NewRunModal } from './components/NewRunModal.tsx';
import { RunRecord, SessionRecord, HeuristicFinding } from './types.ts';
import { Radio } from 'lucide-react';

interface LiveEvent {
  persona: string;
  event: {
    index: number;
    action_type: string;
    target_description?: string;
    inner_speech?: string;
    hesitation_seconds?: number;
    is_dead_click?: boolean;
  };
}

export function App() {
  const [activeTab, setActiveTab] = useState<'overview' | 'replay' | 'heatmaps' | 'heuristics' | 'personas' | 'compare'>('overview');
  const [runs, setRuns] = useState<RunRecord[]>([]);
  const [selectedRun, setSelectedRun] = useState<RunRecord | null>(null);
  const [sessions, setSessions] = useState<SessionRecord[]>([]);
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [heuristics, setHeuristics] = useState<HeuristicFinding[]>([]);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [liveEvent, setLiveEvent] = useState<LiveEvent | null>(null);
  const [activeStreamingRunId, setActiveStreamingRunId] = useState<string | null>(null);

  // Load runs on initial mount
  const fetchRuns = () => {
    fetch('/api/runs')
      .then((res) => res.json())
      .then((data: RunRecord[]) => {
        setRuns(data || []);
        if (data && data.length > 0) {
          if (!selectedRun) {
            setSelectedRun(data[0]);
          }
        }
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to fetch runs:', err);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchRuns();
  }, []);

  // When selected run changes, load its sessions and heuristics
  useEffect(() => {
    if (!selectedRun) return;

    fetch(`/api/runs/${selectedRun.id}/sessions`)
      .then((res) => res.json())
      .then((data: SessionRecord[]) => {
        setSessions(data || []);
        if (data && data.length > 0) {
          setSelectedSessionId(data[0].id);
        }
      })
      .catch((err) => console.error('Failed to fetch sessions:', err));

    fetch(`/api/runs/${selectedRun.id}/heuristics`)
      .then((res) => res.json())
      .then((data: HeuristicFinding[]) => {
        setHeuristics(data || []);
      })
      .catch((err) => console.error('Failed to fetch heuristics:', err));
  }, [selectedRun?.id]);

  // Connect to SSE stream when a run is active
  useEffect(() => {
    if (!activeStreamingRunId) return;

    const eventSource = new EventSource(`/api/runs/${activeStreamingRunId}/stream`);

    eventSource.addEventListener('step', (e) => {
      try {
        const payload: LiveEvent = JSON.parse(e.data);
        setLiveEvent(payload);
      } catch (err) {
        console.error(err);
      }
    });

    eventSource.addEventListener('complete', () => {
      setActiveStreamingRunId(null);
      setLiveEvent(null);
      fetchRuns();
    });

    eventSource.onerror = () => {
      eventSource.close();
      setActiveStreamingRunId(null);
    };

    return () => {
      eventSource.close();
    };
  }, [activeStreamingRunId]);

  const handleSelectSessionAndReplay = (sessionId: string) => {
    setSelectedSessionId(sessionId);
    setActiveTab('replay');
  };

  const handleRunCreated = (newRun: RunRecord) => {
    setRuns((prev) => [newRun, ...prev]);
    setSelectedRun(newRun);
    setActiveStreamingRunId(newRun.id);
    setActiveTab('overview');
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Header
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        onOpenNewRun={() => setIsModalOpen(true)}
      />

      {/* Live SSE Execution Ticker Bar */}
      {liveEvent && (
        <div style={{
          background: 'linear-gradient(90deg, #1e3a8a, #3b82f6)',
          color: 'white',
          padding: '8px 2rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1rem',
          fontSize: '0.85rem',
          fontWeight: 600,
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <Radio size={16} className="animate-pulse" />
            <span style={{ textTransform: 'uppercase', fontSize: '0.75rem', letterSpacing: '0.05em' }}>
              Live Swarm Stream
            </span>
          </div>
          <div>
            <strong>[{liveEvent.persona}]</strong> #{liveEvent.event.index} ({liveEvent.event.action_type}): {liveEvent.event.target_description || 'Active Step'}
            {liveEvent.event.inner_speech && (
              <span style={{ fontStyle: 'italic', opacity: 0.9, marginLeft: '0.5rem' }}>
                &ldquo;{liveEvent.event.inner_speech}&rdquo;
              </span>
            )}
          </div>
        </div>
      )}

      {/* Run Selector subheader if multiple runs exist */}
      {runs.length > 1 && (
        <div style={{
          background: '#090d16',
          borderBottom: '1px solid #1e293b',
          padding: '8px 2rem',
          display: 'flex',
          alignItems: 'center',
          gap: '1rem',
          fontSize: '0.85rem',
        }}>
          <span style={{ color: '#94a3b8' }}>Selected Study:</span>
          <select
            value={selectedRun?.id || ''}
            onChange={(e) => {
              const r = runs.find((item) => item.id === e.target.value);
              if (r) setSelectedRun(r);
            }}
            style={{
              background: '#141e33',
              border: '1px solid #1e293b',
              color: 'white',
              padding: '4px 10px',
              borderRadius: '6px',
              outline: 'none',
              cursor: 'pointer',
            }}
          >
            {runs.map((r) => (
              <option key={r.id} value={r.id}>
                {r.name} (SUS: {r.average_sus.toFixed(1)}, Grade: {r.grade || 'Running'})
              </option>
            ))}
          </select>
        </div>
      )}

      <main style={{ flex: 1, backgroundColor: '#090d16' }}>
        {loading ? (
          <div style={{ padding: '5rem', textAlign: 'center', color: '#94a3b8' }}>
            Connecting to uxswarm daemon...
          </div>
        ) : (
          <>
            {activeTab === 'overview' && (
              <RunOverview
                run={selectedRun}
                sessions={sessions}
                onSelectSession={handleSelectSessionAndReplay}
              />
            )}

            {activeTab === 'replay' && (
              <SessionReplay
                sessions={sessions}
                selectedSessionId={selectedSessionId}
                onSelectSession={setSelectedSessionId}
              />
            )}

            {activeTab === 'heatmaps' && (
              <FrictionHeatmap sessions={sessions} />
            )}

            {activeTab === 'compare' && (
              <ABCompare runs={runs} />
            )}

            {activeTab === 'heuristics' && (
              <HeuristicCard findings={heuristics} runId={selectedRun?.id} />
            )}

            {activeTab === 'personas' && (
              <PersonaBuilder />
            )}
          </>
        )}
      </main>

      <NewRunModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onRunCreated={handleRunCreated}
      />
    </div>
  );
}
export default App;
