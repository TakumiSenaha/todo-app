import { NextRequest, NextResponse } from "next/server";
import { authenticatedFetch } from "@/lib/auth";

export async function PUT(request: NextRequest) {
  try {
    const body = await request.json();

    const response = await authenticatedFetch("/api/v1/profile", {
      method: "PUT",
      body,
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(
        { error: data.message || "Failed to update profile" },
        { status: response.status },
      );
    }

    return NextResponse.json(data);
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

    console.error("Profile update error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 },
    );
  }
}
