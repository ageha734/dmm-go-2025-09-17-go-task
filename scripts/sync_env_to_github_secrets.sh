#!/bin/bash

set -euo pipefail

# カラー定義
COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

# 設定
DEFAULT_ENV_FILE=".env"
DEFAULT_ENVIRONMENT="production"

# 使用方法を表示
show_help() {
    echo -e "${COLOR_GREEN}GitHub Environment Secrets 同期スクリプト${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  $0 [OPTIONS]"
    echo ""
    echo -e "${COLOR_YELLOW}オプション:${COLOR_NC}"
    echo -e "  -f, --file FILE        環境変数ファイルを指定 (デフォルト: ${DEFAULT_ENV_FILE})"
    echo -e "  -e, --environment ENV  GitHub Environment名を指定 (デフォルト: ${DEFAULT_ENVIRONMENT})"
    echo -e "  -r, --repo REPO        リポジトリ名を指定 (owner/repo形式)"
    echo -e "  -t, --token TOKEN      GitHub Personal Access Tokenを指定"
    echo -e "  -d, --dry-run          実際には実行せず、実行予定の内容を表示"
    echo -e "  -h, --help             このヘルプを表示"
    echo ""
    echo -e "${COLOR_YELLOW}前提条件:${COLOR_NC}"
    echo -e "  - GitHub CLI (gh) がインストールされていること"
    echo -e "  - GitHub Personal Access Token が設定されていること"
    echo -e "  - リポジトリへの admin 権限があること"
    echo ""
    echo -e "${COLOR_YELLOW}例:${COLOR_NC}"
    echo -e "  $0 -f .env -e production -r owner/repo"
    echo -e "  $0 --dry-run"
    echo -e "  $0 --help"
}

# 環境変数ファイルを解析
parse_env_file() {
    local file="$1"
    local -A env_vars=()

    while IFS= read -r line || [ -n "$line" ]; do
        # コメント行をスキップ
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        # 空行をスキップ
        [[ "$line" =~ ^[[:space:]]*$ ]] && continue

        # KEY=VALUE 形式の行を解析
        if [[ "$line" =~ ^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*=[[:space:]]*(.*)$ ]]; then
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"

            # 値から引用符を削除
            value="${value%\"}"
            value="${value#\"}"
            value="${value%\'}"
            value="${value#\'}"

            env_vars["$key"]="$value"
        fi
    done < "$file"

    # 連想配列を出力（キー=値の形式）
    for key in "${!env_vars[@]}"; do
        echo "$key=${env_vars[$key]}"
    done
}

# GitHub CLI の存在確認
check_github_cli() {
    if ! command -v gh >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: GitHub CLI (gh) がインストールされていません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}インストール方法: https://cli.github.com/${COLOR_NC}" >&2
        exit 1
    fi
}

# GitHub認証の確認
check_github_auth() {
    if ! gh auth status >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: GitHub CLI の認証が設定されていません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}認証方法: gh auth login${COLOR_NC}" >&2
        exit 1
    fi
}

# リポジトリの存在確認
check_repository() {
    local repo="$1"

    if ! gh repo view "$repo" >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: リポジトリ '$repo' にアクセスできません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}リポジトリ名が正しいか、アクセス権限があるか確認してください。${COLOR_NC}" >&2
        exit 1
    fi
}

# Environment の存在確認・作成
ensure_environment() {
    local repo="$1"
    local environment="$2"

    echo -e "${COLOR_BLUE}Environment '$environment' の確認中...${COLOR_NC}"

    # GitHub CLI を使用してEnvironmentの存在を確認
    if ! gh api "repos/$repo/environments/$environment" >/dev/null 2>&1; then
        echo -e "${COLOR_YELLOW}Environment '$environment' が存在しません。作成します...${COLOR_NC}"

        # Environment を作成
        gh api --method PUT "repos/$repo/environments/$environment" \
            --field wait_timer=0 \
            --field prevent_self_review=false \
            --field reviewers='[]' \
            --field deployment_branch_policy='null' >/dev/null

        echo -e "${COLOR_GREEN}Environment '$environment' を作成しました。${COLOR_NC}"
    else
        echo -e "${COLOR_GREEN}Environment '$environment' が存在します。${COLOR_NC}"
    fi
}

# シークレットを設定
set_secret() {
    local repo="$1"
    local environment="$2"
    local key="$3"
    local value="$4"
    local dry_run="$5"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} $key を Environment '$environment' に設定予定"
        return
    fi

    echo -e "${COLOR_BLUE}$key を設定中...${COLOR_NC}"

    # GitHub CLI を使用してシークレットを設定
    echo "$value" | gh secret set "$key" --repo "$repo" --env "$environment"

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ $key を設定しました${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ $key の設定に失敗しました${COLOR_NC}" >&2
        return 1
    fi
}

