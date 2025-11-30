#!/bin/sh
set -euo

log() {
    level="$1"
    shift
    echo "$(date '+%Y-%m-%d %H:%M:%S') [$level] $*"
}

check_crond_up() {
    if ! pgrep crond >/dev/null 2>&1; then
        log ERROR "crond is not running"
        exit 1
    fi

    log INFO "crond is running"
}

log INFO "starting crond in the background"
crond -f -p -m off & # add -x pars,proc if we need to debug the container
check_crond_up

if [ "$#" -gt 0 ]; then
    exec "$@"
fi