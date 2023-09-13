#!/usr/bin/env bash
set -e

dnf -y install https://repo.mongodb.org/yum/redhat/8/mongodb-org/5.0/x86_64/RPMS/mongodb-org-shell-5.0.9-1.el8.x86_64.rpm 1>&2
dnf -y install https://repo.mongodb.org/yum/redhat/8/mongodb-org/5.0/x86_64/RPMS/mongodb-cli-1.25.0.x86_64.rpm 1>&2
rpm -qa --queryformat '%|EPOCH?{%{EPOCH}}:{}|¶%{NAME}¶%{VERSION}¶%{RELEASE}¶%|ARCH?{%{ARCH}}:{}|¶%|SOURCERPM?{%{SOURCERPM}}:{}|¶%{SIZE}¶%{LICENSE}¶%|VENDOR?{%{VENDOR}}:{}|¶%{SUMMARY}¶%{SIGMD5}\n'
cp /var/lib/rpm/rpmdb.sqlite /mnt/export/rpmdb.sqlite
