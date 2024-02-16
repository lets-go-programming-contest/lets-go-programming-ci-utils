#!/bin/bash

dir=$(dirname "$0")
. "$dir/common.sh"

student=$(get_students "$HEAD" | head -n 1)
tasks=$(get_tasks "$HEAD" "$student")

for task in $tasks; do
  mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "BUILD_MODE")

  case "$mode" in
    "none" )
      printf "Build mode for %s set as NONE, build will be skipped.\n" "$task"
    ;;
  "make" )
      printf "Run build %s with Makefile by test system.\n" "$task"
      make -f "$TEST_DIR_COMMON/$task/Makefile" build BIN="$WORKDIR/bin/$task" DIR="$WORKDIR/$student/$task"
    ;;
  "student" )
      printf "Run build %s with Makefile by student.\n" "$task"
      make -f "$student/$task/Makefile" build BIN="$WORKDIR/bin/$task" DIR="$WORKDIR/$student/$task"
    ;;
  * )
    printf "Run build %s process in default mode.\n" "$task"
        find "$student/$task/cmd" -name 'main.go' | while read -r main_file; do
          service_dir=$(dirname "$main_file")
          service_name=$(basename "$service_dir")
          mkdir -p "$WORKDIR/bin/$task"
          go build -o "$WORKDIR/bin/$task/$service_name" "$main_file"
        done
    ;;
  esac
done



