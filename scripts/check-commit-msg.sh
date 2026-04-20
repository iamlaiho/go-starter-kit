#!/usr/bin/env bash
set -euo pipefail

MSG_FILE="$1"
MSG=$(cat "$MSG_FILE")

PATTERN="^(feat|fix|docs|style|refactor|test|chore|ci|perf|build|revert)(\(.+\))?: .{1,72}"

if ! echo "$MSG" | grep -qE "$PATTERN"; then
  echo "ERROR: Commit message does not follow Conventional Commits format."
  echo ""
  echo "Expected: <type>(<scope>): <description>"
  echo "Types: feat, fix, docs, style, refactor, test, chore, ci, perf, build, revert"
  echo ""
  echo "Got: $MSG"
  exit 1
fi
