#!/bin/bash

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
#set -o pipefail

# same cert is mounted in each pod /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
# CERT=$(cat ca.crt  | base64 | tr -d '\n')
#export CA_BUNDLE=$(kubectl config view --raw --flatten -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')
export CA_BUNDLE=$(kubectl config view --raw -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')

echo $CA_BUNDLE
#sed "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" webhook-registration.yaml.template > webhook-registration.yaml

# if command -v envsubst >/dev/null 2>&1; then
#     envsubst
# else
#     sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
# fi