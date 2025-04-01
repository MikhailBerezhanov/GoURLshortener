// Basic implementation of memory db

package mem_db

import (
	"fmt"
	"url_shortener/url"
)

type MemDB struct {
	data map[string]url.Record
	id   int
}

func NewMemDB() *MemDB {
	return &MemDB{make(map[string]url.Record), 0}
}

// TODO: mutual access from differenr goroutines

func (m *MemDB) InsertRecord(r url.Record) (id string, err error) {
	// TODO: set id
	m.data[r.ShortCode] = r
	id = fmt.Sprintf("%d", m.id)
	m.id++
	return
}

func (m *MemDB) SelectRecord(shortURL string) (rec url.Record, err error) {
	if rec, ok := m.data[shortURL]; ok {
		return rec, nil
	}

	return url.Record{}, fmt.Errorf("no record for shortURL %q found", shortURL)
}