# メイン処理
main() {
    local env_file="$DEFAULT_ENV_FILE"
    local environment="$DEFAULT_ENVIRONMENT"
    local repo=""
    local token=""
    local dry_run="false"

    # 引数解析
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--file)
                env_file="$2"
                shift 2
                ;;
            -e|--environment)
                environment="$2"
                shift 2
                ;;
            -r|--repo)
                repo="$2"
                shift 2
                ;;
            -t|--token)
                token="$2"
                shift 2
                ;;
            -d|--dry-run)
                dry_run="true"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                echo -e "${COLOR_RED}エラー: 不明なオプション '$1'${COLOR_NC}" >&2
                show_help
                exit 1
                ;;
        esac
    done

    # 必須パラメータの確認
    if [ -z "$repo" ]; then
        # 現在のリポジトリから自動取得を試行
        if git remote get-url origin >/dev/null 2>&1; then
            repo_url=$(git remote get-url origin)
            if [[ "$repo_url" =~ github\.com[:/]([^/]+/[^/]+)(\.git)?$ ]]; then
                repo="${BASH_REMATCH[1]}"
                repo="${repo%.git}"
                echo -e "${COLOR_BLUE}リポジトリを自動検出: $repo${COLOR_NC}"
            fi
        fi

        if [ -z "$repo" ]; then
            echo -e "${COLOR_RED}エラー: リポジトリ名が指定されていません。${COLOR_NC}" >&2
            echo -e "${COLOR_YELLOW}-r オプションでリポジトリ名を指定してください (例: owner/repo)${COLOR_NC}" >&2
            exit 1
        fi
    fi

    # 環境変数ファイルの存在確認
    if [ ! -f "$env_file" ]; then
        echo -e "${COLOR_RED}エラー: 環境変数ファイル '$env_file' が見つかりません。${COLOR_NC}" >&2
        exit 1
    fi

    # GitHub CLI と認証の確認
    check_github_cli
    check_github_auth

    # リポジトリの確認
    check_repository "$repo"

    # GitHub Personal Access Token の設定
    if [ -n "$token" ]; then
        export GH_TOKEN="$token"
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  🚀 ${COLOR_GREEN}GitHub Environment Secrets 同期${COLOR_NC}"
    echo -e "  ファイル: ${COLOR_YELLOW}$env_file${COLOR_NC}"
    echo -e "  リポジトリ: ${COLOR_YELLOW}$repo${COLOR_NC}"
    echo -e "  Environment: ${COLOR_YELLOW}$environment${COLOR_NC}"
    if [ "$dry_run" = "true" ]; then
        echo -e "  モード: ${COLOR_YELLOW}DRY RUN${COLOR_NC}"
    fi
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    # Environment の確認・作成
    if [ "$dry_run" = "false" ]; then
        ensure_environment "$repo" "$environment"
    fi

    # 環境変数の解析と設定
    echo -e "${COLOR_BLUE}環境変数を解析中...${COLOR_NC}"

    local count=0
    local success_count=0

    while IFS='=' read -r key value; do
        if [ -n "$key" ] && [ -n "$value" ]; then
            count=$((count + 1))

            if set_secret "$repo" "$environment" "$key" "$value" "$dry_run"; then
                success_count=$((success_count + 1))
            fi
        fi
    done < <(parse_env_file "$env_file")

    echo ""
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  📊 ${COLOR_GREEN}同期結果${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"

    if [ "$dry_run" = "true" ]; then
        echo -e "処理予定の環境変数: ${COLOR_YELLOW}$count${COLOR_NC} 個"
        echo -e "${COLOR_YELLOW}実際の同期を実行するには --dry-run オプションを外してください。${COLOR_NC}"
    else
        echo -e "処理した環境変数: ${COLOR_YELLOW}$count${COLOR_NC} 個"
        echo -e "成功: ${COLOR_GREEN}$success_count${COLOR_NC} 個"

        if [ $success_count -eq $count ]; then
            echo -e "${COLOR_GREEN}✅ すべての環境変数の同期が完了しました！${COLOR_NC}"
        else
            local failed_count=$((count - success_count))
            echo -e "失敗: ${COLOR_RED}$failed_count${COLOR_NC} 個"
            echo -e "${COLOR_YELLOW}⚠️  一部の環境変数の同期に失敗しました。${COLOR_NC}"
        fi
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
}

# スクリプト実行
main "$@"
