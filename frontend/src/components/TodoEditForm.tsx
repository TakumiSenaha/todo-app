// Todo編集フォームコンポーネント

"use client";

import { useState, memo, useCallback } from "react";
import { Todo, UpdateTodoRequest } from "@/services/api";
import { PRIORITY_OPTIONS } from "@/constants/todo";

interface TodoEditFormProps {
  todo: Todo;
  onSave: (updateData: UpdateTodoRequest) => Promise<void>;
  onCancel: () => void;
  isLoading: boolean;
}

function TodoEditForm({
  todo,
  onSave,
  onCancel,
  isLoading,
}: TodoEditFormProps) {
  const [editData, setEditData] = useState({
    title: todo.title,
    due_date: todo.due_date || "",
    priority: todo.priority,
  });

  const handleSave = useCallback(async () => {
    const updateData: UpdateTodoRequest = {
      title: editData.title.trim(),
      priority: editData.priority,
    };

    if (editData.due_date) {
      updateData.due_date = editData.due_date;
    }

    await onSave(updateData);
  }, [editData, onSave]);

  const handleCancel = useCallback(() => {
    setEditData({
      title: todo.title,
      due_date: todo.due_date || "",
      priority: todo.priority,
    });
    onCancel();
  }, [todo, onCancel]);

  return (
    <div className="space-y-3">
      {/* タイトル入力 */}
      <input
        type="text"
        value={editData.title}
        onChange={(e) => setEditData({ ...editData, title: e.target.value })}
        className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Todo title..."
        disabled={isLoading}
      />

      {/* 期日入力 */}
      <input
        type="date"
        value={editData.due_date}
        onChange={(e) => setEditData({ ...editData, due_date: e.target.value })}
        className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        disabled={isLoading}
      />

      {/* 優先度選択 */}
      <select
        value={editData.priority}
        onChange={(e) =>
          setEditData({
            ...editData,
            priority: parseInt(e.target.value),
          })
        }
        className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        disabled={isLoading}
      >
        {PRIORITY_OPTIONS.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>

      {/* アクションボタン */}
      <div className="flex space-x-2 pt-2">
        <button
          onClick={handleSave}
          disabled={isLoading || !editData.title.trim()}
          className="bg-green-500 hover:bg-green-600 disabled:bg-green-300 text-white px-4 py-2 rounded transition-colors flex items-center space-x-1"
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
              d="M5 13l4 4L19 7"
            />
          </svg>
          <span>保存</span>
        </button>

        <button
          onClick={handleCancel}
          disabled={isLoading}
          className="bg-gray-500 hover:bg-gray-600 disabled:bg-gray-400 text-white px-4 py-2 rounded transition-colors flex items-center space-x-1"
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
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
          <span>キャンセル</span>
        </button>
      </div>
    </div>
  );
}

export default memo(TodoEditForm);
