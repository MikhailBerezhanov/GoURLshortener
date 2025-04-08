// Implementation of URL record structure
// for JSON encoding and decoding

package url

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339

type DateTime struct {
	time.Time
}

type Record struct {
	Id        string   `json:"id" bson:"_id"`
	URL       string   `json:"url" bson:"url"`
	ShortCode string   `json:"shortCode" bson:"shortCode"`
	CreatedAt DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt DateTime `json:"updatedAt" bson:"updatedAt"`
	//	  "id": "1",
	//	  "url": "https://www.example.com/some/long/url",
	//	  "shortCode": "abc123",
	//	  "createdAt": "2021-09-01T12:00:00Z",
	//	  "updatedAt": "2021-09-01T12:00:00Z"
	AccessCount int `json:"accessCount,omitempty" bson:"accessCount,omitempty"`
}

var ErrRecordNotExist = errors.New("record does not exist")

type RecordStore interface {
	//

	InsertRecord(r *Record) error

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

func generateId(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewRecord(url string) *Record {
	return &Record{
		Id:        generateId(10),
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

func (r *Record) String() string {
	if data, err := json.Marshal(r); err == nil {
		return string(data)
	}

	return fmt.Sprintf("%v", *r)
}
