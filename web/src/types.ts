export interface RunRecord {
  id: string;
  name: string;
  target_url: string;
  goal: string;
  status: 'running' | 'completed' | 'failed' | 'queued';
  total_personas: number;
  average_sus: number;
  average_friction: number;
  grade: string;
  mandatory_redesign: boolean;
  created_at: string;
  completed_at?: string;
}

export interface SessionRecord {
  id: string;
  run_id: string;
  persona_id: string;
  persona_name: string;
  archetype: string;
  viewport_width: number;
  viewport_height: number;
  is_mobile: boolean;
  sus_score: number;
  sus_grade: string;
  friction_index: number;
  total_steps: number;
  total_duration: number;
  hesitation_count: number;
  dead_clicks: number;
  backtracks: number;
  goal_completed: boolean;
  dropped_off: boolean;
  drop_off_reason?: string;
  created_at: string;
}

export interface StepRecord {
  id: number;
  session_id: string;
  step_index: number;
  timestamp: string;
  action_type: string;
  target_selector?: string;
  target_description?: string;
  target_role?: string;
  x: number;
  y: number;
  is_dead_click: boolean;
  hesitation_seconds: number;
  url: string;
  inner_speech: string;
  frustration: number;
  screenshot_path?: string;
  frame_data_uri?: string;
}

export interface HeuristicFinding {
  heuristic_id: number;
  heuristic_name: string;
  severity: 'CRITICAL' | 'MAJOR' | 'MINOR';
  title: string;
  evidence: string;
  affected_personas: string[];
  recommendation: string;
}

export interface Persona {
  id: string;
  name: string;
  age: number;
  archetype: string;
  patience_threshold_seconds: number;
  technical_literacy: number;
  attention_span: string;
  viewport: {
    width: number;
    height: number;
    is_mobile: boolean;
  };
  traits: string[];
}
