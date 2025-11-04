package temporal

import (
	"go.temporal.io/sdk/client"
)

var (
	DefaultClient client.Client
	DefaultHost   = "localhost:7233"
)

func Init() error {
	var err error
	DefaultClient, err = client.Dial(client.Options{HostPort: DefaultHost})
	return err
}
func Close() { if DefaultClient != nil { DefaultClient.Close() } }