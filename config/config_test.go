package config_test

import (
	"lxdexplorer-api/config"
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {

	c, _ := config.LoadConfig()

	if c == nil {
		t.Error("Expected a config object, got nil")
	}
}

func TestMongoDBURI(t *testing.T) {
	c, _ := config.LoadConfig()

	if c.MongoDB.URI != "mongodb://localhost:27017" {
		t.Errorf("Expected %s, got %s", "mongodb://localhost:27017", c.MongoDB.URI)
	}
}

func TestEnvMongoDBOverride(t *testing.T) {

	os.Setenv("MONGODB_URI", "mongodb://localhost:27018")

	c, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if c.MongoDB.URI != "mongodb://localhost:27018" {
		t.Errorf("Expected %s, got %s", "mongodb://localhost:27018", c.MongoDB.URI)
	}
}
