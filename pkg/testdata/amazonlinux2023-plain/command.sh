#!/usr/bin/env bash
set -e

rpm -qa --queryformat '%|EPOCH?{%{EPOCH}}:{}|¶%{NAME}¶%{VERSION}¶%{RELEASE}¶%|ARCH?{%{ARCH}}:{}|¶%|SOURCERPM?{%{SOURCERPM}}:{}|¶%{SIZE}¶%{LICENSE}¶%|VENDOR?{%{VENDOR}}:{}|¶%{SUMMARY}¶%{SIGMD5}\n'
cp /var/lib/rpm/rpmdb.sqlite /mnt/export/rpmdb.sqlite
