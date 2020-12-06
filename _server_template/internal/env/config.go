package env

import (
	"time"
)

// Config is struct with unmarshalled application config
type Config struct {
	Control struct {
		// Time for gracefull app exit
		DieTimeout time.Duration
	}
	Log struct {
		// Where to put logs
		File string
		// Global log mask: panic|fatal|error|warn|info|debug|trace
		Mask string
	}

	Router struct {
		GenerateDoc bool
		Addr        string
	}

	Dummy struct {
		StringSlice []string
	}

	Auth struct {
		// JSON file from https://console.firebase.google.com/
		FirebaseAccountKeyFile string
	}
}
