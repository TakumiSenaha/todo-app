# Task 5: Setup SQLC Configuration and Queries

## Goal
Goバックエンドでデータベースを操作するための`SQLC`の設定ファイルと、基本的なCRUD操作を行うためのSQLクエリファイルを作成します。
追加改定指示まで読んでください。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Tool**: `SQLC`
- **DB Design**: `users`と`todos`テーブルが存在します。
- **Directory**: `backend/`

## Instructions
1.  `backend/`ディレクトリに`sqlc.yaml`という設定ファイルを作成します。内容は以下のようにしてください。
    - `version`: "2"
    - `sql`:
        - `schema`: "db/migration/"
        - `queries`: "db/query/"
        - `engine`: "postgresql"
        - `gen`:
            - `go`:
                - `package`: "db"
                - `out`: "db/sqlc"
                - `emit_json_tags`: true
                - `emit_interface`: true

2.  `backend/db/query/`というディレクトリを作成してください。

3.  `backend/db/query/user.sql`というファイルを作成し、以下のSQLクエリを記述してください。
    - `CreateUser`: 新規ユーザーを作成する
    - `GetUserByUsername`: ユーザー名でユーザー情報を取得する

4.  `backend/db/query/todo.sql`というファイルを作成し、以下のSQLクエリを記述してください。
    - `CreateTodo`: 新規ToDoを作成する
    - `GetTodo`: IDでToDoを一件取得する
    - `ListTodos`: 特定ユーザーの全ToDoを取得する
    - `UpdateTodo`: ToDoを更新する
    - `DeleteTodo`: ToDoを削除する

# 追加改定指示
# Task 5 (Revised): Setup SQLC with Clean Architecture

## Goal
クリーンアーキテクチャのディレクトリ構造に合わせて `SQLC` の設定を行い、リポジトリ層で使われるSQLクエリファイルを作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Tool**: `SQLC`
- **Architecture**: クリーンアーキテクチャでは、データベースとの具体的なやり取りは **Interface Adapters** 層の **Repository** が担当します。
- **New File Structure**:
  - `sqlc.yaml` は `backend/` のルートに配置します。
  - SQLクエリファイルは `backend/internal/interface/repository/` に配置します。
  - 生成されるGoコードは `backend/internal/infrastructure/persistence/` に配置します。

## Instructions
1.  `backend/` ディレクトリに `sqlc.yaml` という設定ファイルを作成し、以下の新しいパス構成で記述してください。
    - `version`: "2"
    - `sql`:
        - `schema`: "migrations/"
        - `queries`: "internal/interface/repository/"
        - `engine`: "postgresql"
        - `gen`:
            - `go`:
                - `package`: "persistence"
                - `out`: "internal/infrastructure/persistence"
                - `emit_json_tags`: true
                - `emit_interface`: true

2.  `backend/internal/interface/repository/` ディレクトリを作成してください。

3.  `backend/internal/interface/repository/user.sql` を作成し、ユーザー管理用の基本的なCRUDクエリを記述してください。
    - `CreateUser`: 新規ユーザーを作成する
    - `GetUserByUsername`: ユーザー名でユーザー情報を取得する

4.  `backend/internal/interface/repository/todo.sql` を作成し、ToDo管理用の基本的なCRUDクエリを記述してください。
    - `CreateTodo`: 新規ToDoを作成する
    - `GetTodo`: IDでToDoを一件取得する
    - `ListTodos`: 特定ユーザーの全ToDoを取得する
    - `UpdateTodo`: ToDoを更新する
    - `DeleteTodo`: ToDoを削除する

## Output Format
- `sqlc.yaml`の全内容を`yaml`コードブロックで出力してください。
- `user.sql`の全内容を`sql`コードブロックで出力してください。
- `todo.sql`の全内容を`sql`コードブロックで出力してください。
