#!/bin/sh

touch /var/log/mongo.log
tail -f /var/log/mongo.log &

mongod --fork --logpath /var/log/mongo.log

eval "go $1"
