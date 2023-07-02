package syncnetwork

import (
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/internal/satellite/syncnetwork/manager"
	satelliteProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/satellite"
	supervisorProtocol "kroseida.org/slixx/internal/satellite/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	handshakeProtocol "kroseida.org/slixx/pkg/syncnetwork/protocol/handshake"
	"time"
)

func Listen() error {
	manager.Server = &syncnetwork.Server{
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
	err := manager.Server.Listen()
	if err != nil {
		return err
	}
	return nil
}

func SyncLoop() {
	for {
		action.SyncLogsToSupervisor()
		time.Sleep(5 * time.Second)
	}
}
