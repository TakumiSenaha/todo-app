import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

/**
 * Get authentication token from cookies or localStorage (fallback)
 */
export async function getAuthToken(): Promise<string | null> {
  const cookieStore = await cookies();
  const authToken = cookieStore.get("auth_token");

  // Cookieがある場合はそれを使用
  if (authToken?.value) {
    return authToken.value;
  }

  // サーバーサイドでは localStorage にアクセスできないため、
  // ここではCookieのみをチェック
  return null;
}

/**
 * Create authenticated headers for backend requests
 */
export async function createAuthHeaders(): Promise<HeadersInit> {
  const token = await getAuthToken();
  if (!token) {
    throw new Error("No authentication token found");
  }

  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`, // Bearer token形式に変更
  };
}

/**
 * Make authenticated request to backend
 */
export async function authenticatedFetch(
  endpoint: string,
  options: {
    method?: string;
    body?: unknown;
    headers?: HeadersInit;
  } = {},
): Promise<Response> {
  const token = await getAuthToken();

  const requestOptions: RequestInit = {
    method: options.method || "GET",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include", // Include cookies in request
  };

  // If we have a token, add it as both cookie and header for maximum compatibility
  if (token) {
    requestOptions.headers = {
      ...requestOptions.headers,
      Authorization: `Bearer ${token}`,
      Cookie: `auth_token=${token}`,
    };
  } else {
    throw new Error("No authentication token found");
  }

  if (options.body) {
    requestOptions.body = JSON.stringify(options.body);
  }

  return fetch(`${BACKEND_URL}${endpoint}`, requestOptions);
}

/**
 * Wrapper for API routes that require authentication
 * Returns 401 if no auth token is found
 */
export async function withAuth<T extends unknown[]>(
  handler: (...args: T) => Promise<NextResponse>,
) {
  return async (...args: T): Promise<NextResponse> => {
    try {
      const token = await getAuthToken();
      if (!token) {
        return NextResponse.json(
          { error: "No authentication token" },
          { status: 401 },
        );
      }
      return await handler(...args);
    } catch (error) {
      console.error("Authentication error:", error);
      return NextResponse.json(
        { error: "Authentication failed" },
        { status: 401 },
      );
    }
  };
}

/**
 * Higher-order function to create authenticated API handlers
 */
export function createAuthenticatedHandler(
  handler: (request: NextRequest, token: string) => Promise<NextResponse>,
) {
  return async (request: NextRequest): Promise<NextResponse> => {
    try {
      // Cookieからトークンを取得
      const token = await getAuthToken();

      if (!token) {
        return NextResponse.json(
          { error: "No authentication token" },
          { status: 401 },
        );
      }

      return await handler(request, token);
    } catch (error) {
      console.error("Authentication error:", error);
      return NextResponse.json(
        { error: "Authentication failed" },
        { status: 401 },
      );
    }
  };
}
