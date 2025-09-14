// リファクタリング済みTodoCardコンポーネント

"use client";

import { useState } from "react";
import { Todo, UpdateTodoRequest } from "@/services/api";
import { getTodoBorderColor } from "@/utils/todo";
import TodoDisplay from "./TodoDisplay";
import TodoEditForm from "./TodoEditForm";

interface TodoCardProps {
  todo: Todo;
  onToggleComplete: (id: number) => Promise<void>;
  onUpdate: (id: number, data: UpdateTodoRequest) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export default function TodoCard({
  todo,
  onToggleComplete,
  onUpdate,
  onDelete,
}: TodoCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleToggleComplete = async () => {
    setIsLoading(true);
    try {
      await onToggleComplete(todo.id);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async (updateData: UpdateTodoRequest) => {
    setIsLoading(true);
    try {
      await onUpdate(todo.id, updateData);
      setIsEditing(false);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm("この Todo を削除してもよろしいですか？")) {
      setIsLoading(true);
      try {
        await onDelete(todo.id);
      } finally {
        setIsLoading(false);
      }
    }
  };

  return (
    <div
      className={`bg-white rounded-lg shadow-md p-4 border-l-4 ${
        todo.is_completed ? "opacity-75" : ""
      } ${getTodoBorderColor(todo)}`}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-start space-x-3 flex-1">
          {/* 完了チェックボックス */}
          <button
            onClick={handleToggleComplete}
            disabled={isLoading}
            className={`mt-1 w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${
              todo.is_completed
                ? "bg-green-500 border-green-500 text-white"
                : "border-gray-300 hover:border-green-400"
            }`}
          >
            {todo.is_completed && (
              <svg className="w-3 h-3" viewBox="0 0 20 20" fill="currentColor">
                <path
                  fillRule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clipRule="evenodd"
                />
              </svg>
            )}
          </button>

          {/* Todo コンテンツ */}
          <div className="flex-1">
            {isEditing ? (
              <TodoEditForm
                todo={todo}
                onSave={handleSave}
                onCancel={() => setIsEditing(false)}
                isLoading={isLoading}
              />
            ) : (
              <TodoDisplay todo={todo} />
            )}
          </div>
        </div>

        {/* アクションボタン */}
        <div className="flex items-center space-x-2 ml-3">
          {!isEditing && (
            <>
              <button
                onClick={() => setIsEditing(true)}
                disabled={isLoading}
                className="text-blue-500 hover:text-blue-700 p-2 rounded transition-colors"
                title="編集"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
              </button>
              <button
                onClick={handleDelete}
                disabled={isLoading}
                className="text-red-500 hover:text-red-700 p-2 rounded transition-colors"
                title="削除"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
