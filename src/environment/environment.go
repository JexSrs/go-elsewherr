package environment

import (
	_ "github.com/joho/godotenv/autoload"
	"strings"
)
import (
	"os"
)

type Environment struct {
	Debug bool

	TMDBKey      string
	JustWatchKey string

	RadarrUrl    string
	RadarrKey    string
	RadarrSource string

	SonarrUrl    string
	SonarrKey    string
	SonarrSource string

	TagPrefix string
	Country   string
}

var Env = Environment{
	Debug: false,

	TagPrefix: "go-",

	RadarrSource: "tmdb",

	SonarrSource: "tmdb",
}

func init() {
	debug, isSet := os.LookupEnv("DEBUG")
	if isSet {
		Env.Debug = strings.ToUpper(debug) == "TRUE" || strings.ToUpper(debug) == "1"
	}

	tmdbKey, isSet := os.LookupEnv("TMDB_KEY")
	if isSet {
		Env.TMDBKey = tmdbKey
	}

	justWatchKey, isSet := os.LookupEnv("JUSTWATCH_KEY")
	if isSet {
		Env.JustWatchKey = justWatchKey
	}

	radarrUrl, isSet := os.LookupEnv("RADARR_URL")
	if isSet {
		Env.RadarrUrl = radarrUrl
	}

	radarrKey, isSet := os.LookupEnv("RADARR_KEY")
	if isSet {
		Env.RadarrKey = radarrKey
	}

	radarrSource, isSet := os.LookupEnv("RADARR_SOURCE")
	if isSet {
		Env.RadarrSource = radarrSource
	}

	sonarrUrl, isSet := os.LookupEnv("SONARR_URL")
	if isSet {
		Env.SonarrUrl = sonarrUrl
	}

	sonarrKey, isSet := os.LookupEnv("SONARR_KEY")
	if isSet {
		Env.SonarrKey = sonarrKey
	}

	sonarrSource, isSet := os.LookupEnv("SONARR_SOURCE")
	if isSet {
		Env.SonarrSource = sonarrSource
	}

	tagPrefix, isSet := os.LookupEnv("TAG_PREFIX")
	if isSet {
		Env.TagPrefix = tagPrefix
	}

	country, isSet := os.LookupEnv("COUNTRY")
	if isSet {
		Env.Country = country
	}
}
