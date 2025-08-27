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
    echo -e "${COLOR_GREEN}GitHub Container Registry Docker Image Push Script${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  $0 [OPTIONS]"
    echo ""
    echo -e "${COLOR_YELLOW}オプション:${COLOR_NC}"
    echo -e "  -r, --repo REPO        リポジトリ名を指定 (owner/repo形式)"
    echo -e "  -t, --tag TAG          イメージタグを指定 (デフォルト: ${DEFAULT_TAG})"
    echo -e "  -u, --username USER    GitHub ユーザー名を指定"
    echo -e "  -p, --token TOKEN      GitHub Personal Access Tokenを指定"
    echo -e "  -d, --dry-run          実際には実行せず、実行予定の内容を表示"
    echo -e "  -h, --help             このヘルプを表示"
    echo ""
    echo -e "${COLOR_YELLOW}前提条件:${COLOR_NC}"
    echo -e "  - Docker がインストールされていること"
    echo -e "  - GitHub Personal Access Token (packages:write権限) が設定されていること"
    echo -e "  - リポジトリへの write 権限があること"
    echo ""
    echo -e "${COLOR_YELLOW}例:${COLOR_NC}"
    echo -e "  $0 -r owner/repo -u username -p ghp_token"
    echo -e "  $0 --dry-run"
    echo -e "  $0 --help"
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

build_and_push_image() {
    local repo="$1"
    local tag="$2"
    local dry_run="$3"

    local image_name="${REGISTRY}/${repo}:${tag}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Docker buildx build and push (multi-arch): $image_name"
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Platforms: linux/amd64,linux/arm64"
        return
    fi

    # Docker Buildxが利用可能かチェック
    if ! docker buildx version >/dev/null 2>&1; then
        echo -e "${COLOR_RED}エラー: Docker Buildx が利用できません。${COLOR_NC}" >&2
        echo -e "${COLOR_YELLOW}Docker Buildx をインストールするか、Docker Desktop を使用してください。${COLOR_NC}" >&2
        exit 1
    fi

    # マルチアーキテクチャビルダーを作成または使用
    local builder_name="multiarch-builder"
    if ! docker buildx inspect "$builder_name" >/dev/null 2>&1; then
        echo -e "${COLOR_BLUE}マルチアーキテクチャビルダーを作成中...${COLOR_NC}"
        docker buildx create --name "$builder_name" --driver docker-container --use
        docker buildx inspect --bootstrap
    else
        echo -e "${COLOR_BLUE}既存のマルチアーキテクチャビルダーを使用: $builder_name${COLOR_NC}"
        docker buildx use "$builder_name"
    fi

    echo -e "${COLOR_BLUE}Docker イメージをマルチアーキテクチャでビルド・プッシュ中...${COLOR_NC}"
    echo -e "${COLOR_BLUE}対象プラットフォーム: linux/amd64, linux/arm64${COLOR_NC}"

    docker buildx build \
        --platform linux/amd64,linux/arm64 \
        --tag "$image_name" \
        --push \
        .

    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ マルチアーキテクチャビルド・プッシュに成功しました: $image_name${COLOR_NC}"
        echo -e "${COLOR_GREEN}  - linux/amd64${COLOR_NC}"
        echo -e "${COLOR_GREEN}  - linux/arm64${COLOR_NC}"
    else
        echo -e "${COLOR_RED}❌ マルチアーキテクチャビルド・プッシュに失敗しました${COLOR_NC}" >&2
        exit 1
    fi
}

update_compose_file() {
    local repo="$1"
    local tag="$2"
    local dry_run="$3"

    local image_name="${REGISTRY}/${repo}:${tag}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}[DRY RUN]${COLOR_NC} Update compose.yml with image: $image_name"
        return
    fi

    echo -e "${COLOR_BLUE}compose.yml を更新中...${COLOR_NC}"

    # compose.ymlのapiサービスのimageを更新
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|image: api:latest|image: $image_name|g" compose.yml
    else
        # Linux
        sed -i "s|image: api:latest|image: $image_name|g" compose.yml
    fi

    echo -e "${COLOR_GREEN}✅ compose.yml を更新しました${COLOR_NC}"
}

main() {
    local repo=""
    local tag="$DEFAULT_TAG"
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

    # リポジトリ名の自動検出
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

    check_docker

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  🚀 ${COLOR_GREEN}GitHub Container Registry Push${COLOR_NC}"
    echo -e "  リポジトリ: ${COLOR_YELLOW}$repo${COLOR_NC}"
    echo -e "  タグ: ${COLOR_YELLOW}$tag${COLOR_NC}"
    echo -e "  イメージ: ${COLOR_YELLOW}${REGISTRY}/${repo}:${tag}${COLOR_NC}"
    if [ "$dry_run" = "true" ]; then
        echo -e "  モード: ${COLOR_YELLOW}DRY RUN${COLOR_NC}"
    fi
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    if [ "$dry_run" = "false" ]; then
        docker_login "$username" "$token" "$dry_run"
    fi

    build_and_push_image "$repo" "$tag" "$dry_run"
    update_compose_file "$repo" "$tag" "$dry_run"

    echo ""
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  📊 ${COLOR_GREEN}完了${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"

    if [ "$dry_run" = "true" ]; then
        echo -e "${COLOR_YELLOW}実際のプッシュを実行するには --dry-run オプションを外してください。${COLOR_NC}"
    else
        echo -e "${COLOR_GREEN}✅ Docker イメージのプッシュが完了しました！${COLOR_NC}"
        echo -e "イメージ: ${COLOR_YELLOW}${REGISTRY}/${repo}:${tag}${COLOR_NC}"
        echo ""
        echo -e "${COLOR_BLUE}次のステップ:${COLOR_NC}"
        echo -e "1. GitHub Actions でこのイメージを使用できます"
        echo -e "2. ローカルでも 'docker-compose up' でこのイメージを使用できます"
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
}

main "$@"
