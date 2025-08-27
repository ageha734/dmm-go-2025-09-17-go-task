#!/bin/bash

set -euo pipefail

COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

DEFAULT_ENV_FILE=".env"

show_help() {
    echo -e "${COLOR_GREEN}GitHub Repository Secrets 同期スクリプト${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  $0 [OPTIONS]"
    echo ""
    echo -e "${COLOR_YELLOW}オプション:${COLOR_NC}"
    echo -e "  -f, --file FILE        環境変数ファイルを指定 (デフォルト: ${DEFAULT_ENV_FILE})"
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
    echo -e "  $0 -f .env -r owner/repo"
    echo -e "  $0 --dry-run"
    echo -e "  $0 --help"
}

parse_env_file() {
    local file="$1"

    while IFS= read -r line || [ -n "$line" ]; do
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ "$line" =~ ^[[:space:]]*$ ]] && continue

        if [[ "$line" =~ ^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*=[[:space:]]*(.*)$ ]]; then
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"

            value="${value%\"}"
            value="${value#\"}"
            value="${value%\'}"
            value="${value#\'}"

            echo "$key=$value"
        fi
    done < "$file"
}

check_github_cli() {
    if ! command -v gh >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: GitHub CLI (gh) がインストールされていません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}インストール方法: https://cli.github.com/${COLOR_NC}" >&2
        exit 1
    fi
}

check_github_auth() {
    if ! gh auth status -h github.com >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: GitHub CLI の認証が設定されていません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}認証方法: gh auth login${COLOR_NC}" >&2
        exit 1
    fi
}

check_repository() {
    local repo="$1"

    if ! gh repo view "$repo" >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: リポジトリ '$repo' にアクセスできません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}リポジトリ名が正しいか、アクセス権限があるか確認してください。${COLOR_NC}" >&2
        exit 1
    fi
}

set_secret() {
    local repo="$1"
    local key="$2"
    local value="$3"
    local dry_run="$4"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} $key をRepository secretsに設定予定"
        return
    fi

    echo -e "${COLOR_BLUE}$key を設定中...${COLOR_NC}"

    echo "$value" | gh secret set "$key" --repo "$repo"

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ $key を設定しました${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ $key の設定に失敗しました${COLOR_NC}" >&2
        return 1
    fi
}

main() {
    local env_file="$DEFAULT_ENV_FILE"
    local repo=""
    local token=""
    local dry_run="false"

    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--file)
                env_file="$2"
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

    if [ -z "$repo" ]; then
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

    if [ ! -f "$env_file" ]; then
        echo -e "${COLOR_RED}エラー: 環境変数ファイル '$env_file' が見つかりません。${COLOR_NC}" >&2
        exit 1
    fi

    check_github_cli
    check_github_auth

    check_repository "$repo"

    if [ -n "$token" ]; then
        export GH_TOKEN="$token"
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  🚀 ${COLOR_GREEN}GitHub Repository Secrets 同期${COLOR_NC}"
    echo -e "  ファイル: ${COLOR_YELLOW}$env_file${COLOR_NC}"
    echo -e "  リポジトリ: ${COLOR_YELLOW}$repo${COLOR_NC}"
    if [ "$dry_run" = "true" ]; then
        echo -e "  モード: ${COLOR_YELLOW}DRY RUN${COLOR_NC}"
    fi
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    echo -e "${COLOR_BLUE}環境変数を解析中...${COLOR_NC}"

    local count=0
    local success_count=0

    while IFS='=' read -r key value; do
        if [ -n "$key" ] && [ -n "$value" ]; then
            count=$((count + 1))

            if set_secret "$repo" "$key" "$value" "$dry_run"; then
                success_count=$((success_count + 1))
            fi
        fi
    done < <(parse_env_file "$env_file")

    echo ""
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  📊 ${COLOR_GREEN}同期結果${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"

    if [ "$dry_run" = "true" ]; then
        echo -e "処理予定のsecrets: ${COLOR_YELLOW}$count${COLOR_NC} 個"
        echo -e "${COLOR_YELLOW}実際の同期を実行するには --dry-run オプションを外してください。${COLOR_NC}"
    else
        echo -e "処理したsecrets: ${COLOR_YELLOW}$count${COLOR_NC} 個"
        echo -e "成功: ${COLOR_GREEN}$success_count${COLOR_NC} 個"

        if [ $success_count -eq $count ]; then
            echo -e "${COLOR_GREEN}✅ すべてのsecretsの同期が完了しました！${COLOR_NC}"
        else
            local failed_count=$((count - success_count))
            echo -e "失敗: ${COLOR_RED}$failed_count${COLOR_NC} 個"
            echo -e "${COLOR_YELLOW}⚠️  一部のsecretsの同期に失敗しました。${COLOR_NC}"
        fi
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
}

main "$@"
