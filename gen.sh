#!/bin/sh

if [ "$1" = "all" ]; then
    GEN_LUT=1 go generate ./gen
else
    go generate ./gen
fi
