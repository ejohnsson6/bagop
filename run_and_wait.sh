#!/bin/sh
while true
do
    ./bagop -b
    ./bagop -c
    sleep $SLEEP
done