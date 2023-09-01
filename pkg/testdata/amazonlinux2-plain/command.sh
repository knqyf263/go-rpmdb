#!/usr/bin/env bash
set -e

# rpm -qa --queryformat '{\"Epoch\":\"%{EPOCH}\",\"Name\":\"%{NAME}\",\"Version\":\"%{VERSION}\",\"Release\":\"%{RELEASE}\",\"Arch\":\"%{ARCH}\",\"SourceRpm\":\"%{SOURCERPM}\",\"Size\":%{SIZE},\"License\":\"%{LICENSE}\",\"Vendor\":\"%{VENDOR}\",\"Summary\":\"%{SUMMARY}\",\"SigMD5\":\"%{SIGMD5}\"\}\n'
rpm -qa --queryformat '%|EPOCH?{%{EPOCH}}:{}|¶%{NAME}¶%{VERSION}¶%{RELEASE}¶%|ARCH?{%{ARCH}}:{}|¶%|SOURCERPM?{%{SOURCERPM}}:{}|¶%{SIZE}¶%{LICENSE}¶%|VENDOR?{%{VENDOR}}:{}|¶%{SUMMARY}¶%{SIGMD5}\n'
cp /var/lib/rpm/Packages /mnt/export/Packages
