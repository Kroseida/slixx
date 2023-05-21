package syncnetwork

import (
	"kroseida.org/slixx/internal/satellite/application"
	satelliteProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/satellite"
	supervisorProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/pkg/satellite"
	"kroseida.org/slixx/pkg/satellite/protocol"
	handshakeProtocol "kroseida.org/slixx/pkg/satellite/protocol/handshake"
)

var server satellite.Server

func Listen() error {
	server := satellite.Server{
		BindAddress: application.CurrentSettings.Satellite.Network.BindAddress,
		Token:       application.CurrentSettings.Satellite.AuthenticationToken,
		Handler: map[string]protocol.Handler{
			protocol.HandshakeProtocol: &handshakeProtocol.ServerHandler{
				Token: application.CurrentSettings.Satellite.AuthenticationToken,
			},
			protocol.SupervisorProtocol: &supervisorProtocol.Handler{},
			protocol.SatelliteProtocol:  &satelliteProtocol.Handler{},
		},
		Logger: application.Logger,
	}
	err := server.Listen()
	if err != nil {
		return err
	}
	return nil
}
