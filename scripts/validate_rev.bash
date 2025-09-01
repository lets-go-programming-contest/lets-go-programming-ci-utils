#!/bin/bash

if [[ -z "$1" ]]; then
  echo "Rev not set." >&2
  exit 1
fi

if [[ ! $(git rev-parse --verify --quiet "$1") ]]; then
  echo "Failed calculate rev. Rev $1 not found." >&2
  exit 1
fi

echo "$1"