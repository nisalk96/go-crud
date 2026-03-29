package handlers

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"restapi/internal/models"
	"restapi/internal/repository"
	"restapi/internal/storage"

	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	Repo            *repository.MovieRepository
	Covers          *storage.CoverStorage
	PublicCoverPath string // e.g. /api/v1/files/covers
}

type movieResponse struct {
	ID                 string  `json:"id"`
	Title              string  `json:"title"`
	Rate               float64 `json:"rate"`
	Description        string  `json:"description,omitempty"`
	IMDbLink           string  `json:"imdbLink,omitempty"`
	TrailerYouTubeLink string  `json:"trailerYouTubeLink,omitempty"`
	CoverArtURL        string  `json:"coverArtURL,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

func (h *MovieHandler) toResponse(m *models.Movie) movieResponse {
	out := movieResponse{
		ID:                 m.ID.Hex(),
		Title:              m.Title,
		Rate:               m.Rate,
		Description:        m.Description,
		IMDbLink:           m.IMDbLink,
		TrailerYouTubeLink: m.TrailerYouTubeLink,
		CreatedAt:          m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if m.CoverArt != "" {
		out.CoverArtURL = storage.JoinPublicPath(h.PublicCoverPath, m.CoverArt)
	}
	return out
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	ct, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if strings.HasPrefix(ct, "multipart/form-data") {
		h.createMultipart(w, r)
		return
	}
	h.createJSON(w, r)
}

func (h *MovieHandler) createJSON(w http.ResponseWriter, r *http.Request) {
	var body models.MovieCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if err := validateMovieCreate(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	movie, err := h.Repo.Create(r.Context(), body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create movie"})
		return
	}
	writeJSON(w, http.StatusCreated, h.toResponse(movie))
}

func (h *MovieHandler) createMultipart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	rate := 0.0
	if s := strings.TrimSpace(r.FormValue("rate")); s != "" {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rate"})
			return
		}
		rate = v
	}
	in := models.MovieCreate{
		Title:              title,
		Rate:               rate,
		Description:        strings.TrimSpace(r.FormValue("description")),
		IMDbLink:           strings.TrimSpace(r.FormValue("imdbLink")),
		TrailerYouTubeLink: strings.TrimSpace(r.FormValue("trailerYouTubeLink")),
	}
	if err := validateMovieCreate(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	movie, err := h.Repo.Create(r.Context(), in)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create movie"})
		return
	}

	file, header, err := r.FormFile("cover")
	if err == nil && header != nil {
		name, serr := h.Covers.Save(file, header)
		if serr != nil {
			_ = h.Repo.Delete(r.Context(), movie.ID.Hex())
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": uploadErrMessage(serr)})
			return
		}
		patched, uerr := h.Repo.Update(r.Context(), movie.ID.Hex(), models.MovieUpdate{CoverArt: &name})
		if uerr != nil {
			_ = h.Covers.Remove(name)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save cover"})
			return
		}
		movie = patched
	}

	writeJSON(w, http.StatusCreated, h.toResponse(movie))
}

func validateMovieCreate(in *models.MovieCreate) error {
	if in.Title == "" {
		return errors.New("title is required")
	}
	if in.Rate < 0 || in.Rate > 10 {
		return errors.New("rate must be between 0 and 10")
	}
	if in.IMDbLink != "" && !strings.HasPrefix(in.IMDbLink, "http://") && !strings.HasPrefix(in.IMDbLink, "https://") {
		return errors.New("imdbLink must be a valid http(s) URL")
	}
	if in.TrailerYouTubeLink != "" && !strings.HasPrefix(in.TrailerYouTubeLink, "http://") && !strings.HasPrefix(in.TrailerYouTubeLink, "https://") {
		return errors.New("trailerYouTubeLink must be a valid http(s) URL")
	}
	return nil
}

func uploadErrMessage(err error) string {
	switch {
	case errors.Is(err, storage.ErrFileTooLarge):
		return "cover file too large"
	case errors.Is(err, storage.ErrNotImage):
		return "cover must be jpeg, png, webp, or gif"
	default:
		return "invalid cover upload"
	}
}

func (h *MovieHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movie, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load movie"})
		return
	}
	writeJSON(w, http.StatusOK, h.toResponse(movie))
}

func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	movies, err := h.Repo.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list movies"})
		return
	}
	out := make([]movieResponse, 0, len(movies))
	for i := range movies {
		out = append(out, h.toResponse(&movies[i]))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var patch models.MovieUpdate
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if patch.Title == nil && patch.Rate == nil && patch.Description == nil && patch.IMDbLink == nil &&
		patch.TrailerYouTubeLink == nil && patch.CoverArt == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no fields to update"})
		return
	}
	if patch.Rate != nil && (*patch.Rate < 0 || *patch.Rate > 10) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "rate must be between 0 and 10"})
		return
	}
	if patch.IMDbLink != nil && *patch.IMDbLink != "" &&
		!strings.HasPrefix(*patch.IMDbLink, "http://") && !strings.HasPrefix(*patch.IMDbLink, "https://") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "imdbLink must be a valid http(s) URL"})
		return
	}
	if patch.TrailerYouTubeLink != nil && *patch.TrailerYouTubeLink != "" &&
		!strings.HasPrefix(*patch.TrailerYouTubeLink, "http://") && !strings.HasPrefix(*patch.TrailerYouTubeLink, "https://") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "trailerYouTubeLink must be a valid http(s) URL"})
		return
	}

	prev, _ := h.Repo.GetByID(r.Context(), id)
	if patch.CoverArt != nil && *patch.CoverArt == "" && prev != nil && prev.CoverArt != "" {
		_ = h.Covers.Remove(prev.CoverArt)
	}

	movie, err := h.Repo.Update(r.Context(), id, patch)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update movie"})
		return
	}
	writeJSON(w, http.StatusOK, h.toResponse(movie))
}

func (h *MovieHandler) UploadCover(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}
	file, header, err := r.FormFile("cover")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cover file is required"})
		return
	}

	prev, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load movie"})
		return
	}

	name, serr := h.Covers.Save(file, header)
	if serr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": uploadErrMessage(serr)})
		return
	}

	movie, uerr := h.Repo.Update(r.Context(), id, models.MovieUpdate{CoverArt: &name})
	if uerr != nil {
		_ = h.Covers.Remove(name)
		if errors.Is(uerr, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update movie"})
		return
	}
	if prev.CoverArt != "" && prev.CoverArt != name {
		_ = h.Covers.Remove(prev.CoverArt)
	}

	writeJSON(w, http.StatusOK, h.toResponse(movie))
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	prev, _ := h.Repo.GetByID(r.Context(), id)
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete movie"})
		return
	}
	if prev != nil && prev.CoverArt != "" {
		_ = h.Covers.Remove(prev.CoverArt)
	}
	w.WriteHeader(http.StatusNoContent)
}
