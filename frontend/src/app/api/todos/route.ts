import { NextRequest, NextResponse } from "next/server";
import { authenticatedFetch } from "@/lib/auth";

// GET /api/todos - Get all todos for authenticated user
export async function GET() {
  try {
    const response = await authenticatedFetch("/api/v1/todos");
    const data = await response.json();

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    if (
      error instanceof Error &&
      error.message === "No authentication token found"
    ) {
      return NextResponse.json(
        { error: "Authentication required" },
        { status: 401 },
      );
    }

    console.error("Get todos error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}

// POST /api/todos - Create new todo
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    const response = await authenticatedFetch("/api/v1/todos", {
      method: "POST",
      body,
    });

    const data = await response.json();

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    if (
      error instanceof Error &&
      error.message === "No authentication token found"
    ) {
      return NextResponse.json(
        { error: "Authentication required" },
        { status: 401 },
      );
    }

    console.error("Create todo error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}
