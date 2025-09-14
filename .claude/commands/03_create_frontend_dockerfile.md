# Task 3: Create Frontend Dockerfile (using pnpm)

## Goal
Next.jsフロントエンドアプリケーションをビルドし、本番環境で実行するための `frontend/Dockerfile` を作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Framework**: Next.js
- **Package Manager**: `pnpm` (高速で効率的なモダンパッケージマネージャー)
- **Best Practice**: ビルド成果物のみを実行環境に含める**マルチステージビルド**を採用します。

## Instructions
1. `frontend/Dockerfile` という名前のファイルを作成してください。
2. 以下の2つのステージを持つマルチステージビルドを実装してください。

    - **`builder` ステージ (ビルド用)**:
      - `node:18-alpine` イメージをベースにします。
      - `pnpm` をグローバルにインストールします (`RUN npm install -g pnpm`)。
      - ワーキングディレクトリを `/app` に設定します。
      - `package.json` と `pnpm-lock.yaml` をコピーします。
      - `pnpm install` で依存関係をインストールします。
      - 全てのソースコードをコピーします。
      - `pnpm run build` でアプリケーションをビルドします。

    - **`final` ステージ (実行用)**:
      - `node:18-alpine` イメージをベースにします。
      - `pnpm` をグローバルにインストールします (`RUN npm install -g pnpm`)。
      - ワーキングディレクトリを `/app` に設定します。
      - `builder` ステージから `public` ディレクトリと `.next` ディレクトリをコピーします。
      - `package.json` と `pnpm-lock.yaml` をコピーします。
      - `pnpm install --prod` で本番用の依存関係のみをインストールします。
      - ポート `3000` を公開します。
      - コンテナ起動時に `pnpm start` が実行されるように `CMD` を設定します。

## Output Format
- 作成した `frontend/Dockerfile` の全内容を、単一の `Dockerfile` コードブロックで出力してください。
