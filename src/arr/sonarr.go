package arr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/environment"
	"github.com/JexSrs/go-elsewherr/src/sources"
	"io"
	"net/http"
)

type Sonarr struct {
	URL    string
	APIKey string
}

type SonarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func NewSonarr(url, apiKey string) *Sonarr {
	return &Sonarr{url, apiKey}
}

func (r *Sonarr) Do(method, path string, data map[string]any) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", r.URL, path)

	var req *http.Request
	var err error
	if data != nil {
		jData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", r.APIKey)

	client := &http.Client{}
	return client.Do(req)
}

func (r *Sonarr) CreateTag(tag string) (*Tag, error) {
	resp, err := r.Do(http.MethodPost, "/api/v3/tag", map[string]any{"label": tag, "id": 0})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get response: %d %s", resp.StatusCode, string(dt))
	}

	var t2 RadarrTag
	if err := json.NewDecoder(resp.Body).Decode(&t2); err != nil {
		return nil, err
	}

	return &Tag{
		ID:   t2.ID,
		Name: t2.Label,
	}, nil
}

func (r *Sonarr) GetAllTags() ([]Tag, error) {
	resp, err := r.Do(http.MethodGet, "/api/v3/tag", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get response: %d %s", resp.StatusCode, string(dt))
	}

	var tags []RadarrTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	res := make([]Tag, len(tags))
	for i, t := range tags {
		res[i] = Tag{
			ID:   t.ID,
			Name: t.Label,
		}
	}
	return res, nil
}

func (r *Sonarr) GetEntries() ([]Entry, error) {
	resp, err := r.Do(http.MethodGet, "/api/v3/series", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get response: %d %s", resp.StatusCode, string(dt))
	}

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	res := make([]Entry, len(items))
	for i, item := range items {
		res[i] = Entry{
			ID:       int(Get[float64](item, "id")),
			Title:    Get[string](item, "title"),
			IMDBID:   Get[string](item, "imdbId"),
			TMDBID:   int(Get[float64](item, "tmdbId")),
			Tags:     ToIntArray(Get[[]interface{}](item, "tags")),
			Original: item,
		}
	}
	return res, nil
}

func (r *Sonarr) UpdateEntryTags(entry Entry) error {
	entry.Original["tags"] = entry.Tags

	resp, err := r.Do(http.MethodPut, fmt.Sprintf("/api/v3/series/%d", entry.ID), entry.Original)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		dt, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get response: %d %s", resp.StatusCode, string(dt))
	}
	return nil
}

func (r *Sonarr) GetSource() sources.Source {
	return sources.NewTMDB(environment.Env.TMDBKey, "tv")
}
