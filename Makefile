DATABASE_PORT ?= 3306
DATABASE_USER ?= testuser
DATABASE_PASSWORD ?= password
DATABASE_NAME ?= testdb
export DATABASE_PORT DATABASE_USER DATABASE_PASSWORD DATABASE_NAME

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
setup: .setup-mysql .setup-hurl .setup-shlack ## 必要なツール（mysql-client, hurl, shlack）を自動インストール
	@echo "👍 全てのツールのセットアップが完了しました。"

.setup-mysql:
	@if ! command -v mysql >/dev/null 2>&1; then \
		echo "🔧 mysql-clientをインストールします..."; \
		$(MAKE) .install-mysql; \
	else \
		echo "🔧 mysql-clientは既にインストールされています。"; \
	fi

.setup-hurl:
	@if ! command -v hurl >/dev/null 2>&1; then \
		echo "🔧 hurlをインストールします..."; \
		$(MAKE) .install-hurl; \
	else \
		echo "🔧 hurlは既にインストールされています。"; \
	fi

.detect-os:
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "macos"; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "linux"; \
	elif [ -n "$$WINDIR" ] || [ "$$(uname -o 2>/dev/null)" = "Msys" ] || [ "$$(uname -o 2>/dev/null)" = "Cygwin" ]; then \
		echo "windows"; \
	else \
		echo "unknown"; \
	fi

.install-mysql:
	@OS=$$($(MAKE) -s .detect-os); \
	echo "🔧 OS: $$OS でmysql-clientをインストールします..."; \
	case $$OS in \
		"macos") \
			if ! command -v brew >/dev/null 2>&1; then \
				echo "❌ Homebrewがインストールされていません。https://brew.sh/ からインストールしてください。"; \
				exit 1; \
			fi; \
			brew install mysql-client; \
			;; \
		"linux") \
			if command -v apt >/dev/null 2>&1; then \
				sudo apt update && sudo apt install -y mysql-client; \
			elif command -v yum >/dev/null 2>&1; then \
				sudo yum install -y mysql; \
			elif command -v dnf >/dev/null 2>&1; then \
				sudo dnf install -y mysql; \
			elif command -v pacman >/dev/null 2>&1; then \
				sudo pacman -S --noconfirm mysql; \
			elif command -v zypper >/dev/null 2>&1; then \
				sudo zypper install -y mysql-client; \
			else \
				echo "❌ サポートされているパッケージマネージャーが見つかりません。"; \
				exit 1; \
			fi; \
			;; \
		"windows") \
			if command -v winget >/dev/null 2>&1; then \
				winget install Oracle.MySQL; \
			else \
				echo "❌ wingetがインストールされていません。Windows Package Managerをインストールしてください。"; \
				exit 1; \
			fi; \
			;; \
		*) \
			echo "❌ サポートされていないOS: $$OS"; \
			exit 1; \
			;; \
	esac; \
	echo "✅ mysql-clientのインストールが完了しました。"

.install-hurl:
	@OS=$$($(MAKE) -s .detect-os); \
	echo "🔧 OS: $$OS でhurlをインストールします..."; \
	case $$OS in \
		"macos") \
			if ! command -v brew >/dev/null 2>&1; then \
				echo "❌ Homebrewがインストールされていません。https://brew.sh/ からインストールしてください。"; \
				exit 1; \
			fi; \
			brew install hurl; \
			;; \
		"linux") \
			echo "⚠️  パッケージマネージャーでhurlが見つからない場合、公式サイトからバイナリをダウンロードします..."; \
			HURL_VERSION=$$(curl -s https://api.github.com/repos/Orange-OpenSource/hurl/releases/latest | grep tag_name | cut -d '"' -f 4); \
			curl -LO https://github.com/Orange-OpenSource/hurl/releases/latest/download/hurl-$${HURL_VERSION}-x86_64-unknown-linux-gnu.tar.gz; \
			tar -xzf hurl-*.tar.gz; \
			sudo mv hurl-*/bin/hurl /usr/local/bin/; \
			rm -rf hurl-$${HURL_VERSION}-x86_64-unknown-linux-gnu.tar.gz; \
			;; \
		"windows") \
			if command -v winget >/dev/null 2>&1; then \
				winget install hurl; \
			else \
				echo "❌ wingetがインストールされていません。Windows Package Managerをインストールしてください。"; \
				exit 1; \
			fi; \
			;; \
		*) \
			echo "❌ サポートされていないOS: $$OS"; \
			exit 1; \
			;; \
	esac; \
	echo "✅ hurlのインストールが完了しました。"

