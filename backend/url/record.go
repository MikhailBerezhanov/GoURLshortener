// Implementation of URL record structure
// for JSON encoding \ decoding

package url

import "time"

type Record struct {
	Id        string    `json:"id"`
	URL       string    `json:"url"`
	ShortCode string    `json:"shortCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	//	  "id": "1",
	//	  "url": "https://www.example.com/some/long/url",
	//	  "shortCode": "abc123",
	//	  "createdAt": "2021-09-01T12:00:00Z",
	//	  "updatedAt": "2021-09-01T12:00:00Z"
	AccessCount int `json:"accessCount,omitempty"`
}

// func (r )
// data, err := json.Marshal(rec)
// 	if err != nil {
// 		return err
// 	}

// 	if err := json.Unmarshal([]byte(recordText), &record); err != nil {
// 		return err
// 	}
