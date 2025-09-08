// API service layer for all backend communications

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "";

// Types
export interface User {
  id: number;
  username: string;
  email: string;
}

export interface UpdateProfileRequest {
  username: string;
  email: string;
  current_password: string;
  new_password: string;
}

export interface UpdateProfileResponse {
  user: User;
  message: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  message: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  id: number;
  username: string;
  email: string;
  message: string;
}

// API Error class with enhanced error handling
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public isNetworkError: boolean = false,
    public isAuthError: boolean = false,
  ) {
    super(message);
    this.name = "ApiError";

    // 認証エラーの判定
    this.isAuthError = status === 401 || status === 403;
  }

  // エラータイプの判定メソッド
  static isNetworkError(error: unknown): boolean {
    return error instanceof ApiError && error.isNetworkError;
  }

  static isAuthError(error: unknown): boolean {
    return error instanceof ApiError && error.isAuthError;
  }

  static isTemporaryError(error: unknown): boolean {
    if (error instanceof ApiError) {
      return (
        error.isNetworkError || error.status >= 500 || error.status === 408
      ); // Request Timeout
    }
    return false;
  }
}

// Cookie management utilities
export const cookieManager = {
  getCookie(name: string): string | null {
    if (typeof document === "undefined") return null;
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
  },

  clearCookie(name: string): void {
    if (typeof document === "undefined") return;
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
  },

  hasAuthToken(): boolean {
    return !!this.getCookie("auth_token");
  },
};

// Enhanced API request function with retry logic
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  includeAuth: boolean = false,
  retryCount: number = 0,
  maxRetries: number = 2,
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const defaultHeaders: Record<string, string> = {
    "Content-Type": "application/json",
  };

  const config: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
    credentials: "include", // Cookie自動送信
  };

  try {
    const response = await fetch(url, config);

    // レスポンスのパース（JSON以外のレスポンスも考慮）
    let data;
    try {
      data = await response.json();
    } catch {
      data = { message: "Invalid server response" };
    }

    if (!response.ok) {
      const isNetworkError = response.status >= 500 || response.status === 0;
      const error = new ApiError(
        data.message || `HTTP ${response.status}: ${response.statusText}`,
        response.status,
        data.code,
        isNetworkError,
      );

      // 401エラーの場合、Cookieが無効なので削除
      if (response.status === 401 && includeAuth) {
        cookieManager.clearCookie("auth_token");
      }

      throw error;
    }

    return data;
  } catch (error) {
    if (error instanceof ApiError) {
      // リトライ可能なエラーの場合、リトライを実行
      if (retryCount < maxRetries && ApiError.isTemporaryError(error)) {
        await new Promise((resolve) =>
          setTimeout(resolve, 1000 * (retryCount + 1)),
        ); // 指数バックオフ
        return apiRequest<T>(
          endpoint,
          options,
          includeAuth,
          retryCount + 1,
          maxRetries,
        );
      }
      throw error;
    }

    // ネットワークエラーまたはその他のエラー
    const networkError = new ApiError(
      "Network error or server unavailable",
      0,
      "NETWORK_ERROR",
      true, // ネットワークエラーとしてマーク
    );

    // ネットワークエラーもリトライ可能
    if (retryCount < maxRetries) {
      await new Promise((resolve) =>
        setTimeout(resolve, 1000 * (retryCount + 1)),
      );
      return apiRequest<T>(
        endpoint,
        options,
        includeAuth,
        retryCount + 1,
        maxRetries,
      );
    }

    throw networkError;
  }
}

// Auth API functions
export const authApi = {
  // Login
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    return apiRequest<LoginResponse>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
  },

  // Register
  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    return apiRequest<RegisterResponse>("/api/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  },

  // Logout
  async logout(): Promise<{ status: string; message: string }> {
    return apiRequest<{ status: string; message: string }>("/api/auth/logout", {
      method: "POST",
    });
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    return apiRequest<User>(
      "/api/users/me",
      {
        method: "GET",
      },
      true,
    );
  },

  // Update profile
  async updateProfile(
    profileData: UpdateProfileRequest,
  ): Promise<UpdateProfileResponse> {
    return apiRequest<UpdateProfileResponse>(
      "/api/auth/profile",
      {
        method: "PUT",
        body: JSON.stringify(profileData),
      },
      true,
    );
  },
};

// Export auth API
export const api = {
  auth: authApi,
};
