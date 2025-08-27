#!/bin/bash

set -euo pipefail

# カラー定義
COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_CYAN='\033[0;36m'
COLOR_MAGENTA='\033[0;35m'
COLOR_NC='\033[0m'

# 環境変数の読み込み
load_env_file() {
    local file="$1"
    while IFS= read -r line || [ -n "$line" ]; do
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ "$line" =~ ^[[:space:]]*$ ]] && continue

        if [[ "$line" =~ ^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*=[[:space:]]*(.*)$ ]]; then
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"
            value="${value%\"}"
            value="${value#\"}"
            export "$key"="$value"
        fi
    done < "$file"
}

if [ -f .env ]; then
    load_env_file .env
elif [ -f .env.example ]; then
    load_env_file .env.example
fi

# ヘルプ表示
show_help() {
    echo -e "${COLOR_GREEN}Make vs Task 包括的比較ツール${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  ./comprehensive_comparison.sh [オプション] [タスク名]"
    echo ""
    echo -e "${COLOR_YELLOW}オプション:${COLOR_NC}"
    echo -e "  --all, -a          全てのタスクを比較"
    echo -e "  --performance, -p  パフォーマンス測定のみ"
    echo -e "  --memory, -m       メモリ使用量測定のみ"
    echo -e "  --complexity, -c   複雑性分析のみ"
    echo -e "  --report, -r       詳細レポート生成"
    echo -e "  --help, -h         このヘルプを表示"
    echo ""
    echo -e "${COLOR_YELLOW}利用可能なタスク:${COLOR_NC}"
    echo -e "  setup, mod, build, lint, test-unit, test-e2e, check"
    echo ""
    echo -e "${COLOR_BLUE}例:${COLOR_NC}"
    echo -e "  ./comprehensive_comparison.sh --all"
    echo -e "  ./comprehensive_comparison.sh --performance build"
    echo -e "  ./comprehensive_comparison.sh --memory test-unit"
}

# 必要なツールのチェック
check_dependencies() {
    local missing_tools=()

    if ! command -v make >/dev/null 2>&1; then
        missing_tools+=("make")
    fi

    if ! command -v task >/dev/null 2>&1; then
        missing_tools+=("task")
    fi

    if ! command -v bc >/dev/null 2>&1; then
        missing_tools+=("bc")
    fi

    if ! command -v ps >/dev/null 2>&1; then
        missing_tools+=("ps")
    fi

    if [ ${#missing_tools[@]} -gt 0 ]; then
        echo -e "${COLOR_RED}エラー: 以下のツールがインストールされていません: ${missing_tools[*]}${COLOR_NC}" >&2
        exit 1
    fi
}

# タスクマッピング
get_task_mapping() {
    local task_type="$1"
    case "$task_type" in
        "setup")
            echo "setup:setup"
            ;;
        "mod")
            echo "mod:mod"
            ;;
        "build")
            echo "build:build"
            ;;
        "lint")
            echo "lint:lint"
            ;;
        "test-unit")
            echo "test-unit:test:unit"
            ;;
        "test-e2e")
            echo "test-e2e:test:e2e"
            ;;
        "check")
            echo "check:check"
            ;;
        *)
            echo "unknown:unknown"
            ;;
    esac
}

