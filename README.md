## Transmission Helper
A tool that automate my routine operations with [transmission](https://github.com/transmission/transmission).

## Features
* ✅ Send an email notification to you when a download completes
* Remove seed when a download completes

## Development
```
docker-compose up -d
docker-compose exec app /bin/sh

go run .
```

### Run Test
```
go test
```

### TODO
* Build a binary executable. Can be used along with systemd or crontab
* (v1 config) All config will be done by environment variables
* (v2 config) [TBC] Allow user to configure config. Maybe Read from a config file (with hard-coded path). Or create a CLI to update config.
* Release in GitHub
* Support remote host
* Add Linter
* Add back important tests (marked as TODO)
