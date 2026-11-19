#!/bin/sh
set -eu

# Keep substituted values valid in both JavaScript and nginx configuration.
# Credentials, quotes, fragments and query strings do not belong in API base URLs.
case "$API_URL$BACKEND_URL" in
    *[!A-Za-z0-9_./~%:-]*) echo 'API URLs contain unsupported characters' >&2; exit 1 ;;
esac
case "$API_URL" in
    //*) echo 'API_URL must use an explicit HTTP(S) scheme or a local path' >&2; exit 1 ;;
esac
if ! printf '%s\n' "$API_URL" | grep -Eq '^(/[A-Za-z0-9_./~%:-]*|https?://[A-Za-z0-9_.:-]+(/[A-Za-z0-9_./~%:-]*)?)$'; then
    echo 'API_URL must be an HTTP(S) URL or an absolute path without query or fragment' >&2
    exit 1
fi
if ! printf '%s\n' "$BACKEND_URL" | grep -Eq '^https?://[A-Za-z0-9_.:-]+$'; then
    echo 'BACKEND_URL must be an HTTP(S) origin without path, query or fragment' >&2
    exit 1
fi
printf 'window.__WASA_CONFIG__ = {"apiURL":"%s"};\n' "$API_URL" > /tmp/wasa-config.js
envsubst '${BACKEND_URL}' < /etc/nginx/wasa.conf.template > /tmp/wasa-nginx.conf
exec nginx -c /tmp/wasa-nginx.conf -g 'daemon off;'
