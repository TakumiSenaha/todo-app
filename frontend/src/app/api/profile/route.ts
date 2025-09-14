import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

interface UpdateProfileRequest {
  username?: string;
  email?: string;
  current_password?: string;
  new_password?: string;
}

export async function PUT(request: NextRequest) {
  try {
    const body: UpdateProfileRequest = await request.json();

    // Validate request body
    if (!body.username || !body.email) {
      return NextResponse.json(
        {
          message: "ユーザー名とメールアドレスは必須です",
          errors: {
            username: !body.username ? "ユーザー名は必須です" : undefined,
            email: !body.email ? "メールアドレスは必須です" : undefined,
          },
        },
        { status: 400 },
      );
    }

    // Validate username format
    if (body.username.length < 3 || body.username.length > 20) {
      return NextResponse.json(
        {
          message: "ユーザー名は3-20文字で入力してください",
          errors: { username: "ユーザー名は3-20文字で入力してください" },
        },
        { status: 400 },
      );
    }

    if (!/^[a-zA-Z0-9_]+$/.test(body.username)) {
      return NextResponse.json(
        {
          message: "ユーザー名は英数字とアンダースコアのみ使用できます",
          errors: {
            username: "ユーザー名は英数字とアンダースコアのみ使用できます",
          },
        },
        { status: 400 },
      );
    }

    // Validate email format
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(body.email)) {
      return NextResponse.json(
        {
          message: "有効なメールアドレスを入力してください",
          errors: { email: "有効なメールアドレスを入力してください" },
        },
        { status: 400 },
      );
    }

    // Validate password if changing
    if (body.new_password) {
      if (!body.current_password) {
        return NextResponse.json(
          {
            message: "現在のパスワードを入力してください",
            errors: { currentPassword: "現在のパスワードを入力してください" },
          },
          { status: 400 },
        );
      }

      if (body.new_password.length < 8) {
        return NextResponse.json(
          {
            message: "パスワードは8文字以上で入力してください",
            errors: { newPassword: "パスワードは8文字以上で入力してください" },
          },
          { status: 400 },
        );
      }

      if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(body.new_password)) {
        return NextResponse.json(
          {
            message: "パスワードは英数字の両方を含む必要があります",
            errors: {
              newPassword: "パスワードは英数字の両方を含む必要があります",
            },
          },
          { status: 400 },
        );
      }
    }

    // Forward request to backend
    const backendResponse = await fetch(`${BACKEND_URL}/api/v1/profile`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Cookie: request.headers.get("cookie") || "",
      },
      body: JSON.stringify(body),
      credentials: "include",
    });

    const backendData = await backendResponse.json();

    if (!backendResponse.ok) {
      // Handle backend errors
      if (backendResponse.status === 409) {
        // Conflict error (username or email already exists)
        return NextResponse.json(
          {
            message:
              backendData.message ||
              "ユーザー名またはメールアドレスが既に使用されています",
            errors: backendData.errors || {},
          },
          { status: 409 },
        );
      }

      if (backendResponse.status === 401) {
        // Invalid current password
        return NextResponse.json(
          {
            message: "現在のパスワードが正しくありません",
            errors: { currentPassword: "現在のパスワードが正しくありません" },
          },
          { status: 400 },
        );
      }

      return NextResponse.json(
        {
          message: backendData.message || "プロフィールの更新に失敗しました",
          errors: backendData.errors || {},
        },
        { status: backendResponse.status },
      );
    }

    return NextResponse.json({
      status: "success",
      message: "プロフィールが正常に更新されました",
      user: backendData.user,
    });
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

    console.error("Profile update error:", error);
    return NextResponse.json(
      {
        message: "サーバーエラーが発生しました",
        errors: { general: "サーバーエラーが発生しました" },
      },
      { status: 500 },
    );
  }
}