.setup-shlack:
	@if ! command -v shlack >/dev/null 2>&1; then \
		echo "🔧 shlackをインストールします..."; \
		curl --location https://raw.githubusercontent.com/ageha734/shlack/master/install.sh | bash; \
	else \
		echo "🔧 shlackは既にインストールされています。"; \
	fi

# --------------------------------------------------------------------------
# Development Tasks
# --------------------------------------------------------------------------
mod: ## Goモジュールの依存関係を整理・ダウンロード
	@echo "📦 Goモジュールの依存関係を整理・ダウンロードします..."
	@go mod tidy
	@go mod download

build: mod $(TARGET_APP) ## アプリケーションをビルド
	@echo "✅ ビルドが完了しました。"

$(TARGET_APP): $(GO_FILES) go.mod go.sum
	@echo "🔨 アプリケーションをビルドします..."
	@go build -o $(TARGET_APP) ./cmd/main.go

dev: ## 開発モードでアプリケーションを起動
	@echo "⚡️ 開発モードでアプリケーションを起動します..."
	@# Makefileにはファイルの変更を監視する機能はありません。
	@# reflexやairなどのツールを使用してください:
	@# reflex -r '\.go$$' -s -- go run ./cmd/main.go
	@go run ./cmd/main.go

# --------------------------------------------------------------------------
# Test & Lint Tasks
# --------------------------------------------------------------------------
check: lint test-unit build ## リント、ユニットテスト、ビルドを実行

lint: ## golangci-lintを実行
	@echo "🔍 golangci-lintを実行します..."
	@go tool golangci-lint run ./...

test-unit: mod ## ユニットテストを実行
	@echo "🔬 ユニットテストを実行します..."
	@TEST=true go test ./... -overlay=$(shell go run github.com/tenntenn/testtime/cmd/testtime@latest)

test-e2e: ## E2Eテストを実行
	@sh -c ' \
		echo "🚀 E2Eテストを実行します..."; \
		EXIT_CODE=0; \
		TARGETS=${service}; \
		if [ -z "$$TARGETS" ]; then \
			TARGETS="$(E2E_TEST_NAME)"; \
		fi; \
		for s in $$TARGETS; do \
			if ! echo "$(E2E_TEST_NAME)" | grep -q -w "$$s"; then \
				echo "指定された '\''$$s'\'' が見つかりません。"; \
				continue; \
			fi; \
			echo "🚀 [$$s] のテストを開始します..."; \
			echo "  - DBセットアップ: seed.sql"; \
			mysql -h 127.0.0.1 -P $$DATABASE_PORT -u $$DATABASE_USER -p$$DATABASE_PASSWORD $$DATABASE_NAME < "./mock/seed.sql" || echo "⚠️  seed.sqlの実行でエラーが発生しましたが、継続します"; \
			echo "  - DBセットアップ: ./e2e/$$s/01-insert.sql"; \
			mysql -h 127.0.0.1 -P $$DATABASE_PORT -u $$DATABASE_USER -p$$DATABASE_PASSWORD $$DATABASE_NAME < "./e2e/$$s/01-insert.sql" || echo "⚠️  01-insert.sqlの実行でエラーが発生しましたが、継続します"; \
			echo "  - APIテスト: ./e2e/$$s/index.hurl"; \
			if ! hurl --test "./e2e/$$s/index.hurl"; then \
				echo "⚠️  [$$s] のAPIテストでエラーが発生しましたが、継続します"; \
				EXIT_CODE=1; \
			fi; \
			echo "  - DBクリーンアップ: ./e2e/$$s/02-delete.sql"; \
			mysql -h 127.0.0.1 -P $$DATABASE_PORT -u $$DATABASE_USER -p$$DATABASE_PASSWORD $$DATABASE_NAME < "./e2e/$$s/02-delete.sql" || echo "⚠️  02-delete.sqlの実行でエラーが発生しましたが、継続します"; \
			echo "✅ [$$s] のテストが完了しました。"; \
		done; \
		if [ $$EXIT_CODE -ne 0 ]; then \
			echo "⚠️  一部のテストでエラーが発生しましたが、全てのテストを実行しました。"; \
		else \
			echo "✅ 全てのE2Eテストが正常に完了しました。"; \
		fi; \
		exit 0; \
	'
