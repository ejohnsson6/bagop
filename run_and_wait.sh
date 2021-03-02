#!/bin/sh
while true
do
    bagop -b -ttl=$TTL
    bagop -c
    sleep $SLEEP
done