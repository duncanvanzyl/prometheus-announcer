# Prometheus Announcer

A gprc server that allows applications to announce themselves to prometheus 
using [prometheus http service discovery](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config).

## Announce
Connect using grpc and Announce the service to have it presented using the prometheus http service
discovery. The app to be discovered is responsible for reannouncing itself before the announcement 
lifetime has expired.

Example Announce Function:
```go
import (
	"github.com/google/uuid"
	pa "github.com/duncanvanzyl/prometheus-announcer"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)
⋮
grpcServer := "localhost:50051"
target := "192.168.1.1:9000"
id := uuid.New().String()
labels := map[string]string{"name": "value"}
⋮
// Note: It is probably worth doing something with the returned error value...
go pa.DialAndAnnounce(ctx, grpcServer, 1*time.Minute, id, target , labels)
```

## Server
There is a functional server in `cmd/server`. It is intended to be run in docker (or kubernetes). 
Create the image with `make docker-build`.

### Environment Variable
The following environment variables can be used:
| Variable     | Function                                                                              | Default        |
| ------------ | ------------------------------------------------------------------------------------- | -------------- |
| HSD_GPRCHOST | Port to allow clients to connect to.                                                  | :50051         |
| HSD_LOGLEVEL | Logging level.                                                                        | info           |
| HSD_HTTPHOST | The host and port for prometheus to use for service discovery.                        | localhost:8080 |
| HSD_LIFETIME | Announcement lifetimes. Announcement expires if not reannounced before this duration. | 2m             |
| HSD_INTERVAL | The interval to check for expired announcements.                                      | 30s            |
| HSD_WITHREST | Enable the REST API.                                                                  | true           |

### REST API
A endpoint for adding targets over http is available at "http://<HSD_HTTPHOST>/v1/announce".  
The REST api is enabled by default, but can be disabled by setting `HST_WITHREST=false`.  

Announcements can be made by sending an appropriate POST request to the server.
```bash
curl -d '{"id":"asdf","targets":["10.13.200.249:2113"],"labels":{"name1":"value1"}}' -H "Content-Type: application/json" -X POST http://localhost:8080/v1/announce
```

## Client
There is a functional client in `cmd/client`. It is intended to be run from the command line, or via systemd.
Create the executable with `make client`.

### Flags
The following flags can be used:
- -id
	- Unique ID for this announcer.  
    If not provided a random uuid will be used.
	- [type] string
- -labels
	- Labels to tag metrics scraped from the supplied targets with.  
    Must be in the form: "name1:value1,name2:value2".  
    Label names may contain ASCII letters, numbers, as well as underscores and must  
    match the regex "[a-zA-Z_][a-zA-Z0-9_]*".  
    Label values may contain any Unicode characters.
	- [type] string

### Environment Variables
The following environment variables can be used:
- HTTPSD_INTERVAL
	- The interval between announcements.
  - [type]        Duration
  - [default]     1m
- HTTPSD_GRPCSERVER
  - The GRPC server to announce to.  
    In the form "host:port"
  - [type]        String
  - [default]     localhost:50051