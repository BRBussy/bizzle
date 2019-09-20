package mongo

import (
	mongoBSON "go.mongodb.org/mongo-driver/bson"
	mongoDriverOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	Limit   int64     `json:"limit"`
	Offset  int64     `json:"offset"`
	Sorting []Sorting `json:"Sorting"`
}

type Sorting struct {
	Field     string
	SortOrder SortOrder
}

func (s *Sorting) IsValid() error {
	reasonsInvalid := make([]string, 0)
	if !(s.SortOrder == SortOrderAscending || s.SortOrder == SortOrderDescending) {
		reasonsInvalid = append(reasonsInvalid, "invalid sort order: "+s.SortOrder.String())
	}
	if s.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field blank")
	}
	if len(reasonsInvalid) > 0 {
		return ErrSortingInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (s *Sorting) ToMongoSortFormat() (*mongoBSON.E, error) {
	if err := s.IsValid(); err != nil {
		return nil, err
	}
	switch s.SortOrder {
	case SortOrderDescending:
		return &mongoBSON.E{
			Key:   s.Field,
			Value: s.SortOrder,
		}, nil
	case SortOrderAscending:
		fallthrough
	default:
		return &mongoBSON.E{
			Key:   s.Field,
			Value: s.SortOrder,
		}, nil
	}
}

type SortOrder string

func (s SortOrder) String() string {
	return string(s)
}

const SortOrderAscending SortOrder = "asc"
const SortOrderDescending SortOrder = "desc"

func (q Query) IsValid() error {
	reasonsInvalid := make([]string, 0)
	for i := range q.Sorting {
		if err := q.Sorting[i].IsValid(); err != nil {
			reasonsInvalid = append(reasonsInvalid, err.Error())
		}
	}
	if len(reasonsInvalid) > 0 {
		return ErrQueryInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (q Query) ToMongoFindOptions() (*mongoDriverOptions.FindOptions, error) {
	// get sorting
	sorting := mongoBSON.D{}
	for i := range q.Sorting {
		sort, err := q.Sorting[i].ToMongoSortFormat()
		if err != nil {
			return nil, err
		}
		sorting = append(sorting, *sort)
	}
	// create find options
	findOptions := new(mongoDriverOptions.FindOptions)

	// populate find options
	findOptions.SetSort(sorting)
	findOptions.SetSkip(q.Offset)
	if q.Limit > 0 {
		findOptions.SetLimit(q.Limit)
	}

	return findOptions, nil
}
