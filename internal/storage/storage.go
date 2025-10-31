package storage

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type InMemoryStore struct {
	mu          sync.RWMutex
	urlToCode   map[string]string
	codeToURL   map[string]string
	domainCount map[string]int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		urlToCode:   make(map[string]string),
		codeToURL:   make(map[string]string),
		domainCount: make(map[string]int),
	}
}

var ErrNotFound = errors.New("not found")

func (s *InMemoryStore) GetByCode(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.codeToURL[code]
	if !ok {
		return "", ErrNotFound
	}
	return u, nil
}

func (s *InMemoryStore) Shorten(rawurl string) (string, error) {
	parsed, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if code, ok := s.urlToCode[rawurl]; ok {
		s.domainCount[domainFromHost(parsed.Host)]++
		return code, nil
	}

	code := encodeBase62(hashToUint64(rawurl))[:8]
	for {
		if _, exists := s.codeToURL[code]; !exists {
			break
		}
		code = encodeBase62(hashToUint64(rawurl + "1"))[:8]
	}

	s.urlToCode[rawurl] = code
	s.codeToURL[code] = rawurl
	s.domainCount[domainFromHost(parsed.Host)]++
	return code, nil
}

func domainFromHost(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	return host
}

func (s *InMemoryStore) TopDomains(n int) []struct{ Domain string; Count int } {
	s.mu.RLock()
	defer s.mu.RUnlock()
	type pair struct {
		Domain string
		Count  int
	}
	pairs := make([]pair, 0, len(s.domainCount))
	for d, c := range s.domainCount {
		pairs = append(pairs, pair{d, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].Domain < pairs[j].Domain
		}
		return pairs[i].Count > pairs[j].Count
	})
	if len(pairs) > n {
		pairs = pairs[:n]
	}
	return pairs
}

func hashToUint64(s string) uint64 {
	h := sha256.Sum256([]byte(s))
	return binary.LittleEndian.Uint64(h[:8])
}

const b62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append(b, b62[n%62])
		n /= 62
	}
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
