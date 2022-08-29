#!/bin/bash
# More safety, by turning some bugs into errors.
# Without `errexit` you don’t need ! and can replace
# PIPESTATUS with a simple $?, but I don’t do that.
set -o errexit -o pipefail -o noclobber -o nounset

echo "Entrypoint script for file-watcher"
echo "Container args: $@"

########################
# STANDARD ENVIRONMENT #
########################
: "${MODULE_NAME:?Need to set MODULE_NAME to string}"
: "${MODULE_TYPE:?Need to set MODULE_TYPE to string}"

# Validate the environment according to module type
if [[ "$MODULE_TYPE" == "INGRESS" ]]
then
    : "${EGRESS_URLS:?Need to set EGRESS_URLS to string}"
    echo "Ingress"
elif [[ "$MODULE_TYPE" == "EGRESS" ]]
then
    : "${EGRESS_URLS:?Need to set EGRESS_URLS to string}"
    : "${INGRESS_PORT:?Need to set INGRESS_PORT to string}"
    : "${INGRESS_PATH:?Need to set INGRESS_PATH to string}"
    echo "Egress"
elif [[ "$MODULE_TYPE" == "PROCESS" ]]
then
    : "${EGRESS_URLS:?Need to set EGRESS_URLS to string}"
    : "${INGRESS_PORT:?Need to set INGRESS_PORT to string}"
    : "${INGRESS_PATH:?Need to set INGRESS_PATH to string}"
    echo "Process"
elif [[ "$MODULE_TYPE" == "FANOUT" ]]
then
    echo "Fanout NOT SUPPORTED"
    exit 1
else
    echo "Unrecognized MODULE_TYPE = $MODULE_TYPE, choose from INGRESS, EGRESS, PROCESS"
    exit 1
fi
echo Environment validated

echo Egress urls: $EGRESS_URLS

#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
# YOUR CODE HERE  #
#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓

# CALL THE MAIN SCRIPT
./app $@