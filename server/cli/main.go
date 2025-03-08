package main

import (
	"flag"
	"log"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/server"
	"github.com/mattermost/focalboard/server/services/auth"
	"github.com/mattermost/focalboard/server/services/config"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func main() {
	username := flag.String("username", "", "Username of new user")
	password := flag.String("password", "", "Password of new user")
	email := flag.String("email", "", "Email of new user")
	role := flag.String("role", "", "Role of new user")
	pConfigFilePath := flag.String(
		"config",
		"",
		"Location of the JSON config file",
	)
	flag.Parse()

	if *username == "" || *password == "" || *email == "" || *role == "" {
		log.Fatal("All params are required")
		return
	}

	config, err := config.ReadConfigFile(*pConfigFilePath)
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
		return
	}

	logger, _ := mlog.NewLogger()
	cfgJSON := config.LoggingCfgJSON
	if config.LoggingCfgFile == "" && cfgJSON == "" {
		// if no logging defined, use default config (console output)
		cfgJSON = defaultLoggingConfig()
	}
	err = logger.Configure(config.LoggingCfgFile, cfgJSON, nil)
	if err != nil {
		log.Fatal("Error in config file for logger: ", err)
		return
	}
	defer func() { _ = logger.Shutdown() }()

	if logger.HasTargets() {
		restore := logger.RedirectStdLog(mlog.LvlInfo, mlog.String("src", "stdlog"))
		defer restore()
	}

	db, err := server.NewStore(config, false, logger)
	if err != nil {
		logger.Fatal("server.NewStore ERROR", mlog.Err(err))
	}

	dbRole, err := db.GetRoleByName(*role)
	if err != nil {
		logger.Fatal("Role not found")
	}

	_, err = db.CreateUser(&model.User{
		ID:          utils.NewID(utils.IDTypeUser),
		Username:    *username,
		Email:       *email,
		Password:    auth.HashPassword(*password),
		MfaSecret:   "",
		AuthService: config.AuthMode,
		AuthData:    "",
		RoleID:      dbRole.ID,
	})
}

func defaultLoggingConfig() string {
	return `
	{
		"def": {
			"type": "console",
			"options": {
				"out": "stdout"
			},
			"format": "plain",
			"format_options": {
				"delim": " ",
				"min_level_len": 5,
				"min_msg_len": 40,
				"enable_color": true,
				"enable_caller": true
			},
			"levels": [
				{"id": 5, "name": "debug"},
				{"id": 4, "name": "info", "color": 36},
				{"id": 3, "name": "warn"},
				{"id": 2, "name": "error", "color": 31},
				{"id": 1, "name": "fatal", "stacktrace": true},
				{"id": 0, "name": "panic", "stacktrace": true}
			]
		}
	}`
}
