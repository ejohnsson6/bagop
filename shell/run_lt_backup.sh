#!/bin/bash
echo "****RUNNING LONG-TERM BACKUP****"
bagop -b -force-color=${BAGOP_COLOR:-true} -v=${BAGOP_VERBOSE:-false} -ttl=$BAGOP_LT_TTL
echo "****RUNNING CLEAN****"
bagop -c -force-color=${BAGOP_COLOR:-true} -v=${BAGOP_VERBOSE:-false}