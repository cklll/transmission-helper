# Transmission Helper
A tool that automate my routine operations with [transmission](https://github.com/transmission/transmission).
* ✅ Send an email notification to you when a download completes
* ✅ Remove seed when a download completes

## Warning
This project is still work in progress. I am using this project to learn Go so you may find some of the codes not following Go convention.

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

Note: I am planning to change the interface to
```
./transmission-helper --config-path=/path/to/config.json
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
* Include non-completed torrents
  * Ideally, we can set a frequent cron job (e.g. 10 minutes) to report completed torrents as soon as it completes.
  * But if none of the torrents completed, it also shouldn't report in-progress torrents every time, but after a certain interval.
  * Note: We can control that interval in the config file. Then perhaps write to a temp file (`~/.transmission-helper/last_notification_time.txt`, note: use `os.UserHomeDir()`) with last sent notification time.
* Support remote host
* Refactor transmission-remote codes to separate package
* Figure out if transmission-remote can output JSON
