#!/usr/bin/env bash
set -e

dnf module install nodejs -y 1>&2

rpm -qa --queryformat '%|EPOCH?{%{EPOCH}}:{}|¶%{NAME}¶%{VERSION}¶%{RELEASE}¶%|ARCH?{%{ARCH}}:{}|¶%|SOURCERPM?{%{SOURCERPM}}:{}|¶%{SIZE}¶%{LICENSE}¶%|VENDOR?{%{VENDOR}}:{}|¶%{SUMMARY}¶%{SIGMD5}¶%{MODULARITYLABEL}\n'
cp /var/lib/rpm/rpmdb.sqlite /mnt/export/rpmdb.sqlite
