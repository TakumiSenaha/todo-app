# Task 7: Create Next.js Frontend Skeleton

## Goal
`create-next-app`で生成されるような、Next.js (App Router) の基本的なファイル構造と骨組みを作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Framework**: Next.js (App Router, TypeScript)
- **Styling**: Tailwind CSS
- **Backend API URL**: `http://localhost:8080` (開発環境)

## Instructions
1.  Next.jsプロジェクトの基本的な設定ファイルを作成してください。
    - `frontend/package.json` (Next.js, React, TypeScript, Tailwind CSSの依存関係を含む)
    - `frontend/tsconfig.json` (Next.js推奨設定)
    - `frontend/tailwind.config.ts` (基本的な設定)
    - `frontend/postcss.config.js` (基本的な設定)

2.  基本的なページ構成を作成してください。
    - `frontend/app/globals.css` (Tailwind CSSのディレクティブを含む)
    - `frontend/app/layout.tsx` (基本的なHTML構造とbodyタグを定義)
    - `frontend/app/page.tsx` (トップページ。"ToDo App"などの見出しを表示するだけのシンプルなもの)

3.  環境変数を設定してください。
    - `frontend/.env.local` を作成し、`NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` と記述してください。

4.  API通信の雛形を作成してください。
    - `frontend/lib/apiClient.ts` を作成し、`axios`または`fetch`を使って、環境変数で設定したベースURLを持つAPIクライアントのインスタンスを作成してください。

# 追加変更
# Task 7 (Revised): Create Next.js Frontend Skeleton (using pnpm)

## Goal
`pnpm`を使い、Next.js (App Router) の基本的なファイル構造と骨組みを作成します。

## Context
- **Framework**: Next.js (App Router, TypeScript)
- **Package Manager**: `pnpm`
- **Styling**: Tailwind CSS
- **Backend API URL**: `http://localhost:8080` (開発環境)

## Instructions
1.  Next.jsプロジェクトの基本的な設定ファイルを作成してください。
    -   `frontend/package.json`: `next`, `react`, `typescript`, `tailwindcss` の依存関係を記述します。`scripts` には `dev`, `build`, `start`, `lint` を含めます。
    -   `frontend/tsconfig.json`: Next.js推奨設定を記述します。
    -   `frontend/tailwind.config.ts`, `frontend/postcss.config.js`: Tailwind CSSの基本的な設定を記述します。

2.  基本的なページ構成を作成してください。
    -   `frontend/app/globals.css`: Tailwind CSSのディレectiveを記述します。
    -   `frontend/app/layout.tsx`: `AuthProvider`（後で作成）で全体をラップする基本的なHTML構造を定義します。
    -   `frontend/app/page.tsx`: トップページ。「ToDo App」などの見出しを表示するだけのシンプルなものにします。

3.  環境変数を設定してください。
    -   `frontend/.env.local` を作成し、`NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` と記述してください。

# 追加指示
# Task 7 (Revised): Create Next.js Frontend Skeleton (using pnpm)

## Goal
`pnpm`を使い、Next.js (App Router) の基本的なファイル構造と骨組みを作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Framework**: Next.js (App Router, TypeScript)
- **Package Manager**: `pnpm`
- **Styling**: Tailwind CSS
- **Backend API URL**: `http://localhost:8080` (開発環境)

## Instructions
1.  Next.jsプロジェクトの基本的な設定ファイルを作成してください。
    -   `frontend/package.json`: `next`, `react`, `typescript`, `tailwindcss` の依存関係を記述します。`scripts` には `dev`, `build`, `start`, `lint` を含めます。
    -   `frontend/tsconfig.json`: Next.js推奨設定を記述します。
    -   `frontend/tailwind.config.ts`, `frontend/postcss.config.js`: Tailwind CSSの基本的な設定を記述します。

2.  基本的なページ構成を作成してください。
    -   `frontend/app/globals.css`: Tailwind CSSのディレクティブを記述します。
    -   `frontend/app/layout.tsx`: `AuthProvider`（後で作成）で全体をラップする基本的なHTML構造を定義します。
    -   `frontend/app/page.tsx`: トップページ。「ToDo App」などの見出しを表示するだけのシンプルなものにします。

3.  環境変数を設定してください。
    -   `frontend/.env.local` を作成し、`NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` と記述してください。
