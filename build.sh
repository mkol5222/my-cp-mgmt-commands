#!/bin/bash

COMMANDS="approve_session show_access_rulebase show_sessions show_packages login discard  install_policy  logout  publish  reject_session  submit_session  verify_policy"

mkdir -p /workspaces/my-cp-mgmt-commands/releases/windows-amd64
mkdir -p /workspaces/my-cp-mgmt-commands/releases/linux-amd64

cd /workspaces/my-cp-mgmt-commands/commands
for CMD in $COMMANDS; do
    echo $CMD
    pushd $CMD
    ls
    ls "./$CMD.go"
    go get
    
    GOARCH=amd64 GOOS=linux go build "./$CMD.go"
    mv "$CMD" /workspaces/my-cp-mgmt-commands/releases/linux-amd64

    GOARCH=amd64 GOOS=windows go build "./$CMD.go"
    mv "$CMD.exe" /workspaces/my-cp-mgmt-commands/releases/windows-amd64

    popd
done