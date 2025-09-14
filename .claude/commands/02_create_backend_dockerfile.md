# Task 2: Create Backend Dockerfile

## Goal
Goバックエンドアプリケーションをビルドし、実行するための `backend/Dockerfile` を作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Language**: Go
- **Best Practice**: パフォーマンスとセキュリティのため、最終的なイメージサイズを最小限に抑える**マルチステージビルド**を採用します。

## Instructions
1.  `backend/Dockerfile` という名前のファイルを作成してください。
2.  以下の2つのステージを持つマルチステージビルドを実装してください。

    - **`builder` ステージ (ビルド用)**:
      - `golang:1.21-alpine` イメージをベースにします。
      - ワーキングディレクトリを `/app` に設定します。
      - `go.mod` と `go.sum` をコピーし、`go mod download` で依存関係をダウンロードします。
      - 全てのソースコードをコピーします。
      - `CGO_ENABLED=0 GOOS=linux go build -o /main .` コマンドでアプリケーションをビルドします。

    - **`final` ステージ (実行用)**:
      - `alpine:latest` イメージをベースにします。
      - `/app` ディレクトリを作成します。
      - `builder` ステージからビルドされたバイナリ `/main` を `/app/main` にコピーします。
      - ポート `8080` を公開します。
      - コンテナ起動時に `./app/main` が実行されるように `CMD` を設定します。

## Output Format
- 作成した `backend/Dockerfile` の全内容を、単一の `Dockerfile` コードブロックで出力してください。
