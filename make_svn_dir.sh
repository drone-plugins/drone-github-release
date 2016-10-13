#!/bin/bash

user=$1
pass=$2
url=$3
baseDir=$4

svn co --username $user --password "$pass" --depth empty --trust-server-cert --non-interactive "$url" $baseDir
