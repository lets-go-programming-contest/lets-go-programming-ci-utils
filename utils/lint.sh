#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

add_in_total() {
  if [ -s "logs/lint-$1-error-log.txt" ]; then
    echo "$1...FAILED" >> total_lint.txt
  else
    echo "$1...OK" >> total_lint.txt
  fi
}

student=$(get_students "${HEAD}" | head -n 1)
tasks=$(get_tasks "${HEAD}" "${student}")

if [ -z "$tasks" ]; then
  echo "No solutions provided, skip!" >&2
  exit 0
fi

for task in $tasks; do
  mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "LINT_MODE")
  case $mode in
    none)
      echo "Mode is set to NONE, skip" > "logs/lint-$task-log.txt"
      add_in_total "$task"
    ;;
    *)
      if [ ! -f "$student/$task/go.mod" ]; then
        echo "The solution directory is not a Go module" > "logs/lint-$task-error-log.txt"
        add_in_total "$task"
        continue
      fi
      pushd "$student/$task" > /dev/null
      golangci-lint run --config "$WORKDIR/.golangci.yml" ./... > "$WORKDIR/logs/lint-$task-error-log.txt"
      popd > /dev/null
      add_in_total "$task"
    ;;
  esac
done

cat total_lint.txt

if grep -q "FAILED" total_lint.txt; then
  exit 1
fi
exit 0
