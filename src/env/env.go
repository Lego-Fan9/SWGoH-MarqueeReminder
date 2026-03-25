package env

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	DISCORD_WEBHOOK string //nolint:revive
	CUSTOM_FORMAT   string //nolint:revive
	COMLINK_URL     string //nolint:revive
	ENV_PATH        string //nolint:revive
	PING_ROLE       string //nolint:revive
	// SWGOH_AE_URL    string //nolint:revive
)

var TESTING = os.Getenv("TESTING")

func init() {
	TestHandler()

	if TESTING != "1" {
		Init()
	}
}

func Init() {
	temp_env := os.Getenv("ENV_PATH") //nolint:revive
	if temp_env != "" {
		ENV_PATH = temp_env
	} else {
		ENV_PATH = ".env"
	}

	if ENV_PATH != "NONE" && ENV_PATH != "" {
		err := godotenv.Load(ENV_PATH)
		if err != nil {
			log.Warnf("Error loading .env: %v", err)
		}
	}

	DISCORD_WEBHOOK = os.Getenv("DISCORD_WEBHOOK")
	if DISCORD_WEBHOOK == "" {
		log.Fatal("Failed to find env: DISCORD_WEBHOOK")
	}

	CUSTOM_FORMAT = os.Getenv("CUSTOM_FORMAT")

	COMLINK_URL = os.Getenv("COMLINK_URL")
	if COMLINK_URL == "" {
		log.Fatal("Failed to find env: COMLINK_URL")
	}

	/* Make sure to add to .env.example
	SWGOH_AE_URL = os.Getenv("SWGOH_AE_URL")
	if SWGOH_AE_URL == "" {
		log.Fatal("Failed to find env: SWGOH_AE_URL")
	}*/

	PING_ROLE = os.Getenv("PING_ROLE")
	if PING_ROLE == "" {
		log.Fatal("Failed to find env: PING_ROLE")
	}
}

// Handles a test to try and find an env, if TESTING == "1"
func TestHandler() {
	if os.Getenv("TESTING") != "1" {
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil {
			ENV_PATH = filepath.Join(dir, ".env")

			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("Could not find env")
		}

		dir = parent
	}
}
