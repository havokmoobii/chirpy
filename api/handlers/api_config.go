package api

import (
	"sync/atomic"

	"github.com/havokmoobii/chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	Secret         string
	PolkaKey       string
}