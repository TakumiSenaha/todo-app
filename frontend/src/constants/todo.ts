// Todo関連の定数定義

export const TODO_PRIORITY = {
  LOW: 0,
  MEDIUM: 1,
  HIGH: 2,
} as const;

export type TodoPriority = (typeof TODO_PRIORITY)[keyof typeof TODO_PRIORITY];

export const PRIORITY_LABELS: Record<TodoPriority, string> = {
  [TODO_PRIORITY.LOW]: "Low",
  [TODO_PRIORITY.MEDIUM]: "Medium",
  [TODO_PRIORITY.HIGH]: "High",
};

export const PRIORITY_COLORS: Record<TodoPriority, string> = {
  [TODO_PRIORITY.LOW]: "bg-gray-100 text-gray-800",
  [TODO_PRIORITY.MEDIUM]: "bg-yellow-100 text-yellow-800",
  [TODO_PRIORITY.HIGH]: "bg-red-100 text-red-800",
};

export const PRIORITY_OPTIONS = [
  { value: TODO_PRIORITY.LOW, label: PRIORITY_LABELS[TODO_PRIORITY.LOW] },
  { value: TODO_PRIORITY.MEDIUM, label: PRIORITY_LABELS[TODO_PRIORITY.MEDIUM] },
  { value: TODO_PRIORITY.HIGH, label: PRIORITY_LABELS[TODO_PRIORITY.HIGH] },
];
