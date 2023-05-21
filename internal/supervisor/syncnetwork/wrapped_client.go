package syncnetwork

import (
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/satellite"
)

type WrappedClient struct {
	model  model.Satellite
	client *satellite.Client
}
