#!/bin/bash
# Needed to get ENV variables exported inside CRON
printenv | sed 's/^\(.*\)$/export \1/g' | grep -E "^export (BAGOP|AWS)" >/root/project_env.sh

# Create new CRONTAB based on $CRON and $LT_CRON, redirect output to stdout/stderr
echo "${CRON?Error CRON env not defined} . /root/project_env.sh; /home/root/run_backup.sh > /tmp/stdout 2> /tmp/stderr" >/tmp/bagop-cron
if [ -n "$LT_CRON" ]; then
  echo "$LT_CRON . /root/project_env.sh; /home/root/run_lt_backup.sh > /tmp/stdout 2> /tmp/stderr" >>/tmp/bagop-cron
fi

# Replace root crontab
cat /tmp/bagop-cron >/etc/crontabs/root

# Create custom stdout and stderr named pipes
mkfifo /tmp/stdout /tmp/stderr
chmod 0666 /tmp/stdout /tmp/stderr

crond
# Have the main Docker process tail the files to produce stdout and stderr
# for the main process that Docker will actually show in docker logs.
tail -f /tmp/stdout &
tail -f /tmp/stderr >&2
