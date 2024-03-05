#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

add_in_total() {
  if [ -s "logs/custom-tests-$1-error-log.txt" ]; then
    echo "$1...FAILED" >> total_custom-tests.txt
  else
    echo "$1...OK" >> total_custom-tests.txt
  fi
}

student=$(get_students "${HEAD}" | head -n 1)
tasks=$(get_tasks "${HEAD}" "${student}")

if [ -z "$tasks" ]; then
  echo "No solutions provided, skip!" >&2
  exit 0
fi

for task in $tasks; do
  mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "CUSTOM_TESTS_MODE")
  case $mode in
    none)
      echo "Mode is set to NONE, skip" > "logs/custom-tests-$task-log.txt" 
      add_in_total "$task"
    ;;
    *)
      if [ -s "$TEST_DIR_COMMON/$task/Makefile" ] && grep -q "custom:" "$TEST_DIR_COMMON/$task/Makefile"; then
        make -C "$TEST_DIR_COMMON/$task" --no-print-directory \
              BIN="$WORKDIR/bin/$task" DIR="$WORKDIR/$student/$task" STATIC="${TEST_DIR_COMMON}/${task}"  \
              custom > "logs/custom-tests-$task-log.txt" 2> "logs/custom-tests-$task-error-log.txt"
        add_in_total "$task"
      else
        echo "Custom tests not found, skip" > "logs/custom-tests-$task-log.txt" 
        add_in_total "$task"
      fi
    ;;
  esac
done

cat total_custom-tests.txt

if grep -q "FAILED" total_custom-tests.txt; then
  exit 1
fi
exit 0

