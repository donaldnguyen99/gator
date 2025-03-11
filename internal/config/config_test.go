package config

import (
	"testing"
)

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
	config, err := Read()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if config == nil {
		t.Error("expected config to be created")
	}
}

// This test will modify the config file
func TestSetUser(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if config == nil {
		t.Error("expected config to be created")
		return
	}
	// save old user name
	oldUser := config.CurrentUserName


	if err := config.SetUser("test"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
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

	// restore old user name
	if err := config.SetUser(oldUser); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
