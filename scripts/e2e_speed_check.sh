#!/bin/bash

set -euo pipefail

COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

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
        elif [[ "$line" =~ ^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*:[[:space:]]*\"?([^\"]*)\"?[[:space:]]*$ ]]; then
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"
            export "$key"="$value"
        fi
    done < "$file"
}

if [ -f .env ]; then
    load_env_file .env
elif [ -f .env.example ]; then
    load_env_file .env.example
fi

SERVICE_NAME="${1:-}"

if [ "${SERVICE_NAME}" = "--help" ] || [ "${SERVICE_NAME}" = "-h" ]; then
    echo -e "${COLOR_GREEN}E2Eテスト速度比較スクリプト${COLOR_NC}"
    echo ""
    echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
    echo -e "  ./e2e_speed_check.sh           # 全サービスのテストを実行"
    echo -e "  ./e2e_speed_check.sh home      # homeサービスのみテストを実行"
    echo -e "  ./e2e_speed_check.sh user      # userサービスのみテストを実行"
    echo -e "  ./e2e_speed_check.sh --help    # このヘルプを表示"
    echo ""
    echo -e "${COLOR_BLUE}このスクリプトは以下を比較します:${COLOR_NC}"
    echo -e "  - make test-e2e の実行時間"
    echo -e "  - task test:e2e の実行時間"
    exit 0
fi

AVAILABLE_SERVICES=$(find ./e2e -maxdepth 1 -mindepth 1 -type d -exec basename {} \; | sort | uniq | tr '\n' ' ')

echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  🚀 ${COLOR_GREEN}E2Eテスト実行時間比較 (Make vs Task)${COLOR_NC}"
if [ -n "${SERVICE_NAME}" ]; then
    echo -e "  対象サービス: ${COLOR_GREEN}${SERVICE_NAME}${COLOR_NC}"
else
    echo -e "  対象サービス: ${COLOR_GREEN}全サービス${COLOR_NC}"
fi
echo -e "  利用可能なサービス: ${COLOR_YELLOW}${AVAILABLE_SERVICES}${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo ""

if [ -n "${SERVICE_NAME}" ]; then
    if ! echo "${AVAILABLE_SERVICES}" | grep -q -w "${SERVICE_NAME}"; then
        echo -e "${COLOR_RED}エラー: 指定されたサービス '${SERVICE_NAME}' が見つかりません。${COLOR_NC}" >&2
        echo -e "利用可能なサービス: ${AVAILABLE_SERVICES}" >&2
        exit 1
    fi
fi

echo -e "🔍 ${COLOR_BLUE}必要なツールをチェック中...${COLOR_NC}"
MISSING_TOOLS=()

if ! command -v make >/dev/null 2>&1; then
    MISSING_TOOLS+=("make")
fi

if ! command -v task >/dev/null 2>&1; then
    MISSING_TOOLS+=("task")
fi

if ! command -v bc >/dev/null 2>&1; then
    MISSING_TOOLS+=("bc")
fi

if [ ${#MISSING_TOOLS[@]} -gt 0 ]; then
    echo -e "${COLOR_RED}エラー: 以下のツールがインストールされていません: ${MISSING_TOOLS[*]}${COLOR_NC}" >&2
    echo -e "${COLOR_YELLOW}注意: hurl と shlack がない場合、テストは失敗しますが実行時間の比較は可能です。${COLOR_NC}"
    echo ""
fi

echo -e "🔍 ${COLOR_BLUE}Docker Composeサービスの状態をチェック中...${COLOR_NC}"
if ! docker compose ps | grep -q "Up"; then
    echo -e "${COLOR_YELLOW}⚠️  Docker Composeサービスが起動していません。起動します...${COLOR_NC}"
    docker compose up -d
    echo -e "⏳ ${COLOR_BLUE}サービスの起動を待機中...${COLOR_NC}"
    sleep 10
fi
echo -e "✅ ${COLOR_GREEN}Docker Composeサービスが起動しています。${COLOR_NC}"
echo ""

MAKE_RESULT=$(mktemp)
TASK_RESULT=$(mktemp)

cleanup() {
    rm -f "${MAKE_RESULT}" "${TASK_RESULT}"
}
trap cleanup EXIT

# --- Makefileでのe2eテスト実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'make test-e2e' ...${COLOR_NC}"
echo "-------------------------------------------------"
START_TIME=$(date +%s.%N)
MAKE_EXIT_CODE=0
if [ -n "${SERVICE_NAME}" ]; then
    { time make test-e2e service="${SERVICE_NAME}"; } 2>&1 | tee "${MAKE_RESULT}" || MAKE_EXIT_CODE=$?
else
    { time make test-e2e; } 2>&1 | tee "${MAKE_RESULT}" || MAKE_EXIT_CODE=$?
fi
END_TIME=$(date +%s.%N)
MAKE_DURATION=$(echo "${END_TIME} - ${START_TIME}" | bc -l)
echo "-------------------------------------------------"
if [ ${MAKE_EXIT_CODE} -eq 0 ]; then
    echo -e "✅ ${COLOR_BLUE}'make test-e2e' finished successfully. Duration: ${MAKE_DURATION}s${COLOR_NC}"
else
    echo -e "⚠️  ${COLOR_YELLOW}'make test-e2e' finished with errors (exit code: ${MAKE_EXIT_CODE}). Duration: ${MAKE_DURATION}s${COLOR_NC}"
fi
echo ""
echo ""

# --- Taskfileでのe2eテスト実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'task test:e2e' ...${COLOR_NC}"
echo "-------------------------------------------------"
START_TIME=$(date +%s.%N)
TASK_EXIT_CODE=0
if [ -n "${SERVICE_NAME}" ]; then
    { time task test:e2e -- "${SERVICE_NAME}"; } 2>&1 | tee "${TASK_RESULT}" || TASK_EXIT_CODE=$?
else
    { time task test:e2e; } 2>&1 | tee "${TASK_RESULT}" || TASK_EXIT_CODE=$?
fi
END_TIME=$(date +%s.%N)
TASK_DURATION=$(echo "${END_TIME} - ${START_TIME}" | bc -l)
echo "-------------------------------------------------"
if [ ${TASK_EXIT_CODE} -eq 0 ]; then
    echo -e "✅ ${COLOR_BLUE}'task test:e2e' finished successfully. Duration: ${TASK_DURATION}s${COLOR_NC}"
else
    echo -e "⚠️  ${COLOR_YELLOW}'task test:e2e' finished with errors (exit code: ${TASK_EXIT_CODE}). Duration: ${TASK_DURATION}s${COLOR_NC}"
fi
echo ""

# --- 結果の比較と表示 ---
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  📊 ${COLOR_GREEN}実行時間比較結果${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
printf "%-15s: %10.3fs\n" "Make" "${MAKE_DURATION}"
printf "%-15s: %10.3fs\n" "Task" "${TASK_DURATION}"
echo "-------------------------------------------------"

if (( $(echo "${MAKE_DURATION} < ${TASK_DURATION}" | bc -l) )); then
    DIFF=$(echo "${TASK_DURATION} - ${MAKE_DURATION}" | bc -l)
    PERCENT=$(echo "scale=1; ${DIFF} / ${MAKE_DURATION} * 100" | bc -l)
    echo -e "🏆 ${COLOR_GREEN}Make が ${DIFF}s (${PERCENT}%) 速い${COLOR_NC}"
elif (( $(echo "${TASK_DURATION} < ${MAKE_DURATION}" | bc -l) )); then
    DIFF=$(echo "${MAKE_DURATION} - ${TASK_DURATION}" | bc -l)
    PERCENT=$(echo "scale=1; ${DIFF} / ${TASK_DURATION} * 100" | bc -l)
    echo -e "🏆 ${COLOR_GREEN}Task が ${DIFF}s (${PERCENT}%) 速い${COLOR_NC}"
else
    echo -e "🤝 ${COLOR_YELLOW}実行時間はほぼ同じです${COLOR_NC}"
fi

echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  ${COLOR_GREEN}E2Eテスト速度比較が完了しました。${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"

echo ""
echo -e "${COLOR_YELLOW}使用方法:${COLOR_NC}"
echo -e "  ./e2e_speed_check.sh           # 全サービスのテストを実行"
echo -e "  ./e2e_speed_check.sh home      # homeサービスのみテストを実行"
echo -e "  ./e2e_speed_check.sh user      # userサービスのみテストを実行"
