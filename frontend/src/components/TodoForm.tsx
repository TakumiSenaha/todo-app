"use client";

import { useState, memo, useCallback } from "react";
import { CreateTodoRequest } from "@/services/api";
import { getPriorityLabel, getPriorityColor } from "@/utils/todo";
import { PRIORITY_OPTIONS } from "@/constants/todo";

interface TodoFormProps {
  onSubmit: (todoData: CreateTodoRequest) => Promise<void>;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

function TodoForm({ onSubmit, onCancel, isSubmitting = false }: TodoFormProps) {
  const [formData, setFormData] = useState<CreateTodoRequest>({
    title: "",
    due_date: "",
    priority: 0,
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  const validateForm = useCallback((): boolean => {
    const newErrors: { [key: string]: string } = {};

    if (!formData.title.trim()) {
      newErrors.title = "Title is required";
    } else if (formData.title.length > 100) {
      newErrors.title = "Title must be 100 characters or less";
    }

    if (formData.due_date) {
      const dueDate = new Date(formData.due_date);
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      if (dueDate < today) {
        newErrors.due_date = "Due date cannot be in the past";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [formData.title, formData.due_date]);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();

      if (!validateForm()) {
        return;
      }

      const submitData: CreateTodoRequest = {
        title: formData.title.trim(),
        priority: formData.priority,
      };

      if (formData.due_date) {
        submitData.due_date = formData.due_date;
      }

      try {
        await onSubmit(submitData);
        setFormData({ title: "", due_date: "", priority: 0 });
        setErrors({});
      } catch {
        setErrors({ submit: "Failed to create todo. Please try again." });
      }
    },
    [formData, onSubmit, validateForm],
  );

  const handleReset = useCallback(() => {
    setFormData({ title: "", due_date: "", priority: 0 });
    setErrors({});
    if (onCancel) {
      onCancel();
    }
  }, [onCancel]);

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Add New Todo</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Title Input */}
        <div>
          <label
            htmlFor="title"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Title *
          </label>
          <input
            type="text"
            id="title"
            value={formData.title}
            onChange={(e) =>
              setFormData({ ...formData, title: e.target.value })
            }
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.title ? "border-red-500" : "border-gray-300"
            }`}
            placeholder="Enter todo title..."
            maxLength={100}
            disabled={isSubmitting}
          />
          {errors.title && (
            <p className="text-red-500 text-sm mt-1">{errors.title}</p>
          )}
        </div>

        {/* Due Date Input */}
        <div>
          <label
            htmlFor="due_date"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Due Date
          </label>
          <input
            type="date"
            id="due_date"
            value={formData.due_date}
            onChange={(e) =>
              setFormData({ ...formData, due_date: e.target.value })
            }
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.due_date ? "border-red-500" : "border-gray-300"
            }`}
            disabled={isSubmitting}
          />
          {errors.due_date && (
            <p className="text-red-500 text-sm mt-1">{errors.due_date}</p>
          )}
        </div>

        {/* Priority Selection */}
        <div>
          <label
            htmlFor="priority"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Priority
          </label>
          <select
            id="priority"
            value={formData.priority}
            onChange={(e) =>
              setFormData({ ...formData, priority: parseInt(e.target.value) })
            }
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={isSubmitting}
          >
            {PRIORITY_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <p className={`text-sm mt-1 ${getPriorityColor(formData.priority)}`}>
            Priority: {getPriorityLabel(formData.priority)}
          </p>
        </div>

        {/* Error Message */}
        {errors.submit && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {errors.submit}
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex space-x-3 pt-2">
          <button
            type="submit"
            disabled={isSubmitting}
            className="bg-blue-500 hover:bg-blue-600 disabled:bg-blue-300 text-white px-6 py-2 rounded-md font-medium transition-colors duration-200"
          >
            {isSubmitting ? "Adding..." : "Add Todo"}
          </button>
          <button
            type="button"
            onClick={handleReset}
            disabled={isSubmitting}
            className="bg-gray-300 hover:bg-gray-400 disabled:bg-gray-200 text-gray-700 px-6 py-2 rounded-md font-medium transition-colors duration-200"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default memo(TodoForm);
