#!/bin/bash
# Needed to get ENV variables exported inside CRON
printenv | sed 's/^\(.*\)$/export \1/g' | grep -E "^export (BAGOP|AWS)" >/root/project_env.sh

# Create new CRONTAB based on $CRON and $LT_CRON, redirect output to stdout/stderr
echo "${CRON?Error CRON env not defined} . /root/project_env.sh; /home/root/run_backup.sh >> /var/bagop/log 2>&1" >/tmp/bagop-cron
if [ -n "$LT_CRON" ]; then
  echo "$LT_CRON . /root/project_env.sh; /home/root/run_lt_backup.sh >> /var/bagop/log 2>&1" >>/tmp/bagop-cron
fi

# Replace root crontab
cat /tmp/bagop-cron >/etc/crontabs/root

# Create logfile
touch /var/bagop/log
chmod 0666 /var/bagop/log

crond -b

# Tail logfile
tail -n 0 -F /var/bagop/log
