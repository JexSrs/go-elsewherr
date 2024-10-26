package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TMDB struct {
	APIKey string
	Path   string
}

type TmdbResponse struct {
	Results map[string]*TmdbRegionProviders `json:"results"`
}

type TmdbRegionProviders struct {
	Flatrate []TmdbProvider `json:"flatrate"`
}

type TmdbProvider struct {
	Name string `json:"provider_name"`
}

func NewTMDB(apiKey, path string) *TMDB {
	return &TMDB{apiKey, path}
}

func (t *TMDB) GetProvidersFor(entryId any, country string) ([]string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/%s/%d/watch/providers?api_key=%s", t.Path, entryId, t.APIKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		// If no entry in TMDB, return an empty array
		if resp.StatusCode == http.StatusNotFound {
			return make([]string, 0), nil
		}

		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get response: %d %s", resp.StatusCode, string(dt))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the TmdbResponse struct
	var res TmdbResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	data, exists := res.Results[country]
	if !exists {
		return make([]string, 0), nil
	}

	result := make([]string, len(data.Flatrate))
	for i, provider := range data.Flatrate {
		result[i] = provider.Name
	}

	return result, nil
}
