#!/bin/bash
set -eo pipefail
shopt -s nullglob

# shellcheck disable=SC2145
echo "./bin/FullNode $@" > command.txt
exec "./bin/FullNode" "$@"
