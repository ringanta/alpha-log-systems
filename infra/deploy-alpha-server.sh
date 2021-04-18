#!/bin/sh

set -euo pipefail

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <AlphaServer Token> <DB Password>"
    exit 1
fi

popd ../
make compile
pushd

ansible-playbook -i hosts deploy-alpha-server.yml \
    -e "server_token=$1" \
    -e "db_password=$2"
