package server

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// GET /v1/proxy/youtube/search?q=...&maxResults=5 — проксирует ответ YouTube Data API без изменений.
func (s *Server) proxyYouTubeSearch(w http.ResponseWriter, r *http.Request) {
	key := os.Getenv("YOUTUBE_API_KEY")
	if key == "" {
		http.Error(w, `{"error":"YOUTUBE_API_KEY not configured"}`, http.StatusServiceUnavailable)
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, `{"error":"missing q"}`, http.StatusBadRequest)
		return
	}
	max := r.URL.Query().Get("maxResults")
	if max == "" {
		max = "5"
	}
	u, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.RawQuery = url.Values{
		"part":       {"snippet"},
		"q":          {q},
		"type":       {"video"},
		"maxResults": {max},
		"key":        {key},
	}.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Error().Err(err).Msg("youtube proxy")
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	copyProxyResponse(w, resp)
}

// GET /v1/proxy/giphy/search?q=...&limit=50 — проксирует ответ Giphy API без изменений.
func (s *Server) proxyGiphySearch(w http.ResponseWriter, r *http.Request) {
	key := os.Getenv("GIPHY_API_KEY")
	if key == "" {
		http.Error(w, `{"error":"GIPHY_API_KEY not configured"}`, http.StatusServiceUnavailable)
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, `{"error":"missing q"}`, http.StatusBadRequest)
		return
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50"
	}
	if _, err := strconv.Atoi(limit); err != nil {
		limit = "50"
	}
	u, err := url.Parse("https://api.giphy.com/v1/gifs/search")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.RawQuery = url.Values{
		"api_key": {key},
		"q":       {q},
		"limit":   {limit},
		"rating":  {"pg"},
	}.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Error().Err(err).Msg("giphy proxy")
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	copyProxyResponse(w, resp)
}

func copyProxyResponse(w http.ResponseWriter, resp *http.Response) {
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/json"
	}
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error().Err(err).Msg("proxy copy body")
	}
}
