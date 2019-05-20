#!/bin/bash
docker build -t test2 --build-arg=ci_sha="$(git rev-parse HEAD | tr -d "\n")" --build-arg=ci_description="$(git log -1 --pretty=%B | tr -d "\n")" --build-arg=ci_version="local-docker" .
