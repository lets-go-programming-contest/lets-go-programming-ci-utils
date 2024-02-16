#!/bin/bash

base_branch=$(git rev-list --max-parents=0 --date-order --reverse "$BASE_BRANCH")
current_branch=$(git rev-list --max-parents=0 --date-order --reverse "$HEAD")

if test $(echo "$current_branch" | wc -l) -ne 1; then
  printf "The branch has multiple initial commits:\n"
  echo "$current_branch"
  printf "Fix your branch structure."
  exit 1
fi


if test $(echo "$base_branch" | wc -l) -ne 1; then
  printf "The base branch has multiple initial commits:\n"
  echo "$base_branch"
  printf "Fix your %s branch structure." "$BASE_BRANCH"
  exit 1
fi

base_branch_root=$(echo "$base_branch" | head -1)
current_branch_root=$(echo "$current_branch" | head -1)

if test "$base_branch_root" != "$current_branch_root"; then
  printf "Initial commits in fork and root not equal:\n"
  printf "Current branch init: %s\n" "$current_branch_root"
  printf "Base branch init: %s\n" "$base_branch_root"
  exit 1
fi

exit 0