#!/bin/bash
while true
do
    for (( c=1; c<$LT_ITER; c++ ))
    do
        bagop -b -ttl=$TTL
        bagop -c
        sleep $SLEEP
    done
    # Run Long-term backup once every $LT_ITER if exported
    if [ -n "$LT_ITER" ]; then
        bagop -b -ttl=$LT_TTL
        bagop -c
        sleep $SLEEP
    fi
done