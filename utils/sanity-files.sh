#!/bin/bash

dir=$(dirname "$0")
. "$dir/common.sh"

no_labs_files=$(get_diff "$HEAD" | grep -v -E "$LAB_FILES_REGEXP_PATTERN")
if test "$no_labs_files"; then
  printf "Affected CI/CD infrastructure files:\n"
  echo "$no_labs_files"
  exit 1
fi
