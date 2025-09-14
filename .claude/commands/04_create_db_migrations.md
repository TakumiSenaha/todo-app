# Task 4: Create Initial Database Migration Files

## Goal
`golang-migrate/migrate` を使用して、`users` と `todos` テーブルを作成するための初期マイグレーションファイルをSQLで記述します。
追加改定指示まで読んでください。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Migration Tool**: `golang-migrate/migrate`
- **Directory**: `backend/migrations/`
- **Database Design**:
  - `users` テーブル:
    - `id` (SERIAL, PK), `username` (VARCHAR, UNIQUE), `email` (VARCHAR, UNIQUE), `password_hash` (VARCHAR), `created_at`, `updated_at` (TIMESTAMPTZ)
  - `todos` テーブル:
    - `id` (SERIAL, PK), `user_id` (INTEGER, FK to users.id, ON DELETE CASCADE), `title` (VARCHAR), `due_date` (DATE), `priority` (INTEGER), `is_completed` (BOOLEAN), `created_at`, `updated_at` (TIMESTAMPTZ)
  - **Best Practice**: `updated_at`カラムを自動更新するためのトリガー関数も作成します。

## Instructions
1.  `backend/migrations/` ディレクトリに以下の4つのファイルを作成してください。
    - `000001_create_users_table.up.sql`
    - `000001_create_users_table.down.sql`
    - `000002_create_todos_table.up.sql`
    - `000002_create_todos_table.down.sql`

2.  各ファイルの内容を以下のように記述してください。
    - **`000001_create_users_table.up.sql`**:
        - `updated_at`を自動更新するトリガー関数 `update_updated_at_column` を作成します。
        - `users` テーブルを作成する `CREATE TABLE` 文を記述します。
        - 作成したテーブルに `updated_at` のトリガーを設定します。
    - **`000001_create_users_table.down.sql`**:
        - `users` テーブルを削除する `DROP TABLE` 文を記述します。
        - トリガー関数を削除する `DROP FUNCTION` 文を記述します。
    - **`000002_create_todos_table.up.sql`**:
        - `todos` テーブルを作成する `CREATE TABLE` 文を記述します。
        - `user_id` を外部キーとして設定し、`ON DELETE CASCADE` を含めます。
        - `user_id`, `due_date`, `priority` カラムにインデックスを作成します。
        - `updated_at` のトリガーを設定します。
    - **`000002_create_todos_table.down.sql`**:
        - `todos` テーブルを削除する `DROP TABLE` 文を記述します。

#追加改定指示

# Task 4: Create Initial Database Migration Files

## Goal
`golang-migrate/migrate` を使用して、`users` と `todos` テーブルを作成するための初期マイグレーションファイルをSQLで記述します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Migration Tool**: `golang-migrate/migrate`
- **Directory**: `backend/migrations/`
- **Database Design**:
  - `users` テーブル:
    - `id` (SERIAL, PK), `username` (VARCHAR, UNIQUE), `email` (VARCHAR, UNIQUE), `password_hash` (VARCHAR), `created_at`, `updated_at` (TIMESTAMPTZ)
  - `todos` テーブル:
    - `id` (SERIAL, PK), `user_id` (INTEGER, FK to users.id, ON DELETE CASCADE), `title` (VARCHAR), `due_date` (DATE), `priority` (INTEGER), `is_completed` (BOOLEAN), `created_at`, `updated_at` (TIMESTAMPTZ)
  - **Best Practice**: `updated_at`カラムを自動更新するためのトリガー関数も作成します。

## Instructions
1.  `backend/migrations/` ディレクトリに以下の4つのファイルを作成してください。
    - `000001_create_users_table.up.sql`
    - `000001_create_users_table.down.sql`
    - `000002_create_todos_table.up.sql`
    - `000002_create_todos_table.down.sql`

2.  各ファイルの内容を以下のように記述してください。
    - **`000001_create_users_table.up.sql`**:
        - `updated_at`を自動更新するトリガー関数 `update_updated_at_column` を作成します。
        - `users` テーブルを作成する `CREATE TABLE` 文を記述します。
        - 作成したテーブルに `updated_at` のトリガーを設定します。
    - **`000001_create_users_table.down.sql`**:
        - `users` テーブルを削除する `DROP TABLE` 文を記述します。
        - トリガー関数を削除する `DROP FUNCTION` 文を記述します。
    - **`000002_create_todos_table.up.sql`**:
        - `todos` テーブルを作成する `CREATE TABLE` 文を記述します。
        - `user_id` を外部キーとして設定し、`ON DELETE CASCADE` を含めます。
        - `user_id`, `due_date`, `priority` カラムにインデックスを作成します。
        - `updated_at` のトリガーを設定します。
    - **`000002_create_todos_table.down.sql`**:
        - `todos` テーブルを削除する `DROP TABLE` 文を記述します。
