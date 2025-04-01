// Implementation of URL record structure
// for JSON encoding and decoding

package url

import (
	"encoding/json"
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339

type DateTime struct {
	time.Time
}

type Record struct {
	Id        string   `json:"id"`
	URL       string   `json:"url"`
	ShortCode string   `json:"shortCode"`
	CreatedAt DateTime `json:"createdAt"`
	UpdatedAt DateTime `json:"updatedAt"`
	//	  "id": "1",
	//	  "url": "https://www.example.com/some/long/url",
	//	  "shortCode": "abc123",
	//	  "createdAt": "2021-09-01T12:00:00Z",
	//	  "updatedAt": "2021-09-01T12:00:00Z"
	AccessCount int `json:"accessCount,omitempty"`
}

type RecordStore interface {
	//

	InsertRecord(r Record) (id string, err error)

	SelectRecord(shortURL string) (rec Record, err error)

	// TODO
	// UpdateRecord
	// DEleteRecord
}

// Redifinitions for custom timestamps format instead of
// time.RFC3339Nano - default selected by `encoding/json` package
func (d *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(TimeFormat))
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	timestamp, err := time.Parse(TimeFormat, s)
	if err != nil {
		return fmt.Errorf("failed to parse time string %q: %v", s, err)
	}

	*d = DateTime{timestamp}
	return nil
}

func NewRecord(url string) *Record {
	return &Record{
		Id:        "newId", // TODO
		URL:       url,
		ShortCode: CreateShortURL(),
		CreatedAt: DateTime{time.Now().UTC()},
		UpdatedAt: DateTime{time.Now().UTC()},
	}
}

func (r *Record) Update() {
	r.ShortCode = CreateShortURL()
	r.UpdatedAt = DateTime{time.Now().UTC()}
}
