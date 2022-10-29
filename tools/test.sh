#!/bin/bash

[ -n "$DEBUG" ] && set -x
set -eo pipefail

TMP_HTTP=$(mktemp --tmpdir HTTP.XXXXXXXX)
TMP_BODY=$(mktemp --tmpdir BODY.XXXXXXXX)
TMP_BODYGRON=$(mktemp --tmpdir BODYGRON.XXXXXXXX)

./smtp2http -listen :5353 &
PID_smtp2http=$!
trap "kill $PID_smtp2http" EXIT

timeout 1 nc -l -p 8080 > $TMP_HTTP &
PID_http=$!

cat <<EOS | curl -s 'smtp://localhost:5353' \
    --mail-from 'mail-from@localhost.localdomain' \
    --mail-rcpt 'mail-rcpt@localhost.localdomain' \
    --upload-file /dev/stdin
From: header-from@localhost.localdomain
To: header-to@localhost.localdomain
Cc: header-cc@localhost.localdomain
Bcc: header-bcc@localhost.localdomain
Subject: subject

body text
EOS

wait $PID_http || true

echo "::HTTP::"
cat "$TMP_HTTP"
echo
echo

cat "$TMP_HTTP" | sed -n -e '/^\r\?$/,$p' > "$TMP_BODY"

echo "::Body::"
cat "$TMP_BODY"
echo
echo

cat "$TMP_BODY" | jq . | gron | cat > "$TMP_BODYGRON"

NEEDLES=(
    'json.addresses.bcc = "header-bcc@localhost.localdomain";'
    'json.addresses.cc = "header-cc@localhost.localdomain";'
    'json.addresses.from = "header-from@localhost.localdomain";'
    'json.addresses.mailfrom = "mail-from@localhost.localdomain";'
    'json.addresses.to = "header-to@localhost.localdomain";'
    'json.body.text = "body text\n\r\n";'
    'json.subject = "subject";'
)

for needle in "${NEEDLES[@]}"
do
    fgrep "$needle" "$TMP_BODYGRON"
done

echo
echo AC
