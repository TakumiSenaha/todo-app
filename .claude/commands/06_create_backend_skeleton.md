# Task 6 (Revised): Create Go Backend Skeleton with Clean Architecture

## Goal
クリーンアーキテクチャの思想に基づき、Goバックエンドアプリケーションの**ディレクトリ構造と各層の基本的なファイル**を一括で作成します。

## 前提知識
プロジェクトルートにある、`CLAUDE.md`を参照して学習してください。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

## Context
- **Architecture**: Clean Architecture
- **Web Server**: `net/http` (標準ライブラリ)
- **Main Entrypoint**: `cmd/api/main.go`

## Instructions
1.  **ディレクトリ構造の作成**: `backend/` 内に、クリーンアーキテクチャの各層に対応する以下のディレクトリとGoファイルを作成してください。
    ```
    backend/
    ├── cmd/api/main.go
    └── internal/
        ├── domain/
        │   ├── user.go
        │   └── todo.go
        ├── usecase/
        │   ├── user_interactor.go
        │   └── user_repository.go # ← DB操作のインターフェース
        ├── interface/
        │   ├── controller/
        │   │   └── user_controller.go
        │   └── repository/
        │       # SQLクエリファイルがここに配置される
        └── infrastructure/
            └── persistence/
                └── user_persistence.go # ← インターフェースの実装
    ```

2.  **`domain/user.go`**: `User`エンティティの構造体を定義してください（DBのタグなどは含めない純粋なもの）。

3.  **`usecase/user_repository.go`**: `UserRepository` という**インターフェース**を定義してください。これには、ユーザーを登録したり検索したりするメソッドのシグネチャ（`CreateUser(...) error`など）を記述します。

4.  **`infrastructure/persistence/user_persistence.go`**:
    - `UserRepository`インターフェースを**実装**する構造体を作成します。
    - ここで初めて、SQLCが生成したコードを呼び出して、具体的なDB操作を記述します。

5.  **`interface/controller/user_controller.go`**:
    - `net/http` のリクエストを処理するハンドラ関数（`RegisterUserHandler`など）を記述します。
    - この関数はユースケース層の処理を呼び出します。

6.  **`cmd/api/main.go`**:
    - アプリケーションのエントリーポイントです。
    - DBに接続します。
    - **依存性の注入 (Dependency Injection)** を行います。
        - `persistence` (リポジトリ実装) のインスタンスを作成
        - `usecase` (ビジネスロジック) のインスタンスにリポジトリを注入
        - `controller` (ハンドラ) のインスタンスにユースケースを注入
    - `net/http`サーバーを起動し、リクエストをハンドラにルーティングします。
    - `/health` ヘルスチェックエンドポイントを実装します。

