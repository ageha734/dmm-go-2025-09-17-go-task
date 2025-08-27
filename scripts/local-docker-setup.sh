#!/bin/bash

set -euo pipefail

COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

REGISTRY="ghcr.io"
DEFAULT_TAG="latest"

show_help() {
    echo -e "${COLOR_GREEN}Local Docker Setup Script${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  $0 [OPTIONS]"
    echo ""
    echo -e "${COLOR_YELLOW}オプション:${COLOR_NC}"
    echo -e "  -r, --repo REPO        リポジトリ名を指定 (owner/repo形式)"
    echo -e "  -t, --tag TAG          イメージタグを指定 (デフォルト: ${DEFAULT_TAG})"
    echo -e "  -m, --mode MODE        実行モード (build|pull|status) (デフォルト: build)"
    echo -e "  -u, --username USER    GitHub ユーザー名を指定 (pullモード時)"
    echo -e "  -p, --token TOKEN      GitHub Personal Access Tokenを指定 (pullモード時)"
    echo -e "  -d, --dry-run          実際には実行せず、実行予定の内容を表示"
    echo -e "  -h, --help             このヘルプを表示"
    echo ""
    echo -e "${COLOR_YELLOW}実行モード:${COLOR_NC}"
    echo -e "  build                  ローカルでDockerイメージをビルドしてcompose.ymlを更新"
    echo -e "  pull                   GitHub Container Registryからイメージをプルしてcompose.ymlを更新"
    echo -e "  status                 現在のcompose.ymlの設定を表示"
    echo ""
    echo -e "${COLOR_YELLOW}例:${COLOR_NC}"
    echo -e "  $0 --mode build                    # ローカルビルド"
    echo -e "  $0 --mode pull -u username -p token # レジストリからプル"
    echo -e "  $0 --mode status                   # 現在の設定を確認"
    echo -e "  $0 --dry-run"
}

check_docker() {
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: Docker がインストールされていません。${COLOR_NC}" >&2
        exit 1
    fi
}

docker_login() {
    local username="$1"
    local token="$2"
    local dry_run="$3"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Docker login to ${REGISTRY}"
        return
    fi

    echo -e "${COLOR_BLUE}GitHub Container Registry にログイン中...${COLOR_NC}"
    echo "$token" | docker login "$REGISTRY" -u "$username" --password-stdin

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ ログインに成功しました${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ ログインに失敗しました${COLOR_NC}" >&2
        exit 1
    fi
}

build_local_image() {
    local dry_run="$1"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Docker build local image: api:latest"
        return
    fi

    echo -e "${COLOR_BLUE}ローカルでDockerイメージをビルド中...${COLOR_NC}"

    # 現在のアーキテクチャを検出
    local arch=$(uname -m)
    case $arch in
        x86_64)
            platform="linux/amd64"
            ;;
        arm64|aarch64)
            platform="linux/arm64"
            ;;
        *)
            platform="linux/amd64"
            echo -e "${COLOR_YELLOW}⚠️  不明なアーキテクチャ ($arch)、linux/amd64を使用します${COLOR_NC}"
            ;;
    esac

    echo -e "${COLOR_BLUE}対象プラットフォーム: $platform${COLOR_NC}"

    # Docker Buildxが利用可能な場合は使用、そうでなければ通常のbuildを使用
    if docker buildx version >/dev/null 2>&1; then
        docker buildx build --platform "$platform" -t api:latest --load .
    else
        docker build -t api:latest .
    fi

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ ビルドに成功しました${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ ビルドに失敗しました${COLOR_NC}" >&2
        exit 1
    fi
}

pull_registry_image() {
    local repo="$1"
    local tag="$2"
    local dry_run="$3"

    local image_name="${REGISTRY}/${repo}:${tag}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Docker pull: $image_name"
        return
    fi

    echo -e "${COLOR_BLUE}GitHub Container Registry からイメージをプル中...${COLOR_NC}"
    docker pull "$image_name"

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ プルに成功しました: $image_name${COLOR_NC}"

        # ローカルタグを付与
        docker tag "$image_name" api:latest
        echo -e "${COLOR_GREEN}✅ ローカルタグ api:latest を付与しました${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ プルに失敗しました${COLOR_NC}" >&2
        exit 1
    fi
}

