package main

import (
	"sync/atomic"
	"github.com/DavidMWeaver4/Chirpy/internal/database"

)

type apiConfig struct {
	fileserverHits atomic.Int32
	db			*database.Queries
}
