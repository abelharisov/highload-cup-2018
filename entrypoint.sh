#!/bin/sh

mkdir /json

mongod > /dev/null &

eval "go $1"
