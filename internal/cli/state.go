package cli

import (
	"github.com/donaldnguyen99/gator/internal/config"
	"github.com/donaldnguyen99/gator/internal/database"
)

type state struct {
	db *database.Queries
	config *config.Config
}