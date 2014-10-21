#!/bin/bash

export GOBIN="${HOME}/tmp/mygo"

for f in $(ls *.go);  do
    go install ${f}
done
