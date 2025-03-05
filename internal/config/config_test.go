package config

import (
	"encoding/json"
	"testing"
)

func createDefaultConfigFile() {
	file, err := openToWriteConfigFile()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// encoding with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&struct {
		DbUrl string `json:"db_url"`
	}{
		DbUrl: "postgres://example",
	})
	if err != nil {
		panic(err)
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Error("expected config to be created")
	} else {
		if config.CurrentUserName != "" {
			t.Errorf("expected current user name to be empty, got %v", config.CurrentUserName)
		}
		if config.DbUrl != "" {
			t.Errorf("expected db url to be empty, got %v", config.DbUrl)
		}
	}
}

func TestRead(t *testing.T) {
	createDefaultConfigFile()
	config, err := Read()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if config == nil {
		t.Error("expected config to be created")
	} else {
		if config.DbUrl != "postgres://example" {
			t.Errorf("expected db url to be postgres://example, got %v", config.DbUrl)
		}
		if config.CurrentUserName != "" {
			t.Errorf("expected current user name to be empty, got %v", config.CurrentUserName)
		}
	}
}

func TestSetUser(t *testing.T) {
	createDefaultConfigFile()
	config, err := Read()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if config == nil {
		t.Error("expected config to be created")
		return
	}

	if err := config.SetUser("test"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	defer createDefaultConfigFile()
	if config.CurrentUserName != "test" {
		t.Errorf("expected current user name to be test, got %v", config.CurrentUserName)
	}

	config, err = Read()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if config == nil {
		t.Error("expected config to be created")
		return
	}
	if config.CurrentUserName != "test" {
		t.Errorf("expected current user name to be test, got %v", config.CurrentUserName)
	}
}
