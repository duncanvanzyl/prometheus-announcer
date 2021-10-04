package configstore

import "github.com/hashicorp/go-hclog"

var logger = hclog.Default()

func SetLogger(l hclog.Logger) {
	logger = l.Named("cs")
}
