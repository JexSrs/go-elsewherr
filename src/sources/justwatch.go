package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/utils"
	"io"
	"net/http"
	"strings"
)

const (
	graphqlSearchQuery = `
		query GetSearchTitles(
		  $searchTitlesFilter: TitleFilter!,
		  $country: Country!,
		  $language: Language!,
		  $first: Int!,
		  $formatPoster: ImageFormat,
		  $formatOfferIcon: ImageFormat,
		  $profile: PosterProfile,
		  $backdropProfile: BackdropProfile,
		  $filter: OfferFilter!,
		) {
		  popularTitles(
			country: $country
			filter: $searchTitlesFilter
			first: $first
			sortBy: POPULAR
			sortRandomSeed: 0
		  ) {
			edges {
			  node {
				...TitleDetails
				__typename
			  }
			  __typename
			}
			__typename
		  }
		}`
	graphqlDetailsFragment = `
		fragment TitleDetails on MovieOrShow {
		  id
		  objectId
		  objectType
		  content(country: $country, language: $language) {
			title
			fullPath
			originalReleaseYear
			originalReleaseDate
			runtime
			shortDescription
			genres {
			  shortName
			  __typename
			}
			externalIds {
			  imdbId
			  tmdbId
			  __typename
			}
			posterUrl(profile: $profile, format: $formatPoster)
			backdrops(profile: $backdropProfile, format: $formatPoster) {
			  backdropUrl
			  __typename
			}
			ageCertification
			scoring {
			  imdbScore
			  imdbVotes
			  tmdbPopularity
			  tmdbScore
			  tomatoMeter
			  certifiedFresh
			  jwRating
			  __typename
			}
			interactions {
			  likelistAdditions
			  dislikelistAdditions
			  __typename
			}
			__typename
		  }
		  streamingCharts(country: $country) {
			edges {
			  streamingChartInfo {
				rank
				trend
				trendDifference
				daysInTop3
				daysInTop10
				daysInTop100
				daysInTop1000
				topRank
				updatedAt
				__typename
			  }
			  __typename
			}
			__typename
		  }
		  offers(country: $country, platform: WEB, filter: $filter) {
			...TitleOffer
		  }
		  __typename
		}`
	graphqlOfferFragment = `
		fragment TitleOffer on Offer {
		  id
		  monetizationType
		  presentationType
		  retailPrice(language: $language)
		  retailPriceValue
		  currency
		  lastChangeRetailPriceValue
		  type
		  package {
			id
			packageId
			clearName
			technicalName
			icon(profile: S100, format: $formatOfferIcon)
			__typename
		  }
		  standardWebURL
		  elementCount
		  availableTo
		  deeplinkRoku: deeplinkURL(platform: ROKU_OS)
		  subtitleLanguages
		  videoTechnology
		  audioTechnology
		  audioLanguages
		  __typename
		}`
)

type JustWatch struct {
	rootUrl string
	client  *http.Client
}

type JustWatchTitleResponse struct {
	Data struct {
		PopularTitles struct {
			Edges []struct {
				Node struct {
					Offers []struct {
						ID               string `json:"id"`
						MonetizationType string `json:"monetizationType"`
						PresentationType string `json:"presentationType"`
						Package          struct {
							ClearName     string `json:"clearName"`
							TechnicalName string `json:"technicalName"`
						} `json:"package"`
					} `json:"offers"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"popularTitles"`
	} `json:"data"`
}

func NewJustWatch() *JustWatch {
	return &JustWatch{
		rootUrl: "https://apis.justwatch.com/graphql",
		client:  &http.Client{},
	}
}

func (t *JustWatch) GetProvidersFor(entry utils.Entry, country string) ([]string, error) {
	country = strings.ToUpper(country)

	requestBody := map[string]interface{}{
		"operationName": "GetSearchTitles",
		"variables": map[string]interface{}{
			"first": 10,
			"searchTitlesFilter": map[string]string{
				"searchQuery": entry.Title,
			},
			"language":        "en",
			"country":         country,
			"formatPoster":    "JPG",
			"formatOfferIcon": "PNG",
			"profile":         "S718",
			"backdropProfile": "S1920",
			"filter": map[string]bool{
				"bestOnly": true,
			},
		},
		"query": graphqlSearchQuery + graphqlDetailsFragment + graphqlOfferFragment,
	}

	data, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, t.rootUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("justwatch: request error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("justwatch: request error: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		// If no entry in TMDB, return an empty array
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		dt, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("justwatch: request error: %d %s", resp.StatusCode, string(dt))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("justwatch: request error: %v", err)
	}

	// Parse response
	var response JustWatchTitleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("justwatch: request error: %v", err)
	}

	var names []string
	if len(response.Data.PopularTitles.Edges) > 0 {
		edge := response.Data.PopularTitles.Edges[0]

		providers := make(map[string]bool)
		for _, offer := range edge.Node.Offers {
			if offer.MonetizationType == "FLATRATE" {
				providers[offer.Package.ClearName] = true
			}
		}

		for key := range providers {
			names = append(names, key)
		}
	}

	return names, nil
}
