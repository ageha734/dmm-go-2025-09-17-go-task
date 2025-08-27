#!/bin/bash

set -euo pipefail

COLOR_BLUE='\033[0;34m'
COLOR_GREEN='\033[0;32m'
COLOR_NC='\033[0m'

if [ -z "${1-}" ]; then
  echo "エラー: 比較対象のタスク名を指定してください。" >&2
  echo "使用法: ./compare_speed_check.sh <task_type>"
  echo "例:     ./compare_speed_check.sh unit"
  echo "例:     ./compare_speed_check.sh e2e"
  echo "例:     ./compare_speed_check.sh build"
  exit 1
fi

TASK_TYPE=$1

# タスクタイプに応じてMakefileとTaskfileのターゲット名をマッピング
case "$TASK_TYPE" in
  "unit")
    MAKE_TARGET="test-unit"
    TASK_TARGET="test:unit"
    ;;
  "e2e")
    MAKE_TARGET="test-e2e"
    TASK_TARGET="test:e2e"
    ;;
  "build")
    MAKE_TARGET="build"
    TASK_TARGET="build"
    ;;
  "lint")
    MAKE_TARGET="lint"
    TASK_TARGET="lint"
    ;;
  "check")
    MAKE_TARGET="check"
    TASK_TARGET="check"
    ;;
  *)
    echo "エラー: サポートされていないタスクタイプ: $TASK_TYPE" >&2
    echo "サポートされているタスクタイプ: unit, e2e, build, lint, check"
    exit 1
    ;;
esac

# --- ヘッダー表示 ---
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  ⏱️  ${COLOR_GREEN}Makefile vs Taskfile 実行時間比較${COLOR_NC}"
echo -e "  対象タスク: ${COLOR_GREEN}${TASK_TYPE}${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo ""

# --- Makefileの実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'make ${MAKE_TARGET}' ...${COLOR_NC}"
echo "-------------------------------------------------"
# timeコマンドでmakeの実行時間を計測
{ time make "${MAKE_TARGET}"; } 2>&1
echo "-------------------------------------------------"
echo -e "✅ ${COLOR_BLUE}'make ${MAKE_TARGET}' finished.${COLOR_NC}"
echo ""
echo ""

# --- Taskfileの実行と時間計測 ---
echo -e "▶️  ${COLOR_BLUE}Executing 'task ${TASK_TARGET}' ...${COLOR_NC}"
echo "-------------------------------------------------"
# timeコマンドでtaskの実行時間を計測
{ time task "${TASK_TARGET}"; } 2>&1
echo "-------------------------------------------------"
echo -e "✅ ${COLOR_BLUE}'task ${TASK_TARGET}' finished.${COLOR_NC}"
echo ""

echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
echo -e "  ${COLOR_GREEN}比較が完了しました。${COLOR_NC}"
echo -e "${COLOR_BLUE}=================================================${COLOR_NC}"
