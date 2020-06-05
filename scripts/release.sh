#!/usr/bin/env bash

set -eu
set -o pipefail

echo -n "Build timestamp: "
TIMESTAMP=$(date +%s)
echo $TIMESTAMP

echo -n "Git commit: "

  COMMIT=$(git rev-parse HEAD)
  RELEASE=true
  
echo $COMMIT

echo -n "Tagged version: "
  VERSION=v1.5.2
  RELEASE=true
  
echo $VERSION

echo Running "go $@"
exec go "$1" -ldflags \
	"-s -w -X storj.io/private/version.buildTimestamp=$TIMESTAMP
         -X storj.io/private/version.buildCommitHash=$COMMIT
         -X storj.io/private/version.buildVersion=$VERSION
         -X storj.io/private/version.buildRelease=$RELEASE" "${@:2}"