# go-rotate-backups

> Still a WIP

go-rotate-backups backups up files to a target and rotates them. It currently only has support for `local` and `s3`.

```shell
$ tree -L 1 backups 
backups
├── daily
├── monthly
├── weekly
└── yearly
```

## Installation

```shell
# homebrew
brew install stenic/tap/go-rotate-backups

# gofish
gofish rig add https://github.com/stenic/fish-food
gofish install github.com/stenic/fish-food/go-rotate-backups

# scoop
scoop bucket add go-rotate-backups https://github.com/stenic/scoop-bucket.git
scoop install go-rotate-backups

# go
go install github.com/stenic/go-rotate-backups@latest

# docker 
docker pull ghcr.io/stenic/go-rotate-backups:latest

# dockerfile
COPY --from=ghcr.io/stenic/go-rotate-backups:latest /go-rotate-backups /usr/local/bin/
```

> For even more options, check the [releases page](https://github.com/stenic/go-rotate-backups/releases).


## Run

```shell
# Installed
go-rotate-backups -h

# Docker
docker run -ti ghcr.io/stenic/go-rotate-backups:latest -h

# Kubernetes
kubectl run go-rotate-backups --image=ghcr.io/stenic/go-rotate-backups:latest --restart=Never -ti --rm -- -h
```

## Documentation

```shell
go-rotate-backups -h

Usage:
  go-rotate-backups [flags] files...

Flags:
      --daily int          Amount of daily backups to keep (default 7)
      --driver string      Driver selection (file, s3) (default "file")
  -h, --help               help for go-rotate-backups
      --monthly int        Amount of monthly backups to keep (default 12)
      --target string      Base location where backup live (default "./backups")
  -v, --verbosity string   Log level (debug, info, warn, error, fatal, panic (default "info")
      --version            version for go-rotate-backups
      --weekly int         Amount of weekly backups to keep (default 4)
      --yearly int         Amount of yearly backups to keep (default 5)
```

### How does it work

When running the command, the provided files will be copied to the backup location. After the copy is
completed, it will check the rotation configuration to cleanup unneeded backups.

__Backup__

Depending on the current date, it will upload to different folders:


| Date | Target |
| --- | --- |
| first day of the year | yearly/${date}_${time} |
| first day of the month | monthly/${date}_${time} |
| first day of the week | weekly/${date}_${time} |
| other | daily/${date}_${time} |

__Rotate__

Rotate will keep the `n` most recent files in the backup folder and clear out the others.

See `yearly`, `monthly`, `weekly` and `daily` for setting the different rotation settings.


### Drivers

__local__

The local driver uses a local folder to store it's data. It does not require any special configuration.


__s3__

The s3 driver uses [Amazon S3](https://aws.amazon.com/s3/) to store it's data. To configure the s3
credentials you will need to set the following environment variables:

| Name | Description |
| --- | --- |
| `GRB_S3_BUCKET` | The bucket name |
| `AWS_*`  | AWS credentials to use. This can be either `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` or `AWS_PROFILE` if you want to use the shared credential file |


## Badges

[![Release](https://img.shields.io/github/release/stenic/go-rotate-backups.svg?style=for-the-badge)](https://github.com/stenic/go-rotate-backups/releases/latest)
[![Software License](https://img.shields.io/github/license/stenic/go-rotate-backups?style=for-the-badge)](./LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/stenic/go-rotate-backups/Release?style=for-the-badge)](https://github.com/stenic/go-rotate-backups/actions?workflow=build)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## License

[License](./LICENSE)
