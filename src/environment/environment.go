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

	TMDBKey string

	RadarrUrl string
	RadarrKey string

	SonarrUrl string
	SonarrKey string

	TagPrefix string
	Country   string
}

var Env = Environment{
	Debug: false,

	TagPrefix: "go-elsewherr-",
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

	radarrUrl, isSet := os.LookupEnv("RADARR_URL")
	if isSet {
		Env.RadarrUrl = radarrUrl
	}

	radarrKey, isSet := os.LookupEnv("RADARR_KEY")
	if isSet {
		Env.RadarrKey = radarrKey
	}

	sonarrUrl, isSet := os.LookupEnv("SONARR_URL")
	if isSet {
		Env.SonarrUrl = sonarrUrl
	}

	sonarrKey, isSet := os.LookupEnv("SONARR_KEY")
	if isSet {
		Env.SonarrKey = sonarrKey
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
