package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lghtr35/reservation-engine/models"
	"github.com/lghtr35/reservation-engine/util"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := zerolog.New(os.Stdout)

	configuration := models.Configuration{}
	err := configuration.ReadAndFillSelf(logger)
	if err != nil {
		panic(err)
	}
	hasher, err := util.NewHasher(&configuration)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(configuration.DbConnectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&models.Source{},
		&models.Secret{},
		&models.ApiToken{},
		&models.Reservation{},
		&models.Customer{},
	)
	if err != nil {
		panic(err)
	}

	h := Handler{
		logger: &logger,
		db:     db,
		hasher: hasher,
	}

	g := gin.New()
	api := g.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			jwt := v1.Group("/")
			{
				jwt.Use(jwtAuthMiddleware(&configuration, db, &logger))
				// Customers
				jwt.GET("/customers", h.ReadAllCustomers)
				jwt.POST("/customers", h.CreateCustomer)
				jwt.PATCH("/customers", h.UpdateCustomer)
				jwt.GET("/customers/:id", h.ReadCustomer)
				jwt.DELETE("/customers/:id", h.DeleteCustomer)
			}
			apiKey := v1.Group("/")
			{
				apiKey.Use(apiKeyAuthMiddleware(db, &logger))
				// Reservations
				apiKey.GET("/reservations", h.ReadAllReservations)
				apiKey.POST("/reservations", h.CreateReservation)
				apiKey.PATCH("/reservations", h.UpdateReservation)
				apiKey.GET("/reservations/:id", h.ReadReservation)
				apiKey.DELETE("/reservations/:id", h.DeleteReservation)
				// Sources
				apiKey.GET("/sources", h.ReadAllSources)
				apiKey.POST("/sources", h.CreateSource)
				apiKey.PATCH("/sources", h.UpdateSource)
				apiKey.GET("/sources/:id", h.ReadSource)
				apiKey.DELETE("/sources/:id", h.DeleteSource)
			}
		}
	}
	g.Run()
}
