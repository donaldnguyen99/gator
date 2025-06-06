package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/donaldnguyen99/gator/internal/projectpath"
	"github.com/joho/godotenv"
)

const (
	configFileName = ".gatorconfig.json"
)

var (
	dbURL = getDbURL()
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func NewConfig() *Config {
	return &Config{}
}

func Read() (*Config, error) {
	file, err := openToReadConfigFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	file, err := openToWriteConfigFile(configFileName)
	if err != nil {
		return fmt.Errorf("unable to set user %v: %w", user, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("unable to set user %v: %w", user, err)
	}
	
	return nil
}

func getDbURL() string {
	envFile, err := godotenv.Read(filepath.Join(projectpath.Root, ".env"))
	if err != nil {
		value, exists := os.LookupEnv("GATOR_POSTGRES_URL")
		if exists {
			return value + "?sslmode=disable"
		} else {
			panic(fmt.Errorf("env variable GATOR_POSTGRES_URL not set"))
		}
	}

	postgresUser     := envFile["POSTGRES_USER"]
	postgresPassword := envFile["POSTGRES_PASSWORD"]
	postgresDb       := envFile["POSTGRES_DB"]
	postgresHost     := envFile["POSTGRES_HOST"]
	postgresPort     := envFile["POSTGRES_PORT"]
	postgresSslmode  := envFile["POSTGRES_SSLMODE"]
	if postgresSslmode == "disable" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDb, postgresSslmode)
	} else {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDb)
	}
}


func createDefaultConfigFile(configFileName string) error {
	file, err := openToWriteConfigFile(configFileName)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// encoding with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&struct {
		DbUrl string `json:"db_url"`
	}{
		DbUrl: dbURL,
	})
	if err != nil {
		return fmt.Errorf("failed to encode config file: %w", err)
	}
	return nil
}

func openToReadConfigFile(configFileName string) (*os.File, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(userHomeDir, configFileName)

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Config file not found at %s, creating default config file...\n", filepath)
		if err2 := createDefaultConfigFile(configFileName); err2 != nil {
			return nil, err2
		}
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func openToWriteConfigFile(configFileName string) (*os.File, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(userHomeDir, configFileName)

	file, err := os.OpenFile(
		filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644,
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}
