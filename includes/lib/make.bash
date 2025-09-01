#!/bin/bash

run_make_target() {
    local sources_dir="$1"
    local mk_file="$2"
    local target_name="$3"

    printf "Run make target ${target_name} from file ${mk_file}.\n"
    make -C ${sources_dir} --no-print-directory -f ${mk_file} ${target_name}
    make_status_code=$?
    return ${make_status_code}
}