update_compose_for_local() {
    local dry_run="$1"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Update compose.yml for local development"
        return
    fi

    echo -e "${COLOR_BLUE}compose.yml をローカル開発用に更新中...${COLOR_NC}"

    # compose.ymlのapiサービスのimageをローカル用に戻す
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|image: ghcr\.io/[^:]*:[^[:space:]]*|image: api:latest|g" compose.yml
    else
        # Linux
        sed -i "s|image: ghcr\.io/[^:]*:[^[:space:]]*|image: api:latest|g" compose.yml
    fi

    echo -e "${COLOR_GREEN}✅ compose.yml をローカル開発用に更新しました${COLOR_NC}"
}

update_compose_with_registry_image() {
    local repo="$1"
    local tag="$2"
    local dry_run="$3"

    local image_name="${REGISTRY}/${repo}:${tag}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Update compose.yml with registry image: $image_name"
        return
    fi

    echo -e "${COLOR_BLUE}compose.yml をレジストリイメージ用に更新中...${COLOR_NC}"

    # compose.ymlのapiサービスのimageをレジストリイメージに更新
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|image: api:latest|image: $image_name|g" compose.yml
    else
        # Linux
        sed -i "s|image: api:latest|image: $image_name|g" compose.yml
    fi

    echo -e "${COLOR_GREEN}✅ compose.yml をレジストリイメージ用に更新しました${COLOR_NC}"
}

show_status() {
    echo -e "${COLOR_BLUE}現在のcompose.yml設定:${COLOR_NC}"
    echo ""

    if [ -f "compose.yml" ]; then
        # apiサービスのimage設定を抽出
        local image_line=$(grep -A 10 "api:" compose.yml | grep "image:" | head -1)
        if [ -n "$image_line" ]; then
            echo -e "  API Image: ${COLOR_YELLOW}$(echo "$image_line" | sed 's/.*image: *//')${COLOR_NC}"

            if echo "$image_line" | grep -q "ghcr.io"; then
                echo -e "  モード: ${COLOR_BLUE}Registry Image${COLOR_NC}"
            else
                echo -e "  モード: ${COLOR_GREEN}Local Build${COLOR_NC}"
            fi
        else
            echo -e "  ${COLOR_RED}API imageの設定が見つかりません${COLOR_NC}"
        fi
    else
        echo -e "  ${COLOR_RED}compose.yml が見つかりません${COLOR_NC}"
    fi

    echo ""
    echo -e "${COLOR_BLUE}利用可能なDockerイメージ:${COLOR_NC}"
    docker images | grep -E "(api|ghcr.io)" || echo -e "  ${COLOR_YELLOW}関連するイメージが見つかりません${COLOR_NC}"
}

