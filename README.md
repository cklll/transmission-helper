# Transmission Helper
A tool that automate my routine operations with [transmission](https://github.com/transmission/transmission).
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

### Lint & Code format
We use the builtin gofmt to format our code.
You can auto-format the code by running
```
go fmt .
```

### TODO
* Build a binary executable. Can be used along with systemd or crontab
* (v1 config) All config will be done by environment variables
* (v2 config) Allow user to configure config. See [Approach #1](https://stackoverflow.com/a/35419545)
* Support remote host
* Add Linter
* Add back important tests (marked as TODO)
