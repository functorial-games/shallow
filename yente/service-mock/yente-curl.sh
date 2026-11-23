#!/bin/sh
# A curl client, also usable inside a Shizuku/rish shell if that shell has curl.
# No Shizuku permission is required for ordinary HTTPS requests.
set -eu
if [ "$#" -ne 4 ]; then
    printf 'usage: %s HTTPS_BASE ROUTE EVENT_ID TEXT\n' "$0" >&2
    exit 2
fi
: "${YENTE_CA:?set YENTE_CA to the mock server certificate/CA path}"
: "${YENTE_MOCK_TOKEN:?set YENTE_MOCK_TOKEN for authenticated mock access}"
base=${1%/}
route=$2
id=$3
text=$4
case $route in *[!A-Za-z0-9_.-]*|'') echo 'bad route' >&2; exit 2;; esac
case $id in *[!A-Za-z0-9_.-]*|'') echo 'bad event id' >&2; exit 2;; esac
printf '%s' "$text" | curl --fail-with-body --silent --show-error \
    --cacert "$YENTE_CA" \
    --header "Authorization: Bearer $YENTE_MOCK_TOKEN" \
    --header "Content-Type: text/plain" \
    --header "X-Yente-Event-ID: $id" \
    --data-binary @- "$base/v1/cli/$route"
