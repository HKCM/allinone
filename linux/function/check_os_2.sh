#!/usr/bin/env bash
_ostype="$(uname -s)"

case "$_ostype" in
Linux)
    _ostype=unknown-linux-gnu
    ;;
Darwin)
    if [[ $_cputype = arm64 ]]; then
    _cputype=aarch64
    fi
    _ostype=apple-darwin
    ;;
*)
    err "machine architecture is currently unsupported"
    ;;
esac