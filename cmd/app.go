package cmd

import (
	"self-update/updater"
	"self-update/web"
)

func Start() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	updater := updater.Updater{
		CurrentVersion: config.Version,
		Repo:           config.Repo,
	}

	server := web.Server{
		Host:    config.Server.Host,
		Port:    config.Server.Port,
		Updater: updater,
	}

	return server.Serve()
}
