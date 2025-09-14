import { useState } from "react";
import { useRouter } from "next/navigation";
import { authApi, UpdateProfileRequest, ApiError } from "@/services/api";
import { useAuth } from "@/contexts/AuthContext";

export interface UseProfileReturn {
  updateProfile: (data: UpdateProfileRequest) => Promise<void>;
  isUpdating: boolean;
  error: string | null;
  fieldErrors: Record<string, string>;
  success: string | null;
  clearMessages: () => void;
}

export function useProfile(): UseProfileReturn {
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [success, setSuccess] = useState<string | null>(null);
  const router = useRouter();
  const { checkAuth } = useAuth();

  const updateProfile = async (data: UpdateProfileRequest): Promise<void> => {
    setIsUpdating(true);
    setError(null);
    setFieldErrors({});
    setSuccess(null);

    try {
      const response = await authApi.updateProfile(data);
      setSuccess(response.message || "Profile updated successfully!");

      // Refresh user data in auth context
      await checkAuth();

      // Redirect to dashboard after a short delay
      setTimeout(() => {
        router.push("/dashboard");
      }, 2000);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
        if (err.hasFieldErrors()) {
          setFieldErrors(err.fieldErrors || {});
        }
      } else {
        setError("Failed to update profile");
      }
    } finally {
      setIsUpdating(false);
    }
  };

  const clearMessages = () => {
    setError(null);
    setFieldErrors({});
    setSuccess(null);
  };

  return {
    updateProfile,
    isUpdating,
    error,
    fieldErrors,
    success,
    clearMessages,
  };
}
