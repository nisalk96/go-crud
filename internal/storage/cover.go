package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileTooLarge = errors.New("file too large")
	ErrNotImage     = errors.New("only jpeg, png, webp, and gif images are allowed")
)

type CoverStorage struct {
	Dir      string
	MaxBytes int64
}

// Save reads the multipart file, validates type, writes to disk, returns stored filename.
func (s *CoverStorage) Save(file multipart.File, header *multipart.FileHeader) (filename string, err error) {
	if header.Size > 0 && header.Size > s.MaxBytes {
		return "", ErrFileTooLarge
	}
	ext, ok := extFromHeader(header)
	if !ok {
		return "", ErrNotImage
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, s.MaxBytes+1))
	if err != nil {
		return "", err
	}
	if int64(len(data)) > s.MaxBytes {
		return "", ErrFileTooLarge
	}
	if len(data) == 0 {
		return "", ErrNotImage
	}

	var nameRand [16]byte
	if _, err := rand.Read(nameRand[:]); err != nil {
		return "", err
	}
	name := hex.EncodeToString(nameRand[:]) + ext
	dst := filepath.Join(s.Dir, name)
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", err
	}
	return name, nil
}

func (s *CoverStorage) Remove(filename string) error {
	if filename == "" {
		return nil
	}
	if strings.Contains(filename, "..") || strings.Contains(filename, string(filepath.Separator)) || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return errors.New("invalid filename")
	}
	path := filepath.Join(s.Dir, filename)
	return os.Remove(path)
}

func extFromHeader(header *multipart.FileHeader) (string, bool) {
	ct := strings.ToLower(header.Header.Get("Content-Type"))
	switch {
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg", true
	case strings.Contains(ct, "png"):
		return ".png", true
	case strings.Contains(ct, "webp"):
		return ".webp", true
	case strings.Contains(ct, "gif"):
		return ".gif", true
	}
	name := strings.ToLower(header.Filename)
	for _, e := range []string{".jpg", ".jpeg", ".png", ".webp", ".gif"} {
		if strings.HasSuffix(name, e) {
			if e == ".jpeg" {
				return ".jpg", true
			}
			return e, true
		}
	}
	return "", false
}

// JoinPublicPath returns a URL path for API responses (forward slashes).
func JoinPublicPath(prefix, filename string) string {
	if filename == "" {
		return ""
	}
	p := strings.TrimSuffix(prefix, "/")
	return fmt.Sprintf("%s/%s", p, filename)
}
