"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { authApi, ApiError, tokenManager } from "@/services/api";

interface User {
  id: number;
  username: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (
    username: string,
    email: string,
    password: string,
  ) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = async () => {
    try {
      if (!tokenManager.hasToken()) {
        setUser(null);
        setIsLoading(false);
        return;
      }

      const userData = await authApi.getCurrentUser();
      setUser(userData);
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

  const login = async (username: string, password: string) => {
    setIsLoading(true);
    try {
      const data = await authApi.login({ username, password });
      setUser(data.user);

      // Store token using tokenManager
      if (data.token) {
        tokenManager.setToken(data.token);
      }
    } catch (error) {
      setUser(null);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (
    username: string,
    email: string,
    password: string,
  ) => {
    setIsLoading(true);
    try {
      const data = await authApi.register({ username, email, password });

      // For register, we need to create a user object from the response
      const user = {
        id: data.id,
        username: data.username,
        email: data.email,
      };
      setUser(user);

      // Registration doesn't return a token, so we need to login afterwards
      await login(username, password);
    } catch (error) {
      setUser(null);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    setIsLoading(true);
    try {
      await authApi.logout();
      setUser(null);

      // Clear token
      tokenManager.clearToken();
    } catch {
      // バックエンドのログアウトが失敗してもローカル状態をクリア
      setUser(null);
      tokenManager.clearToken();
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const value = {
    user,
    isLoading,
    login,
    register,
    logout,
    checkAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
