#!/usr/bin/env bash

INIT=${1:-systemd}

CONFIG_PATH="${XDG_CONFIG_HOME:-$HOME/.config}"
SETTINGS_PATH="${CONFIG_PATH}/photocatalog"

SYSTEMD_UNIT_PATH="${CONFIG_PATH}/systemd/user/"

if "${INIT}" == "systemd"
then
    cp ./init/systemd/photocatalog.service $SYSTEMD_UNIT_PATH/photocatalog.service
    if test ! -f "${SETTINGS_PATH}"
    then
        echo "TARGET=<specify target dir>\nMONITOR=<specify dir to monitor>\nMODE=hardlink" > "${SETTINGS_PATH}"
        ${EDITOR} "${SETTINGS_PATH}"
        exit $?
    else
        exit 0
    fi
fi

echo "Unknown init"
exit 2
