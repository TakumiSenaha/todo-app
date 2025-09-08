import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  authApi,
  LoginRequest,
  RegisterRequest,
  User,
  ApiError,
  tokenManager,
} from "@/services/api";

export interface UseAuthReturn {
  user: User | null;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
  clearError: () => void;
}

export function useAuth(): UseAuthReturn {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  // Check if user is logged in on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        if (!tokenManager.hasToken()) {
          setUser(null);
          setIsLoading(false);
          return;
        }

        const user = await authApi.getCurrentUser();
        setUser(user);
      } catch (error) {
        setUser(null);
        
        // 認証エラーの場合はトークンを削除
        if (ApiError.isAuthError(error)) {
          tokenManager.clearToken();
        } else if (!ApiError.isTemporaryError(error)) {
          // 永続的なエラーの場合も安全のためトークンを削除
          tokenManager.clearToken();
        }
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await authApi.login(credentials);
      setUser(response.user);

      // Store token using tokenManager
      if (response.token) {
        tokenManager.setToken(response.token);
      }

      router.push("/dashboard");
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Login failed");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (userData: RegisterRequest): Promise<void> => {
    setIsLoading(true);
    setError(null);

    try {
      await authApi.register(userData);

      // After successful registration, login automatically
      await login({ username: userData.username, password: userData.password });
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Registration failed");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async (): Promise<void> => {
    setIsLoading(true);
    setError(null);

    try {
      await authApi.logout();
      setUser(null);

      // Clear token
      tokenManager.clearToken();

      router.push("/login");
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Logout failed");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const clearError = () => {
    setError(null);
  };

  return {
    user,
    isLoading,
    login,
    register,
    logout,
    error,
    clearError,
  };
}
