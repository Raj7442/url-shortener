package storage

import "testing"

func TestShortenAndGet(t *testing.T) {
	s := NewInMemoryStore()
	url := "https://example.com/path"
	code, err := s.Shorten(url)
	if err != nil {
		t.Fatalf("Shorten failed: %v", err)
	}
	got, err := s.GetByCode(code)
	if err != nil {
		t.Fatalf("GetByCode failed: %v", err)
	}
	if got != url {
		t.Fatalf("expected %s got %s", url, got)
	}
}

func TestDeterministic(t *testing.T) {
	s := NewInMemoryStore()
	url := "https://example.com/test"
	c1, _ := s.Shorten(url)
	c2, _ := s.Shorten(url)
	if c1 != c2 {
		t.Fatalf("codes differ for same URL: %s vs %s", c1, c2)
	}
}

func TestTopDomains(t *testing.T) {
	s := NewInMemoryStore()
	s.Shorten("https://youtube.com/watch?v=1")
	s.Shorten("https://youtube.com/watch?v=2")
	s.Shorten("https://udemy.com/course/1")
	list := s.TopDomains(3)
	if len(list) < 2 {
		t.Fatal("expected multiple domains")
	}
}
