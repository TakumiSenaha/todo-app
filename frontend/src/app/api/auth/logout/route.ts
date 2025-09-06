import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function POST(_request: NextRequest) {
  try {
    // Get auth token from cookies
    const cookieStore = await cookies();
    const authToken = cookieStore.get("auth_token");

    // Forward request to Go backend with cookie
    const response = await fetch(`${BACKEND_URL}/api/v1/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: authToken ? `auth_token=${authToken.value}` : "",
      },
    });

    const data = await response.json();

    // Create response
    const nextResponse = NextResponse.json(data, { status: response.status });

    // Forward cookie clearing from backend
    const setCookieHeader = response.headers.get("set-cookie");
    if (setCookieHeader) {
      nextResponse.headers.set("set-cookie", setCookieHeader);
    }

    return nextResponse;
  } catch (error) {
    console.error("Logout API error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}
