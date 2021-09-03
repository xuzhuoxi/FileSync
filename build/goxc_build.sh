#!/usr/bin/env bash

CRTDIR=$(cd `dirname $0`; pwd)

goxc -os="linux darwin windows freebsd openbsd" -arch="amd64 arm" -n=FileSync -pv=1.0.0 -wd=${CRTDIR}/../src -d=${CRTDIR}/./release -include=*.go,README*,LICENSE*