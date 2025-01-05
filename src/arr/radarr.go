package arr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/environment"
	"github.com/JexSrs/go-elsewherr/src/sources"
	"github.com/JexSrs/go-elsewherr/src/utils"
	"io"
	"net/http"
)

type Radarr struct {
	URL    string
	APIKey string
	client *http.Client
}

type RadarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func NewRadarr(url, apiKey string) *Radarr {
	return &Radarr{
		url,
		apiKey,
		&http.Client{},
	}
}

func (r *Radarr) Do(method, path string, data map[string]any) (*http.Response, error) {
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

	return r.client.Do(req)
}

func (r *Radarr) CreateTag(tag string) (*utils.Tag, error) {
	resp, err := r.Do(http.MethodPost, "/api/v3/tag", map[string]any{"label": tag, "id": 0})
	if err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("radarr: request error: %d %s", resp.StatusCode, string(dt))
	}

	var t2 RadarrTag
	if err := json.NewDecoder(resp.Body).Decode(&t2); err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}

	return &utils.Tag{
		ID:   t2.ID,
		Name: t2.Label,
	}, nil
}

func (r *Radarr) GetAllTags() ([]utils.Tag, error) {
	resp, err := r.Do(http.MethodGet, "/api/v3/tag", nil)
	if err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("radarr: request error: %d %s", resp.StatusCode, string(dt))
	}

	var tags []RadarrTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}

	res := make([]utils.Tag, len(tags))
	for i, t := range tags {
		res[i] = utils.Tag{
			ID:   t.ID,
			Name: t.Label,
		}
	}
	return res, nil
}

func (r *Radarr) GetEntries() ([]utils.Entry, error) {
	resp, err := r.Do(http.MethodGet, "/api/v3/movie", nil)
	if err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("radarr: request error: %d %s", resp.StatusCode, string(dt))
	}

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("radarr: request error: %v", err)
	}

	res := make([]utils.Entry, len(items))
	for i, item := range items {
		res[i] = utils.Entry{
			ID:       int(utils.Get[float64](item, "id")),
			Title:    utils.Get[string](item, "title"),
			IMDBID:   utils.Get[string](item, "imdbId"),
			TMDBID:   int(utils.Get[float64](item, "tmdbId")),
			Tags:     utils.ToIntArray(utils.Get[[]interface{}](item, "tags")),
			Original: item,
		}
	}
	return res, nil
}

func (r *Radarr) UpdateEntryTags(entry utils.Entry) error {
	entry.Original["tags"] = entry.Tags

	resp, err := r.Do(http.MethodPut, fmt.Sprintf("/api/v3/movie/%d", entry.ID), entry.Original)
	if err != nil {
		return fmt.Errorf("radarr: request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		dt, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("radarr: request error: %d %s", resp.StatusCode, string(dt))
	}
	return nil
}

func (r *Radarr) GetSource() sources.Source {
	if environment.Env.RadarrSource == "tmdb" {
		return sources.NewTMDB(environment.Env.TMDBKey, "movie")
	} else if environment.Env.RadarrSource == "justwatch" {
		return sources.NewJustWatch()
	}

	panic(fmt.Errorf("unknown sonarr source %s", environment.Env.RadarrSource))
}
