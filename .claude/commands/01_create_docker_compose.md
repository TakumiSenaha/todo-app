# Task 1 (Revised): Create docker-compose.yml

## Goal
プロジェクトの全体的なコンテナ構成を定義する `docker-compose.yml` ファイルを作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Architecture**:
  - **Backend**: Go (`net/http` standard library) with **Clean Architecture**.
  - **Frontend**: Next.js (with **pnpm** as package manager).
  - **Database**: PostgreSQL.
- **Directory Structure**:
  ```
  todo-app/
  ├── backend/
  │   ├── cmd/api/main.go
  │   └── Dockerfile
  ├── frontend/
  │   └── Dockerfile
  └── docker-compose.yml
  ```

## Instructions
1.  プロジェクトのルートに `docker-compose.yml` という名前のファイルを作成してください。
2.  ファイルには `db`, `backend`, `frontend` の3つのサービスを定義します。
3.  各サービスは以下の仕様を満たすように記述してください。

    - **`db` サービス**:
      - `image`: `postgres:15-alpine` を使用
      - `volumes`: データベースのデータを永続化するため、`postgres-data:/var/lib/postgresql/data` を設定
      - `environment`:
        - `POSTGRES_USER`: `user`
        - `POSTGRES_PASSWORD`: `password`
        - `POSTGRES_DB`: `todo_db`
      - `ports`: `5432:5432`

    - **`backend` サービス**:
      - `build`: `backend` ディレクトリの `Dockerfile` を使用 (`context: ./backend`)
      - `ports`: `8080:8080`
      - `environment`:
        - `DB_SOURCE`: `postgresql://user:password@db:5432/todo_db?sslmode=disable`
      - `depends_on`: `db` サービスが起動してから開始するように設定

    - **`frontend` サービス**:
      - `build`: `frontend` ディレクトリの `Dockerfile` を使用 (`context: ./frontend`)
      - `ports`: `3000:3000`
      - `depends_on`: `backend` サービス

4.  トップレベルに `postgres-data` という名前付きvolumeを定義してください。
