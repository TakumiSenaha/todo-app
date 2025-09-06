import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function GET(_request: NextRequest) {
  try {
    // Get auth token from cookies
    const cookieStore = await cookies();
    const authToken = cookieStore.get("auth_token");

    if (!authToken) {
      return NextResponse.json(
        { error: "No authentication token" },
        { status: 401 },
      );
    }

    // Forward request to Go backend with cookie
    const response = await fetch(`${BACKEND_URL}/api/v1/me`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Cookie: `auth_token=${authToken.value}`,
      },
    });

    const data = await response.json();

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error("Me API error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}
