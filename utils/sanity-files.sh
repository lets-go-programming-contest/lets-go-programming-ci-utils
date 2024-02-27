#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

no_lab_files=$(get_diff "$HEAD" | grep -v -E "$LAB_FILES_REGEXP_PATTERN")

if test "${no_lab_files}"; then
  cat >&2 <<EOF
Affected CI/CD infrastructure files:
$(echo ${no_lab_files} | awk '{print "\t- " $1}')
Please, tidy up the files you're not allowed to touch before adding changes!
EOF
  exit 1
fi

lab_files=$(get_labs_files "$HEAD")
if test "${lab_files}"; then
  cat >&2 <<EOF
The following lab files are submitted for checking:
$(echo ${lab_files} | awk '{print "\t- " $1}')
EOF
  exit 0
else
  printf 'Warning: Lab files is not represented in this commit...'
fi