export const LABEL_KEYS = {
  NAME:    'monty.io/name',
  VERSION: 'monty.io/agent-version'
};

export type Labels = { [key: string]: string };

export interface Reference {
  id: string;
}

export interface Timestamp {
  seconds: number;
  nanos: number;
}

export interface Duration {
  seconds: number;
  nanos: number;
}

export type State = 'success' | 'warning' | 'error';

export interface Status {
  state: State;
  message: string;
  longMessage?: string;
}
