#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/config.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"
source "${TEST_DIR_UTILS}/includes/lib/go.bash"
source "${TEST_DIR_UTILS}/includes/lib/make.bash"

student=$(get_students "${TARGET_BRANCH}" "${HEAD}")
if [[ -z $student ]]; then
  printf "No target student, build failed.\n" >&2
  exit 1
fi

task=$(get_tasks "${TARGET_BRANCH}" "${HEAD}" "${student}" | head -n 1)
if [[ -z $task ]]; then
  printf "No target task, build failed.\n" >&2
  exit 1
fi

export BIN="${WORKDIR}/bin/${task}"
export DIR="${WORKDIR}/${student}/${task}"
export STATIC="${TEST_DIR_COMMON}/${task}"

mode=$(get_cfg_value "${TEST_DIR_COMMON}/${task}/ci.cfg" "build")

case $mode in
  skip)
    printf "Build stage disabled for current task.\tSkip.\n"
    exit 0
    ;;
  common)
    printf "Using 'build' target from common task files.\n"

    run_make_target "${WORKDIR}/${student}/${task}" "${TEST_DIR_COMMON}/${task}/Makefile" "build"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  student)
    printf "Using 'build' target from student task files.\n"
    
    run_make_target "${WORKDIR}/${student}/${task}" "${WORKDIR}/${student}/${task}/Makefile" "build"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  default)
    printf "Using default go build for current task files.\n"

    go_build "${WORKDIR}/${student}/${task}" "${BIN}" "${LOG}"
    go_build_exit_code=$?
    exit ${go_build_exit_code}
    ;;
  *)
    printf "Unknown execute mode, build failed.\n" >&2
    printf "Contact with admins for update execute mode for task ${task}.\n" >&2
    exit 1
    ;;
esac