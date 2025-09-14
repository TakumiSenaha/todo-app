"use client";

import { memo } from "react";

interface TodoSortProps {
  currentSort: string;
  onSortChange: (sortBy: string) => void;
}

const sortOptions = [
  { value: "", label: "Default (Newest First)" },
  { value: "due_date_asc", label: "Due Date (Ascending)" },
  { value: "due_date_desc", label: "Due Date (Descending)" },
  { value: "priority_desc", label: "Priority (High to Low)" },
  { value: "created_desc", label: "Created Date (Newest)" },
];

function TodoSort({ currentSort, onSortChange }: TodoSortProps) {
  return (
    <div className="flex items-center space-x-4 mb-6">
      <label
        htmlFor="sort-select"
        className="text-sm font-medium text-gray-700"
      >
        Sort by:
      </label>
      <select
        id="sort-select"
        value={currentSort}
        onChange={(e) => onSortChange(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white text-sm"
      >
        {sortOptions.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>

      {/* Sort indicator */}
      <div className="flex items-center text-sm text-gray-500">
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
            d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
          />
        </svg>
        {currentSort ? (
          <span>
            {sortOptions.find((opt) => opt.value === currentSort)?.label}
          </span>
        ) : (
          <span>Default sorting</span>
        )}
      </div>
    </div>
  );
}

export default memo(TodoSort);
