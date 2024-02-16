#!/bin/bash

dir=$(dirname "$0")
. "$dir/common.sh"

student=$(get_students "$HEAD" | head -n 1)
tasks=$(get_tasks "$HEAD" "$student")

for task in $tasks; do
    mode=$(get_cfg_value "$TEST_DIR_COMMON/$task/ci.cfg" "TEST_MODE")

    case "$mode" in
    "none" )
      printf "Build test for %s set as NONE, build will be skipped.\n" "$task"
    ;;
    "make" )
    printf "Run test %s with Makefile by test system.\n" "$task"
    make -f "$TEST_DIR_COMMON/$task/Makefile" test BIN="$WORKDIR/bin/$task" DIR="$workdir/$student/$task"
    ;;
    "student" )
    printf "Run test %s with Makefile by student.\n" "$task"
    make -f "$student/$task/Makefile" test BIN="$WORKDIR/bin/$task" DIR="$workdir/$student/$task"
    ;;
    * )
      printf "Run build %s process in default mode.\n" "$task"
      bin_path="$WORKDIR/bin/$task/service"
      if [ -d "$PATH_TO_COMMON/$task/tests" ] ; then
        mkdir -p "$student/$task/tests/ci"
        cp -Rf "$PATH_TO_COMMON/$task/tests/" "$student/$task/tests/ci/"
      fi
        cd "$student/$task" || exit 1
        go mod tidy
        go test ./... -v -cover -binPath="$bin_path"
        cd "$WORKDIR" || exit 1
    ;;
    esac
done
