set -x
echo "$CRON /app/bagop" > ./bagop-cron

# Copy bagop-cron file to the cron.d directory
mv ./bagop-cron /etc/cron.d/bagop-cron

# Give execution rights on the cron job
chmod 0644 /etc/cron.d/bagop-cron

# Apply cron job
crontab /etc/cron.d/bagop-cron

# Create the log file to be able to run tail
touch /var/log/cron.log

cron && tail -f /var/log/cron.log