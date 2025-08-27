#!/bin/bash

set -euo pipefail

# カラー定義
COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_CYAN='\033[0;36m'
COLOR_NC='\033[0m'

# ヘルプ表示
show_help() {
    echo -e "${COLOR_GREEN}Make vs Task クイック比較ツール${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  ./quick_comparison.sh [タスク名]"
    echo ""
    echo -e "${COLOR_YELLOW}利用可能なタスク:${COLOR_NC}"
    echo -e "  mod      - Goモジュール管理"
    echo -e "  build    - アプリケーションビルド"
    echo -e "  lint     - コード静的解析"
    echo -e "  unit     - ユニットテスト"
    echo -e "  check    - 全体チェック（lint + test + build）"
    echo ""
    echo -e "${COLOR_BLUE}例:${COLOR_NC}"
    echo -e "  ./quick_comparison.sh build"
    echo -e "  ./quick_comparison.sh          # 全タスクを比較"
}

# 必要なツールのチェック
check_tools() {
    local missing=()

    if ! command -v make >/dev/null 2>&1; then
        missing+=("make")
    fi

    if ! command -v task >/dev/null 2>&1; then
        missing+=("task")
    fi

    if ! command -v bc >/dev/null 2>&1; then
        missing+=("bc")
    fi

    if [ ${#missing[@]} -gt 0 ]; then
        echo -e "${COLOR_RED}❌ 以下のツールが必要です: ${missing[*]}${COLOR_NC}" >&2
        exit 1
    fi
}

# タスクマッピング
get_commands() {
    local task="$1"
    case "$task" in
        "mod")
            echo "mod:mod"
            ;;
        "build")
            echo "build:build"
            ;;
        "lint")
            echo "lint:lint"
            ;;
        "unit")
            echo "test-unit:test:unit"
            ;;
        "check")
            echo "check:check"
            ;;
        *)
            echo "unknown:unknown"
            ;;
    esac
}

# 単一タスクの比較
compare_task() {
    local task_name="$1"
    local make_cmd="$2"
    local task_cmd="$3"

    echo -e "${COLOR_CYAN}🔄 比較中: $task_name${COLOR_NC}"
    echo "----------------------------------------"

    # Make実行
    echo -e "${COLOR_BLUE}▶️ make $make_cmd${COLOR_NC}"
    local make_start=$(date +%s.%N)
    local make_exit=0
    make "$make_cmd" >/dev/null 2>&1 || make_exit=$?
    local make_end=$(date +%s.%N)
    local make_time=$(echo "$make_end - $make_start" | bc -l)

    # Task実行
    echo -e "${COLOR_BLUE}▶️ task $task_cmd${COLOR_NC}"
    local task_start=$(date +%s.%N)
    local task_exit=0
    task "$task_cmd" >/dev/null 2>&1 || task_exit=$?
    local task_end=$(date +%s.%N)
    local task_time=$(echo "$task_end - $task_start" | bc -l)

    # 結果表示
    printf "%-10s: %8.3fs (exit: %d)\n" "Make" "$make_time" "$make_exit"
    printf "%-10s: %8.3fs (exit: %d)\n" "Task" "$task_time" "$task_exit"

    # 勝者判定
    if (( $(echo "$make_time < $task_time" | bc -l) )); then
        local diff=$(echo "$task_time - $make_time" | bc -l)
        local percent=$(echo "scale=1; $diff / $make_time * 100" | bc -l)
        echo -e "${COLOR_GREEN}🏆 Make が ${diff}s (${percent}%) 高速${COLOR_NC}"
    elif (( $(echo "$task_time < $make_time" | bc -l) )); then
        local diff=$(echo "$make_time - $task_time" | bc -l)
        local percent=$(echo "scale=1; $diff / $task_time * 100" | bc -l)
        echo -e "${COLOR_GREEN}🏆 Task が ${diff}s (${percent}%) 高速${COLOR_NC}"
    else
        echo -e "${COLOR_YELLOW}🤝 ほぼ同等${COLOR_NC}"
    fi

    echo ""
}

# メイン処理
main() {
    if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
        show_help
        exit 0
    fi

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  ⚡ ${COLOR_GREEN}Make vs Task クイック比較${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    check_tools

    local target_task="${1:-}"
    local tasks_to_run=()

    if [ -n "$target_task" ]; then
        tasks_to_run=("$target_task")
    else
        tasks_to_run=("mod" "build" "lint" "unit")
    fi

    local total_make_time=0
    local total_task_time=0

    for task in "${tasks_to_run[@]}"; do
        local commands=$(get_commands "$task")
        local make_cmd=$(echo "$commands" | cut -d: -f1)
        local task_cmd=$(echo "$commands" | cut -d: -f2)

        if [ "$make_cmd" != "unknown" ]; then
            compare_task "$task" "$make_cmd" "$task_cmd"
        else
            echo -e "${COLOR_RED}❌ 不明なタスク: $task${COLOR_NC}"
        fi
    done

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  ${COLOR_GREEN}クイック比較完了${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}💡 詳細な分析には以下を実行:${COLOR_NC}"
    echo -e "  ./scripts/comprehensive_comparison.sh --all --report"
}

# スクリプト実行
main "$@"
