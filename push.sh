#!/bin/bash

user=$1
pass=$2
baseDir=$3

cd $baseDir
svn add *
svn ci --trust-server-cert --non-interactive --username $user --password "$pass" -m "$HOSTNAME: $DRONE_REPO-$DRONE_COMMIT: $DRONE_BUILD_NUMBER" *'