main() {
    local repo=""
    local tag="$DEFAULT_TAG"
    local mode="build"
    local username=""
    local token=""
    local dry_run="false"

    while [[ $# -gt 0 ]]; do
        case $1 in
            -r|--repo)
                repo="$2"
                shift 2
                ;;
            -t|--tag)
                tag="$2"
                shift 2
                ;;
            -m|--mode)
                mode="$2"
                shift 2
                ;;
            -u|--username)
                username="$2"
                shift 2
                ;;
            -p|--token)
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

    # statusモードの場合は早期リターン
    if [ "$mode" = "status" ]; then
        show_status
        exit 0
    fi

    # リポジトリ名の自動検出（pullモード時のみ必要）
    if [ "$mode" = "pull" ] && [ -z "$repo" ]; then
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

    check_docker

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  🚀 ${COLOR_GREEN}Local Docker Setup${COLOR_NC}"
    echo -e "  モード: ${COLOR_YELLOW}$mode${COLOR_NC}"
    if [ "$mode" = "pull" ]; then
        echo -e "  リポジトリ: ${COLOR_YELLOW}$repo${COLOR_NC}"
        echo -e "  タグ: ${COLOR_YELLOW}$tag${COLOR_NC}"
    fi
    if [ "$dry_run" = "true" ]; then
        echo -e "  実行: ${COLOR_YELLOW}DRY RUN${COLOR_NC}"
    fi
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    case "$mode" in
        "build")
            build_local_image "$dry_run"
            update_compose_for_local "$dry_run"
            ;;
        "pull")
            # GitHub認証情報の確認
            if [ -z "$username" ] && [ "$dry_run" = "false" ]; then
                if [ -n "${GITHUB_ACTOR:-}" ]; then
                    username="$GITHUB_ACTOR"
                    echo -e "${COLOR_BLUE}GitHub Actor を使用: $username${COLOR_NC}"
                else
                    # gitコマンドでユーザー名を自動取得を試行
                    if git config user.name >/dev/null 2>&1; then
                        git_username=$(git config user.name)
                        echo -e "${COLOR_BLUE}Git設定からユーザー名を自動検出: $git_username${COLOR_NC}"
                        username="$git_username"
                    elif git remote get-url origin >/dev/null 2>&1; then
                        # リモートURLからユーザー名を抽出
                        repo_url=$(git remote get-url origin)
                        if [[ "$repo_url" =~ github\.com[:/]([^/]+)/[^/]+(\.git)?$ ]]; then
                            git_username="${BASH_REMATCH[1]}"
                            echo -e "${COLOR_BLUE}Git リモートURLからユーザー名を自動検出: $git_username${COLOR_NC}"
                            username="$git_username"
                        fi
                    fi

                    if [ -z "$username" ]; then
                        echo -e "${COLOR_RED}エラー: GitHub ユーザー名が指定されていません。${COLOR_NC}" >&2
                        echo -e "${COLOR_YELLOW}以下のいずれかの方法でユーザー名を指定してください:${COLOR_NC}" >&2
                        echo -e "${COLOR_YELLOW}  1. -u オプションでユーザー名を指定${COLOR_NC}" >&2
                        echo -e "${COLOR_YELLOW}  2. git config user.name でGitユーザー名を設定${COLOR_NC}" >&2
                        echo -e "${COLOR_YELLOW}  3. GITHUB_ACTOR環境変数を設定${COLOR_NC}" >&2
                        exit 1
                    fi
                fi
            fi

            if [ -z "$token" ] && [ "$dry_run" = "false" ]; then
                if [ -n "${GITHUB_TOKEN:-}" ]; then
                    token="$GITHUB_TOKEN"
                    echo -e "${COLOR_BLUE}GITHUB_TOKEN を使用${COLOR_NC}"
                else
                    echo -e "${COLOR_RED}エラー: GitHub Personal Access Token が指定されていません。${COLOR_NC}" >&2
                    echo -e "${COLOR_YELLOW}-p オプションでトークンを指定するか、GITHUB_TOKEN環境変数を設定してください${COLOR_NC}" >&2
                    exit 1
                fi
            fi

            if [ "$dry_run" = "false" ]; then
                docker_login "$username" "$token" "$dry_run"
            fi
            pull_registry_image "$repo" "$tag" "$dry_run"
            update_compose_with_registry_image "$repo" "$tag" "$dry_run"
            ;;
        *)
            echo -e "${COLOR_RED}エラー: 不明なモード '$mode'${COLOR_NC}" >&2
            echo -e "${COLOR_YELLOW}利用可能なモード: build, pull, status${COLOR_NC}" >&2
            exit 1
            ;;
    esac

    echo ""
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  📊 ${COLOR_GREEN}完了${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}実際の実行を行うには --dry-run オプションを外してください。${COLOR_NC}"
    else
        case "$mode" in
            "build")
                echo -e "${COLOR_GREEN}✅ ローカルビルドが完了しました！${COLOR_NC}"
                echo -e "次のステップ: ${COLOR_BLUE}docker-compose up${COLOR_NC} でアプリケーションを起動できます"
                ;;
            "pull")
                echo -e "${COLOR_GREEN}✅ レジストリからのプルが完了しました！${COLOR_NC}"
                echo -e "次のステップ: ${COLOR_BLUE}docker-compose up${COLOR_NC} でアプリケーションを起動できます"
                ;;
        esac
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
}

main "$@"
