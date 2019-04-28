#!/usr/bin/env bash
docker run -v /root/ndp/dnssecure/bin/:/dnssecure/bin/ -v /root/ndp/dnssecure/dnssecure/src/:/dnssecure/src/ -e "CGO_ENABLED=0" dnssecure-builder
