#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/config.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"
source "${TEST_DIR_UTILS}/includes/lib/go.bash"
source "${TEST_DIR_UTILS}/includes/lib/make.bash"

student=$(get_students "${TARGET_BRANCH}" "${HEAD}")
if [[ -z $student ]]; then
  printf "No target student, lint failed.\n" >&2
  exit 1
fi

task=$(get_tasks "${TARGET_BRANCH}" "${HEAD}" "${student}" | head -n 1)
if [[ -z $task ]]; then
  printf "No target task, lint failed.\n" >&2
  exit 1
fi

export BIN="${WORKDIR}/bin/${task}"
export DIR="${WORKDIR}/${student}/${task}"
export STATIC="${TEST_DIR_COMMON}/${task}"

mode=$(get_cfg_value "${TEST_DIR_COMMON}/${task}/ci.cfg" "lint")

case $mode in
  skip)
    printf "Lint stage disabled for current task.\tSkip.\n"
    exit 0
    ;;
  common)
    printf "Using 'lint' target from common task files.\n"

run_make_target "${WORKDIR}/${student}/${task}" "${TEST_DIR_COMMON}/${task}/Makefile" "lint"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  student)
    printf "Using 'lint' target from student task files.\n"
    
    run_make_target "${WORKDIR}/${student}/${task}" "Makefile" "lint"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  default)
    printf "Using default golangci-lint for current task files.\n"

    common_lint_config="${TEST_DIR_COMMON}/${task}/.golangci.yaml"
    if [[ -f $common_lint_config  ]]; then 
        printf "Run lint with config from common task files.\n"
        go_lint "${common_lint_config}" "${WORKDIR}/${student}/${task}" 
        common_go_lint_exit_code=$?
        if [[ $common_go_lint_exit_code -ne 0 ]]; then
            exit $common_go_lint_exit_code
        fi
    else
        printf "Config file in common task files not found.\tSkip.\n"
    fi

    student_lint_config="${WORKDIR}/${student}/${task}/.golangci.yaml"
    if [[ -f $student_lint_config ]]; then
        printf "Run lint with config from student dir.\n"
        go_lint "${student_lint_config}" "${WORKDIR}/${student}/${task}" 
        student_go_lint_exit_code=$?
        if [[ $student_go_lint_exit_code -ne 0 ]]; then
            exit $student_go_lint_exit_code
        fi
    else
        printf "Config file in student task files not found.\tSkip.\n"
    fi
    ;;
  *)
    printf "Unknown execute mode, lint failed.\n" >&2
    printf "Contact with admins for update execute mode for task ${task}.\n" >&2
    exit 1
    ;;
esac