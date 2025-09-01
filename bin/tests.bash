#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/config.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"
source "${TEST_DIR_UTILS}/includes/lib/go.bash"
source "${TEST_DIR_UTILS}/includes/lib/make.bash"

copy_test_files() {
  local student=$1
  local task=$2

  local source_dir="${TEST_DIR_COMMON}/${task}/tests"
  local target_dir="${WORKDIR}/${student}/${task}/tests/ci"

  printf "${source_dir}\n${target_dir}\n"

  if [[ -d "${source_dir}" ]]; then
    printf "Inject common tests.\n"
    mkdir -p "${target_dir}"
    cp -r "${source_dir}/"* "${target_dir}"
  fi 
}

student=$(get_students "${TARGET_BRANCH}" "${HEAD}")
if [[ -z $student ]]; then
  printf "No target student, tests failed.\n" >&2
  exit 1
fi

task=$(get_tasks "${TARGET_BRANCH}" "${HEAD}" "${student}" | head -n 1)
if [[ -z $task ]]; then
  printf "No target task, tests failed.\n" >&2
  exit 1
fi

export BIN="${WORKDIR}/bin/${task}"
export DIR="${WORKDIR}/${student}/${task}"
export STATIC="${TEST_DIR_COMMON}/${task}"

mode=$(get_cfg_value "${TEST_DIR_COMMON}/${task}/ci.cfg" "tests")

case $mode in
  skip)
    printf "Tests stage disabled for current task.\tSkip.\n"
    exit 0
    ;;
  common)
    printf "Using 'tests' target from common task files.\n"

    run_make_target "${WORKDIR}/${student}/${task}" "${TEST_DIR_COMMON}/${task}/Makefile" "tests"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  student)
    printf "Using 'tests' target from student task files.\n"
    
    run_make_target "${WORKDIR}/${student}/${task}" "${WORKDIR}/${student}/${task}/Makefile" "tests"
    run_make_target_exit_code=$?
    exit ${run_make_target_exit_code}
    ;;
  default)
    printf "Using default go test for current task files.\n"

    copy_test_files "${student}" "${task}"
    go_test "${WORKDIR}/${student}/${task}"
    go_build_exit_code=$?
    exit ${go_build_exit_code}
    ;;
  *)
    printf "Unknown execute mode, tests failed.\n" >&2
    printf "Contact with admins for update execute mode for task ${task}.\n" >&2
    exit 1
    ;;
esac
