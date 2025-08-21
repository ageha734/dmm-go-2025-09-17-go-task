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

```
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

1. リポジトリをクローン
```bash
git clone https://github.com/dmm-com/dmm-go-2025-09-17-go-task.git
cd dmm-go-2025-09-17-go-task
```

2. 依存関係をインストール
```bash
go mod download
```

3. 環境変数を設定
```bash
cp .env.example .env
# .envファイルを編集してデータベース設定を行う
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

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
