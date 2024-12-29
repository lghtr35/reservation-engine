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

type FilterSourcesQuery struct {
	db     *gorm.DB
	logger *zerolog.Logger
	ids    *[]string
	name   *string
	models.Pagination
}

func NewFilterSourcesQuery(db *gorm.DB, logger *zerolog.Logger, ids *[]string, name *string, pagination models.Pagination) *FilterSourcesQuery {
	return &FilterSourcesQuery{db: db, logger: logger, name: name, ids: ids, Pagination: pagination}
}

func (s *FilterSourcesQuery) Execute() (any, error) {
	s.logger.Debug().Msg("FilterSourcesQuery: Started")
	q := s.db.Model(models.Source{})
	if s.ids != nil && len(*s.ids) > 0 {
		q = q.Where("id IN ?", *s.ids)
	}
	if s.name != nil && *s.name != "" {
		q = q.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *s.name))
	}
	offset := s.Pagination.Offset()

	var sources []models.Source
	res := q.Offset(offset).Limit(int(s.Size)).Find(&sources)
	if res.Error != nil {
		return models.NewPaginationResponse(sources, 0, 0), res.Error
	}

	var totalCount int64
	res = q.Count(&totalCount)
	if res.Error != nil {
		return models.NewPaginationResponse(sources, 0, 0), res.Error
	}

	s.logger.Debug().Msg("FilterSourcesQuery: Finished with success")
	return models.NewPaginationResponse(sources, totalCount, s.Page), nil

}

type ReadSourceQuery struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewReadSourceQuery(db *gorm.DB, logger *zerolog.Logger, id string) *ReadSourceQuery {
	return &ReadSourceQuery{db: db, logger: logger, id: id}
}

func (s *ReadSourceQuery) Execute() (any, error) {
	if s.id == "" {
		return models.Source{}, errors.New("ReadSourceQuery: Tried to read one with empty id")
	}
	s.logger.Debug().Msg("ReadSourceQuery: ReadOne started")

	var source models.Source
	res := s.db.Model(models.Source{}).Preload(clause.Associations).First(&source, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("ReadSourceQuery: Could not find the source with this id: %s", s.id)
		}
		return "", res.Error
	}

	s.logger.Debug().Msg("ReadSourceQuery: ReadOne finished with success")
	return source, nil
}
