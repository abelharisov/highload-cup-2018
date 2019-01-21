#!/bin/sh

mongod > /dev/null &

eval "go $1"
