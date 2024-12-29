package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lghtr35/reservation-engine/commands"
	"github.com/lghtr35/reservation-engine/models"
	"github.com/lghtr35/reservation-engine/queries"
	"github.com/lghtr35/reservation-engine/util"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	logger *zerolog.Logger
	hasher *util.Hasher
}

// Queries
func (h *Handler) ReadAllCustomers(c *gin.Context) {
	var request models.ReadAllCustomers
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := queries.NewFilterCustomersQuery(h.db, h.logger, request.IDs, request.Name, request.Pagination)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ReadAllSources(c *gin.Context) {
	var request models.ReadAllSources
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := queries.NewFilterSourcesQuery(h.db, h.logger, request.IDs, request.Name, request.Pagination)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ReadAllReservations(c *gin.Context) {
	var request models.ReadAllReservations
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := queries.NewFilterReservationsQuery(h.db, h.logger, request.IDs, request.ReserveeID, request.ReserverID, request.SourceID, request.Pagination)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ReadCustomer(c *gin.Context) {
	id := c.Param("id")

	q := queries.NewReadCustomerQuery(h.db, h.logger, id)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ReadSource(c *gin.Context) {
	id := c.Param("id")

	q := queries.NewReadSourceQuery(h.db, h.logger, id)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ReadReservation(c *gin.Context) {
	id := c.Param("id")

	q := queries.NewReadReservationQuery(h.db, h.logger, id)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Commands
func (h *Handler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	q := commands.NewDeleteCustomerCommand(h.db, h.logger, id)

	_, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (h *Handler) DeleteSource(c *gin.Context) {
	id := c.Param("id")

	q := commands.NewDeleteSourceCommand(h.db, h.logger, id)

	_, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (h *Handler) DeleteReservation(c *gin.Context) {
	id := c.Param("id")

	q := commands.NewDeleteReservationCommand(h.db, h.logger, id)

	_, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var request models.CreateCustomer
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewCreateCustomerCommand(h.db, h.logger, request.Name, request.Company, request.Email)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sQ := commands.NewCreateSecretCommand(h.db, h.logger, h.hasher, res)
	_, err = sQ.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateSource(c *gin.Context) {
	var request models.CreateSource
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewCreateSourceCommand(h.db, h.logger, request.Name, request.MaxPossibleDuration, request.CustomerID)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tQ := commands.NewCreateApiTokenCommand(h.db, h.logger, h.hasher, request.CustomerID, res)
	_, err = tQ.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateReservation(c *gin.Context) {
	var request models.CreateReservation
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewCreateReservationCommand(h.db, h.logger, request.From, request.To, request.ReserverID, request.ReserveeID)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	var request models.UpdateCustomer
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewUpdateCustomerCommand(h.db, h.logger, request.ID, request.Name, request.Email, request.Company, request.MaxSourceLimit)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateSource(c *gin.Context) {
	var request models.UpdateSource
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewUpdateSourceCommand(h.db, h.logger, request.ID, request.Name, request.MaxPossibleDuration)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateReservation(c *gin.Context) {
	var request models.UpdateReservation
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := commands.NewUpdateReservationCommand(h.db, h.logger, request.ID, request.From, request.To)

	res, err := q.Execute()
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
