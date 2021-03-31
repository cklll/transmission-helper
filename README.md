## Transmission Helper
A tool that automate my routine operations with [transmission](https://github.com/transmission/transmission).

## Road Map
* Send an email notification to you when a download completes
* Remove seed when a download completes
* Build a binary executable. Can be used along with systemd or crontab
* (v1 config) All config will be done by environment variables
* (v2 config) [TBC] Allow user to configure config. Maybe Read from a config file (with hard-coded path). Or create a CLI to update config.
* Release in GitHub

## Development
```
docker-compose up -d
docker-compose exec app /bin/sh

go run main.go
```

### Run Test
```
go test
```
