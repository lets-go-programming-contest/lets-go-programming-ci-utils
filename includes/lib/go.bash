#!/bin/bash

go_check_mod() {
  local dir=$1

  if [[ ! -f "${dir}/go.mod" ]]; then
    printf "The directory is not a Go module.\n" >&2
    return 1
  fi

  return 0
}

go_tidy() {
  local dir=$1

  go -C "$dir" mod tidy
  local go_mod_tidy_status_code=$?
  return ${go_mod_tidy_status_code}
}

go_build() {
  local source_dir=$1
  local out_dir=$2

  go_check_mod "${source_dir}"
  local go_check_mod_exit_code=$?
  if [[ $go_check_mod_exit_code -ne 0 ]]; then
    return ${go_check_mod_exit_code}
  fi

  go_tidy "${source_dir}"
  local go_tidy_exit_code=$?
  if [[ $go_tidy_exit_code -ne 0 ]]; then
    return ${go_tidy_exit_code}
  fi

  local cmds=$(find "${source_dir}/cmd" -maxdepth 1 -mindepth 1 -type d)

  printf "Next targets for build has been found:\n"
  for cmd in ${cmds}; do
    printf -- "\t- %s\n" $cmd
  done

  local cmds_build_exit_code=0
  for cmd in $cmds; do
    local cmd_name=$(basename "${cmd}")
    printf "Build service ${cmd_name}.\n"

    go -C "${source_dir}" build -o "${out_dir}/${cmd_name}" "./cmd/${cmd_name}"
    local go_build_exit_code=$?

    if [[ $go_build_exit_code -eq 0 ]]; then
      printf "Target ${cmd_name} build finished successfully!\n"
    else
      printf "Target ${cmd_name} build failed!\n" >&2
      cmds_build_exit_code=${go_build_exit_code}
    fi
  done

  return ${cmds_build_exit_code}
}

go_lint() {
    local config=$1
    local source_dir=$2

    go_check_mod "${source_dir}"
    local go_check_mod_exit_code=$?
    if [[ $go_check_mod_exit_code -ne 0 ]]; then
      return ${go_check_mod_exit_code}
    fi

    go_tidy "${source_dir}"
    local go_tidy_exit_code=$?
    if [[ $go_tidy_exit_code -ne 0 ]]; then
      return ${go_tidy_exit_code}
    fi

    pushd "${source_dir}" > /dev/null
    golangci-lint run --config "${config}" ./...
    local lint_exit_code=$?
    popd > /dev/null

    return ${lint_exit_code}
}

go_test() {
    local source_dir=$1

    go_check_mod "${source_dir}"
    local go_check_mod_exit_code=$?
    if [[ $go_check_mod_exit_code -ne 0 ]]; then
      return ${go_check_mod_exit_code}
    fi

    go_tidy "${source_dir}"
    local go_tidy_exit_code=$?
    if [[ $go_tidy_exit_code -ne 0 ]]; then
      return ${go_tidy_exit_code}
    fi

    go -C "${source_dir}" test -v -cover ./...
    local go_test_exit_code=$?

    return ${go_test_exit_code}
}