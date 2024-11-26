cpe-insight is a tool to request information from your router from Telenor\
Downloads are available for Linux, FreeBSD and Windows

<img src="screenshots/status.png" alt="Requesting the status endpoint, and visualizing the JSON response with jq" width="50%">

## Installing
### Prebuilt binaries
Download and run the latest version in the [Releases](https://github.com/kivattt/cpe-insight/releases) page

Add it to your path environment variable, or (on Linux/FreeBSD) place the executable in `/usr/local/bin`

### Building from source
This requires Go 1.21.5 or above ([install Go](https://go.dev/dl/))
```
git clone https://github.com/kivattt/cpe-insight
cd cpe-insight
go build
./cpe-insight # cpe-insight.exe on Windows
```

## Usage
This tool outputs responses in JSON.\
You can use the awesome [jq](https://github.com/jqlang/jq/releases) to pretty-print them

Requesting a single endpoint:
```
cpe-insight --password="..." --endpoint status | jq
```

Requesting all endpoints and storing in a `file.json` JSON file:
```
cpe-insight --password="..." --all --output file.json
```
To look through all the responses, open the `file.json` in Firefox

It sends requests to `https://wifi.telenor.no`

## Known issues
Some requests might reply that they are starting a job, and you may need to re-request later to get the results of it

cpe-insight currently only supports requesting the GET endpoints of the CPE Insight API, but not these:
```
/${t}/telemetry-historic/${e}
/${t}/forgot-password/change/${e}/${r}
```
