// Todo関連のユーティリティ関数

import { Todo } from "@/services/api";
import {
  TodoPriority,
  PRIORITY_LABELS,
  PRIORITY_COLORS,
} from "@/constants/todo";

/**
 * 優先度のラベルを取得
 */
export const getPriorityLabel = (priority: number): string => {
  return PRIORITY_LABELS[priority as TodoPriority] || PRIORITY_LABELS[0];
};

/**
 * 優先度の色クラスを取得
 */
export const getPriorityColor = (priority: number): string => {
  return PRIORITY_COLORS[priority as TodoPriority] || PRIORITY_COLORS[0];
};

/**
 * 期限切れかどうかを判定
 */
export const isOverdue = (dueDateString?: string): boolean => {
  if (!dueDateString) return false;

  const dueDate = new Date(dueDateString);
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return dueDate < today;
};

/**
 * Todoの境界線色を取得
 */
export const getTodoBorderColor = (todo: Todo): string => {
  if (todo.is_completed) {
    return "border-green-400";
  }
  if (isOverdue(todo.due_date)) {
    return "border-red-400";
  }
  return "border-blue-400";
};

/**
 * Todoの完了率を計算
 */
export const calculateCompletionRate = (
  todos: Todo[],
): {
  total: number;
  completed: number;
  percentage: number;
} => {
  const total = todos.length;
  const completed = todos.filter((todo) => todo.is_completed).length;
  const percentage = total > 0 ? (completed / total) * 100 : 0;

  return { total, completed, percentage };
};
