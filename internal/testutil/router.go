package testutil

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zin-Theint/hospital-api/internal/router"
)

const jwtKey = "test-secret"

func NewTestRouter(db *pgxpool.Pool) *gin.Engine {
	return router.Setup(db, jwtKey)
}
