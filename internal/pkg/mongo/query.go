package mongo

import (
	mongoBSON "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func CompareQuery(q1, q2 Query) bool {
	if q1.Limit != q2.Limit {
		return false
	}
	if q1.Offset != q2.Offset {
		return false
	}
	if len(q1.Sorting) != len(q2.Sorting) {
		return false
	}

	// for every sorting entry in q1
nextQ1SortingEntry:
	for q1SortingEntryI := range q1.Sorting {
		// look for it in q2
		for q2SortingEntryJ := range q2.Sorting {
			if q1.Sorting[q1SortingEntryI] == q2.Sorting[q2SortingEntryJ] {
				// if it is found, go to next q1 sorting entry
				continue nextQ1SortingEntry
			}
		}
		// if execution reaches here then q1SortingEntryI was not found in q2
		return false
	}
	// if execution reaches here every sorting entry in q1 was found in q2
	return true
}

// ToPipelineStages converts the Query object into mongoDb aggregation pipeline stages - helper function
func (q Query) ToPipelineStages() mongo.Pipeline {
	pipelineStages := make([]mongoBSON.D, 0)

	// if the query is blank, return blank
	if CompareQuery(q, Query{}) {
		return pipelineStages
	}

	// create sort stage
	if len(q.Sorting) > 0 {
		sorting := mongoBSON.D{}
		for _, sort := range q.Sorting {
			sortOrder := -1
			if sort.SortOrder == SortOrderAscending {
				sortOrder = 1
			}
			sorting = append(
				sorting,
				mongoBSON.E{
					Key:   sort.Field,
					Value: sortOrder,
				},
			)
		}
		pipelineStages = append(
			pipelineStages,
			mongoBSON.D{
				{
					Key:   "$sort",
					Value: sorting,
				},
			},
		)
	}

	// create offset stage if required
	if q.Offset > 0 {
		pipelineStages = append(
			pipelineStages,
			mongoBSON.D{
				{
					Key:   "$skip",
					Value: q.Offset,
				},
			},
		)
	}

	// create limit stage if required
	if q.Limit > 0 {
		pipelineStages = append(
			pipelineStages,
			mongoBSON.D{
				{
					Key:   "$limit",
					Value: q.Limit,
				},
			},
		)
	}

	return pipelineStages
}
