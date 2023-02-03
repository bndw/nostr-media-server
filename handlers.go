package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/bndw/nostr-media-server/storage"
)

type handlers struct {
	Config Config
	Store  storageProvider
}

// handleWellKnown returns the nostr.json well-known payload.
func (h *handlers) handleWellKnown(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]any{
		"media": map[string]any{
			"apiPath":           h.Config.APIPath,
			"mediaPath":         h.Config.MediaPath,
			"acceptedMimetypes": h.Config.AcceptedMimetypes,
			"contentPolicy": map[string]any{
				"allowAdultContent":   h.Config.AllowAdultContent,
				"allowViolentContent": h.Config.AllowViolentContent,
			},
		},
		"names": h.Config.Names,
	})

	if err != nil {
		log.Printf("failed to marshal well-known: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// handleGetImage fetches a stored image
func (h *handlers) handleGetImage(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		sum  = chi.URLParam(r, "sum")
		name = chi.URLParam(r, "name")
	)

	f, err := h.Store.Get(ctx, sum, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(fileBytes))
	w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
	w.Write(fileBytes)
}

// handleUpload stores the provided media
func (h *handlers) handleUpload(w http.ResponseWriter, r *http.Request) {
	// TODO: Limit upload size to 1 MB, move to config.
	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)

	fileName, fileBytes, err := getMedia(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		contentType = http.DetectContentType(fileBytes)
		accepted    = false
	)
	if len(h.Config.AcceptedMimetypes) == 0 {
		// No explicit accepted mimetypes, allow all.
		accepted = true
	} else {
		for _, mime := range h.Config.AcceptedMimetypes {
			if strings.EqualFold(contentType, mime) {
				accepted = true
				break
			}
		}
	}
	if !accepted {
		msg := fmt.Sprintf("unaccepted content mimetype %q", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var (
		ctx = r.Context()
		sum = fmt.Sprintf("%x", sha256.Sum256(fileBytes))
	)

	relPath, err := h.Store.Save(ctx, bytes.NewReader(fileBytes), storage.Options{
		Filename: fileName,
		Sha256:   sum,
	})
	if err != nil {
		log.Printf("err: store.Save: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	absPath := filepath.Join(h.Config.MediaPath, relPath)

	data, err := json.Marshal(map[string]any{
		"data": map[string]any{
			"link": absPath,
		},
		"success": true,
		"status":  200,
	})
	if err != nil {
		log.Printf("failed to marshal upload response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getMedia(r *http.Request) (string, []byte, error) {
	err := r.ParseMultipartForm(1 * 1024 * 1024) // 1 MB in memory
	if err != nil {
		return "", nil, err
	}

	fileName := r.Form.Get("filename")
	if fileName == "" {
		return "", nil, fmt.Errorf("must provide filename field")
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil, err
	}

	return fileName, fileBytes, nil
}
