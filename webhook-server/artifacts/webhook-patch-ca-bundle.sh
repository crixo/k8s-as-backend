#!/bin/bash

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
#set -o pipefail


#export CA_BUNDLE=$(kubectl config view --raw --flatten -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')
export CA_BUNDLE=$(kubectl config view --raw -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')

sed "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" webhook-registration.yaml.template > webhook-registration.yaml

# if command -v envsubst >/dev/null 2>&1; then
#     envsubst
# else
#     sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
# fi