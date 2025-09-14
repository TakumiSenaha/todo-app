"use client";

import { useState } from "react";
import { Todo, UpdateTodoRequest } from "@/services/api";

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
  const [editData, setEditData] = useState({
    title: todo.title,
    due_date: todo.due_date || "",
    priority: todo.priority,
  });
  const [isLoading, setIsLoading] = useState(false);

  const getPriorityLabel = (priority: number): string => {
    switch (priority) {
      case 0:
        return "Low";
      case 1:
        return "Medium";
      case 2:
        return "High";
      default:
        return "Low";
    }
  };

  const getPriorityColor = (priority: number): string => {
    switch (priority) {
      case 0:
        return "bg-gray-100 text-gray-800";
      case 1:
        return "bg-yellow-100 text-yellow-800";
      case 2:
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const isOverdue = (dueDateString?: string): boolean => {
    if (!dueDateString) return false;
    const dueDate = new Date(dueDateString);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return dueDate < today && !todo.is_completed;
  };

  const handleToggleComplete = async () => {
    setIsLoading(true);
    try {
      await onToggleComplete(todo.id);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async () => {
    setIsLoading(true);
    try {
      const updateData: UpdateTodoRequest = {
        title: editData.title.trim(),
        priority: editData.priority,
      };

      if (editData.due_date) {
        updateData.due_date = editData.due_date;
      }

      await onUpdate(todo.id, updateData);
      setIsEditing(false);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    setEditData({
      title: todo.title,
      due_date: todo.due_date || "",
      priority: todo.priority,
    });
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (window.confirm("Are you sure you want to delete this todo?")) {
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
        todo.is_completed
          ? "border-green-400 opacity-75"
          : isOverdue(todo.due_date)
            ? "border-red-400"
            : "border-blue-400"
      }`}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-start space-x-3 flex-1">
          {/* Completion Checkbox */}
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

          {/* Todo Content */}
          <div className="flex-1">
            {isEditing ? (
              <div className="space-y-3">
                <input
                  type="text"
                  value={editData.title}
                  onChange={(e) =>
                    setEditData({ ...editData, title: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Todo title..."
                />

                <input
                  type="date"
                  value={editData.due_date}
                  onChange={(e) =>
                    setEditData({ ...editData, due_date: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                <select
                  value={editData.priority}
                  onChange={(e) =>
                    setEditData({
                      ...editData,
                      priority: parseInt(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value={0}>Low</option>
                  <option value={1}>Medium</option>
                  <option value={2}>High</option>
                </select>
              </div>
            ) : (
              <div>
                <h3
                  className={`text-lg font-medium ${
                    todo.is_completed
                      ? "text-gray-500 line-through"
                      : "text-gray-900"
                  }`}
                >
                  {todo.title}
                </h3>

                <div className="mt-2 flex flex-wrap items-center gap-3 text-sm">
                  {/* Priority Badge */}
                  <span
                    className={`px-2 py-1 rounded-full text-xs font-medium ${getPriorityColor(todo.priority)}`}
                  >
                    {getPriorityLabel(todo.priority)}
                  </span>

                  {/* Due Date */}
                  {todo.due_date && (
                    <span
                      className={`flex items-center ${
                        isOverdue(todo.due_date)
                          ? "text-red-600 font-medium"
                          : "text-gray-600"
                      }`}
                    >
                      <svg
                        className="w-4 h-4 mr-1"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                        />
                      </svg>
                      Due: {formatDate(todo.due_date)}
                      {isOverdue(todo.due_date) && " (Overdue)"}
                    </span>
                  )}

                  {/* Created Date */}
                  <span className="text-gray-500">
                    Created: {formatDate(todo.created_at)}
                  </span>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center space-x-2 ml-3">
          {isEditing ? (
            <>
              <button
                onClick={handleSave}
                disabled={isLoading || !editData.title.trim()}
                className="bg-green-500 hover:bg-green-600 disabled:bg-green-300 text-white p-2 rounded transition-colors"
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
              </button>
              <button
                onClick={handleCancel}
                disabled={isLoading}
                className="bg-gray-500 hover:bg-gray-600 text-white p-2 rounded transition-colors"
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
              </button>
            </>
          ) : (
            <>
              <button
                onClick={() => setIsEditing(true)}
                disabled={isLoading}
                className="text-blue-500 hover:text-blue-700 p-2 rounded transition-colors"
                title="Edit"
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
                title="Delete"
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
