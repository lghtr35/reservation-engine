/*
 * Any operation that does not mutate the database belongs to 'queries'.
 */
package queries

import (
	"errors"
	"fmt"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilterCustomersQuery struct {
	db     *gorm.DB
	logger *zerolog.Logger
	ids    *[]string
	name   *string
	models.Pagination
}

func NewFilterCustomersQuery(db *gorm.DB, logger *zerolog.Logger, ids *[]string, name *string, pagination models.Pagination) *FilterCustomersQuery {
	return &FilterCustomersQuery{db: db, logger: logger, name: name, ids: ids, Pagination: pagination}
}

func (s *FilterCustomersQuery) Execute() (any, error) {
	s.logger.Debug().Msg("FilterCustomersQuery: Started")
	q := s.db.Model(models.Customer{})
	if s.ids != nil && len(*s.ids) > 0 {
		q = q.Where("id IN ?", *s.ids)
	}
	if s.name != nil && *s.name != "" {
		q = q.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *s.name))
	}
	offset := s.Pagination.Offset()

	var customers []models.Customer
	res := q.Offset(offset).Limit(int(s.Size)).Find(&customers)
	if res.Error != nil {
		return models.NewPaginationResponse(customers, 0, 0), res.Error
	}

	var totalCount int64
	res = q.Count(&totalCount)
	if res.Error != nil {
		return models.NewPaginationResponse(customers, 0, 0), res.Error
	}

	s.logger.Debug().Msg("FilterCustomersQuery: Finished with success")
	return models.NewPaginationResponse(customers, totalCount, s.Page), nil

}

type ReadCustomerQuery struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewReadCustomerQuery(db *gorm.DB, logger *zerolog.Logger, id string) *ReadCustomerQuery {
	return &ReadCustomerQuery{db: db, logger: logger, id: id}
}

func (s *ReadCustomerQuery) Execute() (any, error) {
	if s.id == "" {
		return models.Customer{}, errors.New("ReadCustomerQuery: Tried to read one with empty id")
	}
	s.logger.Debug().Msg("ReadCustomerQuery: ReadOne started")

	var customer models.Customer
	res := s.db.Model(models.Customer{}).Preload(clause.Associations).First(&customer, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("ReadCustomerQuery: Could not find the customer with this id: %s", s.id)
		}
		return "", res.Error
	}

	s.logger.Debug().Msg("ReadCustomerQuery: ReadOne finished with success")
	return customer, nil
}
