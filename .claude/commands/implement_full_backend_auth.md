# Task 8 (Enhanced Simple Auth): Implement Full Backend & Frontend Authentication

## Goal
リフレッシュトークンを使わない、シンプルかつ堅牢なJWT認証システムをバックエンドとフロントエンドに**一括で実装**します。バリデーション、エラーハンドリング、CORS設定など、実践的な要素をすべて含みます。

## Context
- **Security Model**: A single, stateful JWT (e.g., 24-hour validity) stored in an `HttpOnly`, `Secure`, `SameSite=Strict` cookie. No refresh tokens.
- **Backend**: Go with Clean Architecture, `net/http`, and `go-playground/validator` for validation.
- **Frontend**: Next.js with BFF (Backend for Frontend) pattern, React Context for state management, and Middleware for route protection.

## Instructions

---
### **Part 1: Backend Implementation (Go)**
---

**Step 1: Usecase Layer & Repository Interfaces**
-   `internal/usecase/user_repository.go` に、`UserRepository` インターフェースを定義します。これには `Create`, `FindByUsername` メソッドを含めます。
-   `internal/usecase/user_interactor.go` に、`Register` と `Login` のビジネスロジックを実装します。`Login`は成功時にJWT文字列を返却します。

**Step 2: Controller & Repository Implementation**
-   `internal/interface/controller/user_controller.go`:
    -   `Register`ハンドラ: リクエストボディをバリデーションし (`go-playground/validator`を使用)、パスワードをハッシュ化してUsecaseを呼び出します。
    -   `Login`ハンドラ: Usecaseを呼び出し、成功したら返されたJWTを`HttpOnly` Cookieに設定して返却します。
    -   `Logout`ハンドラ: Cookieをクリアする処理を実装します。
    -   `Me`ハンドラ: 認証ミドルウェアから渡されたユーザーIDを使い、ユーザー情報を返却します。
-   `internal/infrastructure/persistence/user_persistence.go`:
    -   `UserRepository`インターフェースを実装します。SQLCが生成したコードを呼び出します。

**Step 3: Middleware & Main Entrypoint**
-   `internal/interface/middleware/auth_middleware.go`:
    -   `HttpOnly` CookieからJWTを読み取り、検証する`AuthMiddleware`を実装します。検証成功後、`userID`をリクエストのコンテキストに保存します。
-   `cmd/api/main.go`:
    -   **CORSミドルウェア**を導入し、フロントエンドのURL (`http://localhost:3000`) からのアクセスを許可します。
    -   **ロギングミドルウェア**を導入し、リクエストの情報をログ出力します。
    -   依存性の注入（DI）を行い、各層を接続します。
    -   以下のルーティングを設定します。
        -   `POST /api/v1/register`
        -   `POST /api/v1/login`
        -   `POST /api/v1/logout`
        -   `GET /api/v1/me` ( `AuthMiddleware`で保護)

---
### **Part 2: Frontend Implementation (Next.js)**
---

**Step 1: BFF API Routes (`frontend/app/api/`)**
-   クライアントとGoバックエンドを安全に仲介するAPI Route群を作成します。
    -   `auth/login/route.ts`, `auth/register/route.ts`, `auth/logout/route.ts`: Goの対応するAPIを呼び出し、CookieとJSONレスポンスを中継します。エラーハンドリングも適切に行います。
    -   `users/me/route.ts`: ブラウザからのCookieを使ってGoの`/me` APIを呼び出し、ユーザー情報をクライアントに返します。

**Step 2: Global State Management (`frontend/contexts/AuthContext.tsx`)**
-   React Context APIを使い、グローバルな認証状態を管理する`AuthProvider`を実装します。
-   Contextは `{ user, isLoading, login, logout, register }` のような値を提供します。
-   **アプリケーションのロード時**に、BFFの`/api/users/me` APIを呼び出してログイン状態をチェックし、`user`ステートを初期化するロジックを`useEffect`で実装します。

**Step 3: UI Pages & Layout**
-   `frontend/app/login/page.tsx` と `frontend/app/register/page.tsx` に、フォームと`AuthContext`の関数を呼び出すロジックを実装します。
-   `frontend/app/layout.tsx`で、アプリケーション全体を`AuthProvider`でラップします。

**Step 4: Route Protection (`frontend/middleware.ts`)**
-   プロジェクトルートに`middleware.ts`を作成します。
-   保護したいルート（例: `/dashboard`）へのアクセス時に認証Cookieをチェックし、存在しない場合は`/login`にリダイレクトさせます。
