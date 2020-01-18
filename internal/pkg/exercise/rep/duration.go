package rep

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"time"
)

type Duration struct {
	ID       identifier.ID `json:"id" bson:"id"`
	Duration time.Duration `json:"duration" bson:"duration"`
}

func (t *Duration) SetID(id identifier.ID) {
	t.ID = id
}

func (t *Duration) Type() Type {
	return DurationRepType
}

func (t *Duration) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     Type          `json:"type"`
		Duration time.Duration `json:"duration"`
	}{
		Type:     t.Type(),
		Duration: t.Duration,
	})
}

func (t *Duration) ToBSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     Type          `bson:"type"`
		Duration time.Duration `bson:"duration"`
	}{
		Type:     t.Type(),
		Duration: t.Duration,
	})
}
