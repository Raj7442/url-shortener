package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Raj7442/url-shortener/internal/models"
	"github.com/Raj7442/url-shortener/internal/storage"
)

type Handler struct {
	store   *storage.InMemoryStore
	baseURL string
}

func NewHandler(s *storage.InMemoryStore, base string) *Handler {
	return &Handler{store: s, baseURL: strings.TrimRight(base, "/")}
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}

	code, err := h.store.Shorten(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	short := fmt.Sprintf("%s/%s", h.baseURL, code)
	json.NewEncoder(w).Encode(models.ShortenResponse{ShortURL: short})
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	list := h.store.TopDomains(3)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "" || strings.HasPrefix(path, "api/") {
		http.NotFound(w, r)
		return
	}

	orig, err := h.store.GetByCode(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, orig, http.StatusFound)
}
