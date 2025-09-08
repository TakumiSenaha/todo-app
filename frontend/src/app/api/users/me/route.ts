import { NextResponse } from "next/server";
import { authenticatedFetch } from "@/lib/auth";

export async function GET() {
  try {
    const response = await authenticatedFetch("/api/v1/me");
    const data = await response.json();

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    if (
      error instanceof Error &&
      error.message === "No authentication token found"
    ) {
      return NextResponse.json(
        { error: "No authentication token" },
        { status: 401 },
      );
    }

    console.error("Me API error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}
