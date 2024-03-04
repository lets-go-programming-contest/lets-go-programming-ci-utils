#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

add_in_total() {
  if [ $? -ne 0 ] || [ -s "logs/unit-tests-$1-error-log.txt" ]; then
    echo "$1...FAILED" >> total_unit-tests.txt
  else
    echo "$1...OK" >> total_unit-tests.txt
  fi
}

student=$(get_students "${HEAD}" | head -n 1)
tasks=$(get_tasks "${HEAD}" "${student}")

if [ -z "$tasks" ]; then
  echo "No solutions provided, add the solutions to the task directories to continue!" >&2
  exit 0
fi

for task in $tasks; do
  mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "UNIT_TESTS_MODE")
  case $mode in
    none)
      echo "Mode is set to NONE, skip" > "logs/unit-tests-$task-log.txt" 
      add_in_total "$task"
    ;;
    *)
      if [ ! -f "$student/$task/go.mod" ]; then
        echo "The solution directory is not a Go module" > "logs/unit-tests-$task-error-log.txt"
        add_in_total "$task"
        continue
      fi

      if [ -d "$TEST_DIR_COMMON/$task/tests" ]; then
        mkdir -p "$student/$task/tests/ci"
        cp -r "$TEST_DIR_COMMON/$task/tests/"* "$student/$task/tests/ci"
        ls -la "$student/$task/tests/ci"
      fi
      go -C "$student/$task" mod tidy > "logs/unit-tests-$task-log.txt" 2> "logs/unit-tests-$task-error-log.txt"
      if [ -s "logs/unit-tests-$1-error-log.txt" ]; then
        echo "go mod"
        add_in_total "$task"
        continue
      fi
      export BUILD_BIN="${WORKDIR}/bin/${task}/"
      export STATIC_DIR="${TEST_DIR_COMMON}/{$task}"
      go -C "$student/$task" test -v -cover ./... > "logs/unit-tests-$task-trace.txt" 2> "logs/unit-tests-$task-error-log.txt"
      add_in_total "$task"
    ;;
  esac
done

cat total_unit-tests.txt

if grep -q "FAILED" total_unit-tests.txt; then
  exit 1
fi
exit 0

