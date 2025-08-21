ifneq (,$(wildcard ./.env))
    include .env
    export
endif

E2E_TEST_NAME := $(shell find ./e2e -maxdepth 1 -mindepth 1 -type d -exec basename {} \; | sort | uniq)
GO_FILES := $(shell find . -name '*.go')
TARGET_APP := ./.target/app

# --- メインターゲット ---
.DEFAULT_GOAL := help
.PHONY: help setup mod build dev check lint test-unit test-e2e clean

help:
	@echo "利用可能なタスクを一覧表示します"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# --------------------------------------------------------------------------
# Setup Tasks
# --------------------------------------------------------------------------
.PHONY: .setup-mysql .setup-hurl .setup-shlack
setup: .setup-mysql .setup-hurl .setup-shlack
	@echo "👍 全てのツールのセットアップが完了しました。"

.setup-mysql:
	@if ! command -v mysql >/dev/null 2>&1; then \
		echo "🔧 mysql-clientをインストールします..."; \
		echo "TODO: バイナリをダウンロードして、PATHに追加する"; \
	else \
		echo "🔧 mysql-clientは既にインストールされています。"; \
	fi

.setup-hurl:
	@if ! command -v hurl >/dev/null 2>&1; then \
		echo "🔧 hurlをインストールします..."; \
		curl --location https://hurl.dev/install.sh | bash; \
	else \
		echo "🔧 hurlは既にインストールされています。"; \
	fi

.setup-shlack:
	@if ! command -v shlack >/dev/null 2>&1; then \
		echo "🔧 shlackをインストールします..."; \
		curl --location https://raw.githubusercontent.com/dmm-com/shlack/install.sh | bash; \
	else \
		echo "🔧 shlackは既にインストールされています。"; \
	fi

# --------------------------------------------------------------------------
# Development Tasks
# --------------------------------------------------------------------------
mod:
	@echo "📦 Goモジュールの依存関係を整理・ダウンロードします..."
	@go mod tidy
	@go mod download

build: mod $(TARGET_APP)

$(TARGET_APP): $(GO_FILES) go.mod go.sum
	@echo "🔨 アプリケーションをビルドします..."
	@go build -o $(TARGET_APP) ./cmd/main.go

dev:
	@echo "⚡️ 開発モードでアプリケーションを起動します..."
	@# Makefileにはファイルの変更を監視する機能はありません。
	@# reflexやairなどのツールを使用してください:
	@# reflex -r '\.go$$' -s -- go run ./cmd/main.go
	@go run ./cmd/main.go

# --------------------------------------------------------------------------
# Test & Lint Tasks
# --------------------------------------------------------------------------
check: lint test-unit test-e2e

lint:
	@echo "🔍 golangci-lintを実行します..."
	@go tool golangci-lint run ./...

test-unit:
	@echo "🔬 ユニットテストを実行します..."
	@go test ./... -overlay=$(shell go run github.com/tenntenn/testtime/cmd/testtime@latest)

test-e2e:
	@sh -c ' \
		trap "shlack luke \"$$([ $$? -eq 0 ] && echo Success! || echo Failed with exit code $$?)\"" EXIT; \
		echo "🚀 E2Eテストを実行します..."; \
		TARGETS=${service}; \
		if [ -z "$$TARGETS" ]; then \
			TARGETS="$(E2E_TEST_NAME)"; \
		fi; \
		for s in $$TARGETS; do \
			if ! echo "$(E2E_TEST_NAME)" | grep -q -w "$$s"; then \
				echo "指定された '\''$$s'\'' が見つかりません。"; \
				exit 1; \
			fi; \
			echo "🚀 [$$s] のテストを開始します..."; \
			SQL_FILES=$$(find ./e2e/$$s -name "*.sql"); \
			if [ -n "$$SQL_FILES" ]; then \
				for sql_file in $$SQL_FILES; do \
					echo "  - DBセットアップ: $$sql_file"; \
					mysql -u $$DATABASE_USER -h $$DATABASE_HOST -P $$DATABASE_PORT -p$$DATABASE_PASSWORD $$DATABASE_NAME < "$$sql_file"; \
				done; \
			fi; \
			echo "  - APIテスト: ./e2e/$$s/index.hurl"; \
			hurl --test "./e2e/$$s/index.hurl"; \
			echo "✅ [$$s] のテストが完了しました。"; \
		done; \
	'