# パフォーマンス測定
measure_performance() {
    local task_name="$1"
    local tool="$2"
    local command="$3"

    echo -e "${COLOR_CYAN}📊 ${tool} ${command} のパフォーマンス測定中...${COLOR_NC}"

    local temp_file=$(mktemp)
    local start_time=$(date +%s.%N)

    # メモリとCPU使用量を監視するバックグラウンドプロセス
    local monitor_pid=""
    if command -v top >/dev/null 2>&1; then
        (
            max_memory=0
            max_cpu=0
            while true; do
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    # macOS
                    memory=$(ps -o pid,rss,pcpu -p $$ 2>/dev/null | tail -1 | awk '{print $2}' || echo "0")
                    cpu=$(ps -o pid,rss,pcpu -p $$ 2>/dev/null | tail -1 | awk '{print $3}' || echo "0")
                else
                    # Linux
                    memory=$(ps -o pid,rss,pcpu -p $$ 2>/dev/null | tail -1 | awk '{print $2}' || echo "0")
                    cpu=$(ps -o pid,rss,pcpu -p $$ 2>/dev/null | tail -1 | awk '{print $3}' || echo "0")
                fi

                if (( $(echo "$memory > $max_memory" | bc -l) )); then
                    max_memory=$memory
                fi
                if (( $(echo "$cpu > $max_cpu" | bc -l) )); then
                    max_cpu=$cpu
                fi

                echo "$max_memory:$max_cpu" > "$temp_file.monitor"
                sleep 0.1
            done
        ) &
        monitor_pid=$!
    fi

    # 実際のコマンド実行
    local exit_code=0
    if [ "$tool" = "make" ]; then
        { time make "$command"; } 2>&1 | tee "$temp_file" || exit_code=$?
    else
        { time task "$command"; } 2>&1 | tee "$temp_file" || exit_code=$?
    fi

    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc -l)

    # 監視プロセスを停止
    if [ -n "$monitor_pid" ]; then
        kill $monitor_pid 2>/dev/null || true
        wait $monitor_pid 2>/dev/null || true
    fi

    # メモリとCPU使用量を取得
    local max_memory=0
    local max_cpu=0
    if [ -f "$temp_file.monitor" ]; then
        local monitor_data=$(cat "$temp_file.monitor" 2>/dev/null || echo "0:0")
        max_memory=$(echo "$monitor_data" | cut -d: -f1)
        max_cpu=$(echo "$monitor_data" | cut -d: -f2)
        rm -f "$temp_file.monitor"
    fi

    # 実行時間をtimeコマンドの出力から抽出
    local real_time="N/A"
    if grep -q "real" "$temp_file"; then
        real_time=$(grep "real" "$temp_file" | awk '{print $2}' | head -1)
        # 時間形式を秒に変換
        if [[ "$real_time" =~ ([0-9]+)m([0-9.]+)s ]]; then
            real_time=$(echo "scale=3; ${BASH_REMATCH[1]} * 60 + ${BASH_REMATCH[2]}" | bc -l)
        elif [[ "$real_time" =~ ([0-9.]+)s ]]; then
            real_time=${BASH_REMATCH[1]}
        fi
    fi

    # 結果を出力
    echo "$task_name,$tool,$real_time,$max_memory,$max_cpu,$exit_code"

    rm -f "$temp_file"
}

# 複雑性分析
analyze_complexity() {
    echo -e "${COLOR_MAGENTA}🔍 設定ファイルの複雑性分析${COLOR_NC}"

    local makefile_lines=$(wc -l < Makefile 2>/dev/null || echo "0")
    local taskfile_lines=$(wc -l < Taskfile.yml 2>/dev/null || echo "0")

    local makefile_targets=$(grep -c "^[a-zA-Z0-9_-]*:" Makefile 2>/dev/null || echo "0")
    local taskfile_tasks=$(grep -c "^  [a-zA-Z0-9_:-]*:" Taskfile.yml 2>/dev/null || echo "0")

    echo "設定ファイル,行数,ターゲット/タスク数"
    echo "Makefile,$makefile_lines,$makefile_targets"
    echo "Taskfile.yml,$taskfile_lines,$taskfile_tasks"
}

# 学習コスト分析
analyze_learning_cost() {
    echo -e "${COLOR_MAGENTA}📚 学習コスト分析${COLOR_NC}"

    # 構文の複雑さを分析
    local makefile_complexity=0
    local taskfile_complexity=0

    # Makefileの複雑性指標
    if [ -f Makefile ]; then
        # 特殊変数の使用数
        makefile_complexity=$((makefile_complexity + $(grep -c '\$[@<^+?*%]' Makefile 2>/dev/null || echo "0")))
        # 条件分岐の数
        makefile_complexity=$((makefile_complexity + $(grep -c 'ifeq\|ifneq\|ifdef\|ifndef' Makefile 2>/dev/null || echo "0")))
        # シェルコマンドの複雑さ
        makefile_complexity=$((makefile_complexity + $(grep -c '\\$' Makefile 2>/dev/null || echo "0")))
    fi

    # Taskfileの複雑性指標
    if [ -f Taskfile.yml ]; then
        # YAML構造の深さ
        taskfile_complexity=$((taskfile_complexity + $(grep -c '^    ' Taskfile.yml 2>/dev/null || echo "0")))
        # 条件分岐の数
        taskfile_complexity=$((taskfile_complexity + $(grep -c 'status:\|preconditions:' Taskfile.yml 2>/dev/null || echo "0")))
    fi

    echo "ツール,複雑性スコア,可読性"
    echo "Make,$makefile_complexity,低"
    echo "Task,$taskfile_complexity,高"
}

