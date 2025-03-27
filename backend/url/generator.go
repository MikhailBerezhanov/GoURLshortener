// Short URL generator

package url

import (
	"math/rand/v2"
)

const (
	shortLen         = 6
	availableSymbols = "-+=0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CreateShortURL() string {
	shortURL := make([]byte, shortLen)
	for i := range shortURL {
		shortURL[i] = availableSymbols[rand.IntN(len(availableSymbols))]
	}

	return string(shortURL)
}
