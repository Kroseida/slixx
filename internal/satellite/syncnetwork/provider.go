package syncnetwork

import (
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/satellite/application"
	satelliteProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/satellite"
	supervisorProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	handshakeProtocol "kroseida.org/slixx/pkg/syncnetwork/protocol/handshake"
	"time"
)

var server *syncnetwork.Server

func Listen() error {
	server = &syncnetwork.Server{
		BindAddress: application.CurrentSettings.Satellite.Network.BindAddress,
		Token:       application.CurrentSettings.Satellite.AuthenticationToken,
		Handler: map[string]protocol.Handler{
			protocol.Handshake: &handshakeProtocol.ServerHandler{
				Token: application.CurrentSettings.Satellite.AuthenticationToken,
			},
			protocol.Supervisor: &supervisorProtocol.Handler{},
			protocol.Satellite:  &satelliteProtocol.Handler{},
		},
		Logger:  application.Logger,
		Version: common.CurrentVersion,
	}
	err := server.Listen()
	if err != nil {
		return err
	}
	return nil
}

func SyncLoop() {
	for {
		SyncLogsToSupervisor()
		time.Sleep(5 * time.Second)
	}
}
