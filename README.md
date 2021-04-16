# Transmission Helper
A tool that automate my routine operations with [transmission](https://github.com/transmission/transmission).
* ✅ Send an email notification to you when a download completes
* Remove seed when a download completes

## Build & Example Usage
### Build
```sh
go build
env GOOS=linux GOARCH=arm GOARM=5 go build # for my raspberry pi
# output a transmission-helper binary
```

### Run binary
```
TH_SMTP_USER= \
TH_SMTP_PASS= \
TH_SMTP_HOST= \
TH_SMTP_PORT= \
TH_SMTP_SENDER_EMAIL= \
TH_NOTIFY_EMAILS= \
TH_REMOTE_USERNAME= \
TH_REMOTE_PASSWORD= \
./transmission-helper
```

### Systemd
Or to use it as systemd service, refer the example service in `transmission-helper.service.example` and `transmission-helper.timer.example`

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

## TODO
* (v2 config) Allow user to configure config. See [Approach #1](https://stackoverflow.com/a/35419545)
* Support remote host
* Add back important tests (marked as TODO)
* Refactor transmission-remote codes to separate package
