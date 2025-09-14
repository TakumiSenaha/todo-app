"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import Header from "@/components/Header";
import TodoForm from "@/components/TodoForm";
import TodoCard from "@/components/TodoCard";
import TodoSort from "@/components/TodoSort";
import {
  api,
  Todo,
  CreateTodoRequest,
  UpdateTodoRequest,
  ApiError,
} from "@/services/api";

export default function DashboardPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loadingTodos, setLoadingTodos] = useState(true);
  const [currentSort, setCurrentSort] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/login");
    }
  }, [user, isLoading, router]);

  const loadTodos = useCallback(async () => {
    setLoadingTodos(true);
    setError(null);
    try {
      const todosData = await api.todo.getTodos(currentSort);
      setTodos(todosData);
    } catch (error) {
      console.error("Failed to load todos:", error);
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError("Failed to load todos");
      }
    } finally {
      setLoadingTodos(false);
    }
  }, [currentSort]);

  useEffect(() => {
    if (user) {
      loadTodos();
    }
  }, [user, loadTodos]);

  const handleCreateTodo = async (todoData: CreateTodoRequest) => {
    setIsSubmitting(true);
    setError(null);
    try {
      const newTodo = await api.todo.createTodo(todoData);
      setTodos((prev) => [newTodo, ...prev]);
    } catch (error) {
      console.error("Failed to create todo:", error);
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError("Failed to create todo");
      }
      throw error;
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleComplete = async (id: number) => {
    try {
      const updatedTodo = await api.todo.toggleTodoComplete(id);
      setTodos((prev) =>
        prev.map((todo) => (todo.id === id ? updatedTodo : todo)),
      );
    } catch (error) {
      console.error("Failed to toggle todo:", error);
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError("Failed to update todo");
      }
    }
  };

  const handleUpdateTodo = async (
    id: number,
    updateData: UpdateTodoRequest,
  ) => {
    try {
      const updatedTodo = await api.todo.updateTodo(id, updateData);
      setTodos((prev) =>
        prev.map((todo) => (todo.id === id ? updatedTodo : todo)),
      );
    } catch (error) {
      console.error("Failed to update todo:", error);
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError("Failed to update todo");
      }
    }
  };

  const handleDeleteTodo = async (id: number) => {
    try {
      await api.todo.deleteTodo(id);
      setTodos((prev) => prev.filter((todo) => todo.id !== id));
    } catch (error) {
      console.error("Failed to delete todo:", error);
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError("Failed to delete todo");
      }
    }
  };

  const handleSortChange = (sortBy: string) => {
    setCurrentSort(sortBy);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null; // Will redirect to login
  }

  const completedCount = todos.filter((todo) => todo.is_completed).length;
  const totalCount = todos.length;

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="max-w-4xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Todo Dashboard</h1>
          <p className="mt-2 text-gray-600">
            Welcome back, {user.username}! Manage your tasks efficiently.
          </p>
          <div className="mt-4 bg-white rounded-lg shadow p-4">
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-600">
                <span className="font-medium">{totalCount}</span> total tasks,
                <span className="font-medium text-green-600 ml-1">
                  {completedCount}
                </span>{" "}
                completed
              </div>
              <div className="w-64 bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-500 h-2 rounded-full transition-all duration-300"
                  style={{
                    width:
                      totalCount > 0
                        ? `${(completedCount / totalCount) * 100}%`
                        : "0%",
                  }}
                ></div>
              </div>
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg flex items-center justify-between">
            <span>{error}</span>
            <button
              onClick={() => setError(null)}
              className="text-red-500 hover:text-red-700"
            >
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          </div>
        )}

        {/* Todo Creation Form */}
        <TodoForm onSubmit={handleCreateTodo} isSubmitting={isSubmitting} />

        {/* Sort Controls */}
        <TodoSort currentSort={currentSort} onSortChange={handleSortChange} />

        {/* Todo List */}
        <div className="space-y-4">
          {loadingTodos ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
              <span className="ml-2 text-gray-600">Loading todos...</span>
            </div>
          ) : todos.length === 0 ? (
            <div className="text-center py-12 bg-white rounded-lg shadow">
              <svg
                className="mx-auto h-12 w-12 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900">
                No todos yet
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Get started by creating your first todo above.
              </p>
            </div>
          ) : (
            todos.map((todo) => (
              <TodoCard
                key={todo.id}
                todo={todo}
                onToggleComplete={handleToggleComplete}
                onUpdate={handleUpdateTodo}
                onDelete={handleDeleteTodo}
              />
            ))
          )}
        </div>
      </main>
    </div>
  );
}
