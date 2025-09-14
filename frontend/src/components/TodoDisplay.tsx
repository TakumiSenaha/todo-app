// Todo表示コンポーネント

"use client";

import { memo } from "react";
import { Todo } from "@/services/api";
import { formatDate } from "@/utils/date";
import { getPriorityLabel, getPriorityColor, isOverdue } from "@/utils/todo";

interface TodoDisplayProps {
  todo: Todo;
}

function TodoDisplay({ todo }: TodoDisplayProps) {
  return (
    <div>
      {/* タイトル */}
      <h3
        className={`text-lg font-medium ${
          todo.is_completed ? "text-gray-500 line-through" : "text-gray-900"
        }`}
      >
        {todo.title}
      </h3>

      {/* メタ情報 */}
      <div className="mt-2 flex flex-wrap items-center gap-3 text-sm">
        {/* 優先度バッジ */}
        <span
          className={`px-2 py-1 rounded-full text-xs font-medium ${getPriorityColor(todo.priority)}`}
        >
          {getPriorityLabel(todo.priority)}
        </span>

        {/* 期日 */}
        {todo.due_date && (
          <span
            className={`flex items-center ${
              isOverdue(todo.due_date) && !todo.is_completed
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
                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 002 2z"
              />
            </svg>
            期日: {formatDate(todo.due_date)}
            {isOverdue(todo.due_date) && !todo.is_completed && " (期限切れ)"}
          </span>
        )}

        {/* 作成日 */}
        <span className="text-gray-500">
          作成: {formatDate(todo.created_at)}
        </span>
      </div>
    </div>
  );
}

export default memo(TodoDisplay);
