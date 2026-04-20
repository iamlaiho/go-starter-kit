#!/usr/bin/env bash
# Hard-block risky shell commands. Exit code 2 = block (no override).
set -euo pipefail

INPUT=$(cat)
CMD=$(echo "$INPUT" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('tool_input',{}).get('command',''))" 2>/dev/null || true)

BLOCKED_PATTERNS=(
  "rm -rf"
  "git push --force"
  "git push -f"
  "git reset --hard"
  "git clean -f"
  "DROP TABLE"
  "DROP DATABASE"
  "chmod -R 777"
  "> /dev/sd"
  "mkfs"
)

for pattern in "${BLOCKED_PATTERNS[@]}"; do
  if echo "$CMD" | grep -qi "$pattern"; then
    echo "BLOCKED: '$pattern' is not allowed. Review before proceeding." >&2
    exit 2
  fi
done

exit 0
