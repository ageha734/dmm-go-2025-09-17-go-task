#!/bin/bash

set -euo pipefail

COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_NC='\033[0m'

if [ -z "${1-}" ]; then
  echo "エラー: 比較対象のタスク名を指定してください。" >&2
  echo "使用法: ./compare_speed.sh <task_name>"
  echo "例:     ./compare_speed.sh build"
  exit 1
fi

TASK_NAME=$1

# --- ヘッダー表示 ---
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  ⏱️  ${COLOR_GREEN}Makefile vs Taskfile 実行時間比較${COLOR_NC}"
echo -e "  対象タスク: ${COLOR_GREEN}${TASK_NAME}${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo ""

# --- Makefileの実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'make ${TASK_NAME}' ...${COLOR_NC}"
echo "-------------------------------------------------"
# timeコマンドでmakeの実行時間を計測
{ time make "${TASK_NAME}"; } 2>&1
echo "-------------------------------------------------"
echo -e "✅ ${COLOR_BLUE}'make ${TASK_NAME}' finished.${COLOR_NC}"
echo ""
echo ""

# --- Taskfileの実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'task ${TASK_NAME}' ...${COLOR_NC}"
echo "-------------------------------------------------"
# timeコマンドでtaskの実行時間を計測
{ time task "${TASK_NAME}"; } 2>&1
echo "-------------------------------------------------"
echo -e "✅ ${COLOR_BLUE}'task ${TASK_NAME}' finished.${COLOR_NC}"
echo ""

echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  ${COLOR_GREEN}比較が完了しました。${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
