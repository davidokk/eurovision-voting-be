package server

import (
	"eurovision-voting/internal/service"
	"errors"
	"net/http"
)

func (s *Server) uploadMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := r.ParseMultipartForm(12 << 20); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error()})
		return
	}
	kind := r.FormValue("kind")
	file, hdr, err := r.FormFile("file")
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "file required"})
		return
	}
	defer file.Close()

	url, err := s.service.UploadMedia(r.Context(), userID, kind, file, hdr.Header.Get("Content-Type"))
	if err != nil {
		if errors.Is(err, service.ErrMediaNotConfigured) || errors.Is(err, service.ErrInvalidMediaKind) {
			EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error()})
			return
		}
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, map[string]string{"url": url})
}
