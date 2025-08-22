# DMM Go Task - User Management API

Go言語とGin Webフレームワークを使用したユーザー管理APIです。

## 機能

- ユーザーの作成、取得、更新、削除（CRUD操作）
- ユーザー統計情報の取得
- ヘルスチェック機能
- MySQL/SQLiteデータベース対応

## 技術スタック

- **言語**: Go 1.24.0
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: MySQL / SQLite
- **テスト**: testify
- **リンター**: golangci-lint
- **タスクランナー**: Task

## プロジェクト構成

```bash
.
├── cmd/
│   └── main.go              # アプリケーションエントリーポイント
├── handlers/
│   ├── user_handler.go      # ユーザー関連のHTTPハンドラー
│   └── user_handler_test.go # ハンドラーのテスト
├── models/
│   ├── user.go              # ユーザーモデル定義
│   └── user_test.go         # モデルのテスト
├── routes/
│   └── router.go            # ルーティング設定
├── database/
│   └── database.go          # データベース接続設定
├── e2e/
│   ├── home/                # ホーム関連のE2Eテスト
│   └── user/                # ユーザー関連のE2Eテスト
└── scripts/
    ├── compare_speed_check.sh
    └── e2e_speed_check.sh
```

## セットアップ

### 前提条件

- Go 1.24.0以上
- Docker（オプション）

### インストール

#### 1. リポジトリをクローン

```bash
git clone https://github.com/ageha734/dmm-go-2025-09-17-go-task.git
cd dmm-go-2025-09-17-go-task
```

#### 2. 依存関係をインストール

```bash
go mod download
```

#### 3. 環境変数を設定

```bash
cp .env.example .env
```

`.env`ファイルを編集して、以下の環境変数を設定してください：

```bash
# Slack通知設定（オプション）
SLACK_TOKEN=xoxb-your-slack-bot-token

# データベース設定
DATABASE_HOST=127.0.0.1
DATABASE_PORT=3306
DATABASE_USER=root
DATABASE_PASSWORD=password
DATABASE_NAME=testdb
```

**重要**:
- `DATABASE_HOST`は`localhost`ではなく`127.0.0.1`を使用してください（MySQL接続の問題を回避するため）
- 引用符は使用しないでください（例: `DATABASE_HOST=127.0.0.1` ✅、`DATABASE_HOST="127.0.0.1"` ❌）

#### 4. 必要なツールのセットアップ

開発に必要なツール（mysql-client、hurl、shlack）を自動でインストールします：

```bash
# Taskを使用する場合
task setup

# Makeを使用する場合
make setup
```

## 実行方法

### ローカル実行

```bash
go run cmd/main.go
```

### Dockerを使用した実行

```bash
docker-compose up
```

### Taskを使用した実行

```bash
task run
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

### ユーザー管理

- `GET /users` - 全ユーザー取得
- `GET /users/:id` - 特定ユーザー取得
- `POST /users` - ユーザー作成
- `PUT /users/:id` - ユーザー更新
- `DELETE /users/:id` - ユーザー削除

### 統計情報

- `GET /users/stats` - ユーザー統計情報取得

### ヘルスチェック

- `GET /health` - APIの稼働状況確認

## リクエスト例

### ユーザー作成

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "田中太郎",
    "email": "tanaka@example.com",
    "age": 30
  }'
```

### ユーザー取得

```bash
curl http://localhost:8080/users/1
```

## テスト

### 単体テスト実行

```bash
go test ./...
```

### E2Eテスト実行

```bash
task e2e
```

### テストカバレッジ確認

```bash
go test -cover ./...
```

## GitHub Actions との連携

### Environment Secrets の同期

ローカルの`.env`ファイルからGitHub ActionsのEnvironment Secretsに環境変数を自動で同期できます。

#### 前提条件

1. **GitHub CLI のインストール**
   ```bash
   # macOS
   brew install gh

   # Ubuntu/Debian
   sudo apt install gh
   ```

2. **GitHub CLI の認証**
   ```bash
   gh auth login
   ```

3. **リポジトリへのadmin権限**（Environment Secretsの設定に必要）

#### 使用方法

```bash
# ヘルプを表示
./scripts/sync_env_to_github_secrets.sh --help

# dry-runモードで実行（実際には変更せず、実行予定の内容を表示）
./scripts/sync_env_to_github_secrets.sh --dry-run

# 基本的な使用方法（リポジトリは自動検出、デフォルトEnvironment: production）
./scripts/sync_env_to_github_secrets.sh

# 特定のEnvironmentを指定
./scripts/sync_env_to_github_secrets.sh -e staging

# 特定の環境変数ファイルを指定
./scripts/sync_env_to_github_secrets.sh -f .env.production -e production
```

#### 注意事項

- **セキュリティ**: 機密情報を含む環境変数を扱うため、実行前に`--dry-run`で内容を確認してください
- **上書き**: 既存のSecretは警告なしに上書きされます
- **権限**: Environment Secretsの設定にはリポジトリへのadmin権限が必要です

## 開発

### コードフォーマット

```bash
task fmt
```

### リンター実行

```bash
task lint
```

### 利用可能なタスク確認

```bash
task --list
```

### E2Eテスト速度比較

MakeとTaskのE2Eテスト実行時間を比較できます：

```bash
# 全サービスのテストを実行
./scripts/e2e_speed_check.sh

# 特定のサービスのみテスト
./scripts/e2e_speed_check.sh home
./scripts/e2e_speed_check.sh user
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
