// Todo管理用カスタムフック

import { useState, useEffect, useCallback } from "react";
import {
  api,
  Todo,
  CreateTodoRequest,
  UpdateTodoRequest,
} from "@/services/api";
import { handleApiError } from "@/utils/error";

export interface UseTodosReturn {
  todos: Todo[];
  isLoading: boolean;
  error: string | null;
  isSubmitting: boolean;
  currentSort: string;

  // Actions
  loadTodos: () => Promise<void>;
  createTodo: (todoData: CreateTodoRequest) => Promise<void>;
  updateTodo: (id: number, updateData: UpdateTodoRequest) => Promise<void>;
  deleteTodo: (id: number) => Promise<void>;
  toggleComplete: (id: number) => Promise<void>;
  setSort: (sortBy: string) => void;
  clearError: () => void;
}

export function useTodos(): UseTodosReturn {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [currentSort, setCurrentSort] = useState("");

  // Todo一覧の読み込み
  const loadTodos = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const todosData = await api.todo.getTodos(currentSort);
      setTodos(todosData);
    } catch (err) {
      handleApiError(err, setError, "Todoの読み込みに失敗しました");
    } finally {
      setIsLoading(false);
    }
  }, [currentSort]);

  // ソート変更時の再読み込み
  useEffect(() => {
    loadTodos();
  }, [loadTodos]);

  // Todo作成
  const createTodo = useCallback(
    async (todoData: CreateTodoRequest): Promise<void> => {
      setIsSubmitting(true);
      setError(null);

      try {
        const newTodo = await api.todo.createTodo(todoData);
        setTodos((prev) => [newTodo, ...prev]);
      } catch (err) {
        handleApiError(err, setError, "Todoの作成に失敗しました");
        throw err; // 呼び出し元でエラーハンドリングが必要な場合のため
      } finally {
        setIsSubmitting(false);
      }
    },
    [],
  );

  // Todo更新
  const updateTodo = useCallback(
    async (id: number, updateData: UpdateTodoRequest): Promise<void> => {
      try {
        const updatedTodo = await api.todo.updateTodo(id, updateData);
        setTodos((prev) =>
          prev.map((todo) => (todo.id === id ? updatedTodo : todo)),
        );
      } catch (err) {
        handleApiError(err, setError, "Todoの更新に失敗しました");
        throw err;
      }
    },
    [],
  );

  // Todo削除
  const deleteTodo = useCallback(async (id: number): Promise<void> => {
    try {
      await api.todo.deleteTodo(id);
      setTodos((prev) => prev.filter((todo) => todo.id !== id));
    } catch (err) {
      handleApiError(err, setError, "Todoの削除に失敗しました");
      throw err;
    }
  }, []);

  // 完了状態切り替え
  const toggleComplete = useCallback(async (id: number): Promise<void> => {
    try {
      const updatedTodo = await api.todo.toggleTodoComplete(id);
      setTodos((prev) =>
        prev.map((todo) => (todo.id === id ? updatedTodo : todo)),
      );
    } catch (err) {
      handleApiError(err, setError, "Todoの状態変更に失敗しました");
      throw err;
    }
  }, []);

  // ソート設定
  const setSort = useCallback((sortBy: string) => {
    setCurrentSort(sortBy);
  }, []);

  // エラークリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    todos,
    isLoading,
    error,
    isSubmitting,
    currentSort,
    loadTodos,
    createTodo,
    updateTodo,
    deleteTodo,
    toggleComplete,
    setSort,
    clearError,
  };
}
