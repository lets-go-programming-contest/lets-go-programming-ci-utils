#!/bin/bash

dir=$(dirname "$0")
. "$dir/common.sh"

students=$(get_students "$HEAD" )

if test $(echo "$students" | wc -l) -ne 1; then
  printf "Commits contains changes in next students labs:\n"
  echo "$students"
  printf "But commits must contains changes only in one student labs".
fi
