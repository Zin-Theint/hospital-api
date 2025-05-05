package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zin-Theint/hospital-api/internal/handler"
	"github.com/Zin-Theint/hospital-api/internal/middleware"
	"github.com/Zin-Theint/hospital-api/internal/repository"
	"github.com/Zin-Theint/hospital-api/internal/service"
)

func Setup(db *pgxpool.Pool, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// wiring
	staffRepo := repository.StaffRepo{DB: db}
	patientRepo := repository.PatientRepo{DB: db}
	staffSvc := service.StaffService{Repo: staffRepo}
	patientSvc := service.PatientService{Repo: patientRepo}

	authH := handler.AuthHandler{Service: staffSvc, JWTSecret: jwtSecret}
	patientH := handler.PatientHandler{Service: patientSvc}

	// routes
	r.POST("/staff/create", authH.Register)
	r.POST("/staff/login", authH.Login)

	auth := r.Group("/", middleware.Auth(jwtSecret))
	{
		auth.GET("/patient/search", patientH.Search)
	}

	return r
}
