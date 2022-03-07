# bagop

[![Go Report Card](https://goreportcard.com/badge/github.com/swexbe/bagop)](https://goreportcard.com/report/github.com/swexbe/bagop)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/swexbe/bagop)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/swexbe/bagop)

Tool to make automatic backups of any number of docker database containers to AWS Glacier

## Getting Started

Example run command:

```bash
docker run -e BAGOP_VAULT_NAME=myvaultname
    -e AWS_DEFAULT_REGION=myregion
    -e AWS_SECRET_ACCESS_KEY=mysecretaccesskey
    -e AWS_ACCESS_KEY_ID=myaccesskeyid
    -e CRON="0 4 * * 7"
    -e BAGOP_TTL=90
    -v /var/run/docker.sock:/var/run/docker.sock
    -d swexbe/bagop:latest
```

- NOTE: CRON must be a valid cron time/date field
- Set the labels `bagop.enable=true` and `bagop.name=dbname` for any dn containers you wish to automatically backup

## Configuration

How to configure bagop.

### Environment Variables for bagop Container

The application is configured through environment variables.

| Key                   | Required | Description                                                                                                           | Example     |
| --------------------- | -------- | --------------------------------------------------------------------------------------------------------------------- | ----------- |
| AWS_REGION            | yes      | The AWS region in which your vault is located                                                                         | us-east-1   |
| AWS_SECRET_ACCESS_KEY | yes      | Your AWS secret access key                                                                                            | secret      |
| AWS_ACCESS_KEY_ID     | yes      | Your AWS access key id                                                                                                | secret      |
| CRON                  | yes      | How often to run regular backups, accepts any valid CRON expression.                                                  | `0 4 * * *` |
| LT_CRON               | no       | How often to run long-term backups, accepts any valid CRON expression. If not set, no long-term backups will be made. | `0 4 * * 7` |
| BAGOP_VAULT_NAME      | yes      | The name of your AWS Glacier Vault                                                                                    | testvault   |
| BAGOP_TTL             | no       | Time to Live for regular backups (in days). If not set, backups will never expire.                                    | 90          |
| BAGOP_LT_TTL          | no       | Time to Live for long-term backups (in days) If not set, backups will never expire.                                   | 365         |
| BAGOP_VERBOSE         | no       | Run i verbose mode. Defaults to false.                                                                                | false       |
| BAGOP_COLOR           | no       | Run with color output. Defaults to true.                                                                              | true        |

With the above example values set, backups will be made every day at 04:00. A long-term backup will be made every sunday at 04:00. Regular backups will be deleted the next time bagop runs after 90 days. Long-term backups will be deleted after 365 days.

Additional environment variables for configuring the connection to AWS and Docker are also available. Documentation for these can be found in [Docker Go SDK docs](https://pkg.go.dev/github.com/docker/docker/client#NewEnvClient) and [AWS for Go docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

### Labels for Database Containers

To configure backups for an individual container, the following docker labels can be set:

| Key          | Required | Description                                                                                                                 |
| ------------ | -------- | --------------------------------------------------------------------------------------------------------------------------- |
| bagop.enable | yes      | Enable bagop for this container if set to `true`                                                                            |
| bagop.name   | no       | The name of the resulting .sql file for this container, will use docker id if not set                                       |
| bagop.vendor | no       | can be set to `postgres` or `mysql`. Overrides vendor detection and forces bagop to treat the container as the given vendor |

### Volumes

The following directories can be mounted on the filesystem:

| Directory    | Description                                                                                     |
| ------------ | ----------------------------------------------------------------------------------------------- |
| `/extra`     | Anything mounted in this directory will be added to each backup                                 |
| `/var/bagop` | Persistent data will be stored in this directory, i.e. Glacier archive IDs and their expiration |

### Input Parameters

Manual backups can be run using docker exec or interactive shell. The following input parameters are available:

| Flag         | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| -v           | Verbose mode                                                           |
| -b           | Make a backup                                                          |
| -c           | Delete expired containers from Glacier                                 |
| -l           | List all non-expired containers                                        |
| -ttl=        | Time to Live for backup in days. If not set, backup will never expire. |
| -version     | Display version                                                        |
| -force-color | Force color output                                                     |

## Retrieving An Archive

1. Run `bagop -l` inside the container and pick out an archive to retrieve. This can be done using `docker exec my_bagop_container bagop -l`.
2. Follow the [AWS Glacier CLI Documentation for retrieving an archive](https://docs.aws.amazon.com/amazonglacier/latest/dev/downloading-an-archive-two-steps.html) starting from step 1b.


## Docker DB for test command

```bash
docker run -l "bagop.enable=true" -l "bagop.name=testdb" -l "bagop.vendor=postgres" -e "POSTGRES_PASSWORD=password" postgres
```
## Contributions

Would be appreciated
