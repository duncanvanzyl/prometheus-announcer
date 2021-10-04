package configstore

import "context"

type ConfigurationStore interface {
	AddTarget(string, []string, map[string]string) error
	RemoveTarget(string)
	JSON() ([]byte, error)
	Run(context.Context)
}
