#!/bin/bash

[ -n "$DEBUG" ] && set -x
set -eo pipefail

./smtp2http -listen :5353 &
pid_smtp2http=$!
trap "kill $pid_smtp2http" EXIT

timeout 1 nc -l -p 8080 > http &
pid_http=$!

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

wait $pid_http || true

echo "::HTTP::"
cat http
echo
echo

cat http | sed -n -e '/^\r\?$/,$p' > body

echo "::Body::"
cat body
echo
echo

cat body | jq . | gron | cat > body-gron

fgrep 'json.addresses.bcc = "header-bcc@localhost.localdomain";' body-gron
fgrep 'json.addresses.cc = "header-cc@localhost.localdomain";' body-gron
fgrep 'json.addresses.from = "header-from@localhost.localdomain";' body-gron
fgrep 'json.addresses.mailfrom = "mail-from@localhost.localdomain";' body-gron
fgrep 'json.addresses.to = "header-to@localhost.localdomain";' body-gron
fgrep 'json.body.text = "body text\n\n\r\n.\r";' body-gron
fgrep 'json.subject = "subject";' body-gron

echo
echo AC

