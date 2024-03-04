#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

student=$(get_students "${HEAD}" | head -n 1)
tasks=$(get_tasks "${HEAD}" "${student}")

if [ -z "$tasks" ]; then
  echo "No solutions provided, add the solutions to the task directories to continue!" >&2
  exit 0
fi

add_in_total() {
  if [ -s "logs/build-$1-error-log.txt" ]; then
    echo "$1...FAILED" >> total_build.txt
  else
    echo "$1...OK" >> total_build.txt
  fi
}

for task in $tasks; do
  mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "BUILD_MODE")

  case $mode in
    none)
      echo "Mode is set to NONE, skip" > "logs/build-$task-log.txt" 
      add_in_total "$task"
      ;;
    make)
      make -C "$TEST_DIR_COMMON/$task" --no-print-directory \
        BIN="$WORKDIR/bin/$task" DIR="$WORKDIR/$student/$task" STATIC="${TEST_DIR_COMMON}/${task}" \
        build > "logs/build-$task-log.txt" 2> "logs/build-$task-error-log.txt"
      add_in_total "$task"
      ;;
    student)
      make -C "$student/$task" --no-print-directory \
        BIN="$WORKDIR/bin/$task" DIR="$WORKDIR/$student/$task" STATIC="${TEST_DIR_COMMON}/${task}" \
        build > "logs/build-$task-log.txt" 2> "logs/build-$task-error-log.txt"
      add_in_total "$task"
      ;;
    *)
      if [ ! -f "$student/$task/go.mod" ]; then
        echo "The solution directory is not a Go module" > "logs/build-$task-error-log.txt"
        add_in_total "$task-$service_name"
        continue
      fi
      for main_file in $(find "$student/$task/cmd" -name "main.go"); do
        service_dir=$(dirname "$main_file")
        service_name=$(basename "$service_dir")
        go build -o "$WORKDIR/bin/$task/$service_name" "$main_file" > "logs/build-$task-$service_name-log.txt" 2> "logs/build-$task-$service_name-error-log.txt"

        add_in_total "$task-$service_name"
      done
      ;;
  esac
done

cat total_build.txt

if grep -q "FAILED" total_build.txt; then
  exit 1
fi
  exit 0

