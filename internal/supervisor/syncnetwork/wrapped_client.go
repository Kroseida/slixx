package syncnetwork

import (
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork"
)

type WrappedClient struct {
	Model  model.Satellite
	Client *syncnetwork.Client
}