# レポート生成
generate_report() {
    local results_file="$1"
    local report_file="performance_comparison_report_$(date +%Y%m%d_%H%M%S).md"

    echo -e "${COLOR_GREEN}📋 詳細レポートを生成中: $report_file${COLOR_NC}"

    cat > "$report_file" << EOF
# Make vs Task 包括的比較レポート

**生成日時:** $(date '+%Y年%m月%d日 %H:%M:%S')

## 実行環境
- **OS:** $(uname -s)
- **アーキテクチャ:** $(uname -m)
- **Make バージョン:** $(make --version | head -1 2>/dev/null || echo "不明")
- **Task バージョン:** $(task --version 2>/dev/null || echo "不明")

## パフォーマンス比較結果

| タスク | ツール | 実行時間(秒) | メモリ使用量(KB) | CPU使用率(%) | 終了コード |
|--------|--------|-------------|-----------------|-------------|-----------|
EOF

    if [ -f "$results_file" ]; then
        while IFS=',' read -r task tool time memory cpu exit_code; do
            echo "| $task | $tool | $time | $memory | $cpu | $exit_code |" >> "$report_file"
        done < "$results_file"
    fi

    cat >> "$report_file" << EOF

## 複雑性分析

EOF

    analyze_complexity >> "$report_file"

    cat >> "$report_file" << EOF

## 学習コスト分析

EOF

    analyze_learning_cost >> "$report_file"

    cat >> "$report_file" << EOF

## 総合評価

### パフォーマンス
- 実行速度、メモリ使用量、CPU使用率の観点から評価

### 開発体験
- 設定の書きやすさ、可読性、メンテナンス性

### 学習コスト
- 新規参加者の学習のしやすさ

## 推奨事項

プロジェクトの特性に応じて適切なツールを選択することを推奨します。

EOF

    echo -e "${COLOR_GREEN}✅ レポートが生成されました: $report_file${COLOR_NC}"
}

# メイン処理
main() {
    local mode="all"
    local specific_task=""
    local generate_report_flag=false

    # 引数解析
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all|-a)
                mode="all"
                shift
                ;;
            --performance|-p)
                mode="performance"
                shift
                ;;
            --memory|-m)
                mode="memory"
                shift
                ;;
            --complexity|-c)
                mode="complexity"
                shift
                ;;
            --report|-r)
                generate_report_flag=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                specific_task="$1"
                shift
                ;;
        esac
    done

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  🚀 ${COLOR_GREEN}Make vs Task 包括的比較ツール${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo ""

    check_dependencies

    local results_file=$(mktemp)
    echo "task,tool,time,memory,cpu,exit_code" > "$results_file"

    if [ "$mode" = "all" ] || [ "$mode" = "performance" ] || [ "$mode" = "memory" ]; then
        local tasks_to_test=("mod" "build" "lint" "test-unit")

        if [ -n "$specific_task" ]; then
            tasks_to_test=("$specific_task")
        fi

        for task in "${tasks_to_test[@]}"; do
            local mapping=$(get_task_mapping "$task")
            local make_target=$(echo "$mapping" | cut -d: -f1)
            local task_target=$(echo "$mapping" | cut -d: -f2)

            if [ "$make_target" != "unknown" ]; then
                echo -e "${COLOR_YELLOW}🔄 タスク: $task${COLOR_NC}"

                # Make実行
                measure_performance "$task" "make" "$make_target" >> "$results_file"

                # Task実行
                measure_performance "$task" "task" "$task_target" >> "$results_file"

                echo ""
            fi
        done

        # 結果表示
        echo -e "${COLOR_GREEN}📊 パフォーマンス比較結果${COLOR_NC}"
        echo "----------------------------------------"
        column -t -s',' "$results_file"
        echo ""
    fi

    if [ "$mode" = "all" ] || [ "$mode" = "complexity" ]; then
        analyze_complexity
        echo ""
        analyze_learning_cost
        echo ""
    fi

    if [ "$generate_report_flag" = true ]; then
        generate_report "$results_file"
    fi

    rm -f "$results_file"

    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
    echo -e "  ${COLOR_GREEN}比較分析が完了しました${COLOR_NC}"
    echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
}

# スクリプト実行
main "$@"
