#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"

MAINTAINERS_FILE="${TEST_DIR_COMMON}/MAINTAINERS"

no_lab_files=$(get_no_tasks_files "${TARGET_BRANCH}" "${HEAD}")
declare -a affected=()

for file in $no_lab_files; do
  [[ -z $file ]] && continue

  users=$(get_log "${TARGET_BRANCH}" "${HEAD}" "${file}" '%aE')
  for user in $users; do
    [[ -z "$user" ]] && continue

    if ! grep -q "^${user}$" "${MAINTAINERS_FILE}"; then
      affected+=("${file} affected by ${user}")
    fi
  done
done

if [[ ${#affected[@]} -gt 0 ]]; then 
  printf "The following root files are affected:\n" >&2
  printf -- "\t- %s\n" "${affected[@]}" >&2
  exit 1
fi

printf "Affected files list:\n"
printf -- "\t- %s\n" $(get_diff "${TARGET_BRANCH}" "${HEAD}")

printf "OK\n"
exit 0