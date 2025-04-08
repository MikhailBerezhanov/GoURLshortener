// Basic memory db implementing url.RecordStore interface

package datastore

import (
	"url_shortener/url"
)

type MemDB struct {
	data map[string]url.Record
}

func NewMemDB() *MemDB {
	return &MemDB{make(map[string]url.Record)}
}

// TODO: mutual access from differenr goroutines

func (m *MemDB) InsertRecord(r *url.Record) error {
	m.data[r.ShortCode] = *r
	return nil
}

func (m *MemDB) SelectRecord(shortURL string) (rec url.Record, err error) {
	if rec, ok := m.data[shortURL]; ok {
		return rec, nil
	}

	return url.Record{}, url.ErrRecordNotExist
}
