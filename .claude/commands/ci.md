以下を理解して、導入するためのコードを書いてください。

-----

# Lint, Test, CI/CD導入ガイド ⚙️

このガイドでは、プロジェクトの品質を自動的に保ち、開発プロセスを効率化するための「Lint」「Test」「CI」の導入方法を説明します。

## コーディングの注意点
- 末尾（ファイルの最後の行の後）の改行を忘れないようにして下さい。

-----

## 1 Lint（静的解析）

**Lint**は、コードを実行する前に、記述スタイルが統一されているか、潜在的なバグがないかをチェックするプロセスです。

### バックエンド (Go)

  - **ツール**: **`golangci-lint`** を使用します。これはGoのLintツールとして最も標準的で、多くのチェックを一度に実行できます。

  - **設定ファイル (`.golangci.yml`)**:
    プロジェクトのルートに設定ファイルを作成し、チームのルールを定義します。

    ```yaml
    # .golangci.yml
    run:
      timeout: 5m
    linters:
      enable:
        - gofmt
        - goimports
        - revive
        - govet
        - staticcheck
        - unused
        - ineffassign
    issues:
      exclude-rules:
        - path: _test\.go
          linters:
            - funlen
    ```

  - **実行コマンド**:
    `backend`ディレクトリで以下のコマンドを実行します。

    ```bash
    golangci-lint run ./...
    ```

### フロントエンド (Next.js)

  - **ツール**: **ESLint** (コード品質) と **Prettier** (コードフォーマット) を使用します。これらは`create-next-app`で初期設定されています。

  - **設定**:
    `package.json`にルールを追加したり、`.eslintrc.json`や`.prettierrc`ファイルを編集してカスタマイズします。

  - **実行コマンド**:
    `frontend`ディレクトリで以下のコマンドを実行します。

    ```bash
    # ESLintでチェック
    npm run lint

    # Prettierでフォーマット
    npm run format
    ```

-----

## 2Test（テスト） 🧪

**Test**は、コードが期待通りに動作することを保証するためのプロセスです。

### バックエンド (Go)

  - **ツール**: Go標準の`testing`パッケージを使用します。

  - **種類**:

      - **ユニットテスト**: 関数やメソッド単体の動作をテストします。
      - **インテグレーションテスト**: データベースを含めた全体の連携をテストします。CI環境では、テスト実行時に一時的なDBを立ち上げてテストを行います。

  - **実行コマンド**:
    `backend`ディレクトリで以下のコマンドを実行します。

    ```bash
    go test -v ./...
    ```

### フロントエンド (Next.js)

  - **ツール**: **Jest**と**React Testing Library**の組み合わせが標準的です。コンポーネントがユーザーから見て正しく表示・動作するかをテストします。

  - **実行コマンド**:
    `frontend`ディレクトリで`package.json`にテストスクリプトを追加して実行します。

    ```bash
    npm run test
    ```

-----

## 3CI（継続的インテグレーション） 🚀

**CI**は、GitHubにコードをプッシュするたびに、上で設定したLintとTestを**自動的に実行する仕組み**です。これにより、問題のあるコードがマージされるのを防ぎます。

  - **ツール**: **GitHub Actions** を使用します。

  - **設定ファイル (`.github/workflows/ci.yml`)**:
    プロジェクトのルートに以下の設定ファイルを作成します。このファイルは「バックエンド」と「フロントエンド」のチェックを並行して実行します。

<!-- end list -->

```yaml
# .github/workflows/ci.yml

name: CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  # --- Backend Job ---
  backend-ci:
    runs-on: ubuntu-latest
    services:
      # テスト用のPostgreSQLデータベースを一時的に起動
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpassword
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Install golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: v1.55

    - name: Run Lint
      working-directory: ./backend
      run: golangci-lint run ./...

    - name: Run Tests
      working-directory: ./backend
      # テスト用のDB接続情報を環境変数として渡す
      env:
        DB_SOURCE_TEST: "postgresql://testuser:testpassword@localhost:5432/testdb?sslmode=disable"
      run: go test -v ./...

  # --- Frontend Job ---
  frontend-ci:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
        cache-dependency-path: frontend/package-lock.json

    - name: Install Dependencies
      working-directory: ./frontend
      run: npm install

    - name: Run Lint
      working-directory: ./frontend
      run: npm run lint

    - name: Run Tests
      working-directory: ./frontend
      run: npm run test

    - name: Run Build
      working-directory: ./frontend
      run: npm run build
```
