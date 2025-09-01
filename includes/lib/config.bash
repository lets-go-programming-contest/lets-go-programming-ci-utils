#!/bin/bash

SKIP_CONFIG_VALUE="skip"

get_cfg_value() {
    local file="$1"
    local key="$2"

    local value="${SKIP_CONFIG_VALUE}"
    if test -f "${file}" ; then
      local value=$(grep -E "${key}:" "${file}" | awk '{print $2}')
    fi

    echo "$value"
}