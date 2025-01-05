package main

import (
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/arr"
	. "github.com/JexSrs/go-elsewherr/src/environment"
	"github.com/JexSrs/go-elsewherr/src/utils"
	"log"
	"strings"
)

func main() {
	fmt.Println()
	fmt.Println("Starting go-elsewherr...")

	if len(Env.RadarrUrl) != 0 {
		fmt.Println("\nRadarr configuration found")
		radarr := arr.NewRadarr(Env.RadarrUrl, Env.RadarrKey)

		if err := Sync(radarr); err != nil {
			log.Fatal(err)
		}
	}

	if len(Env.SonarrUrl) != 0 {
		fmt.Println("\nSonarr configuration found")
		sonarr := arr.NewSonarr(Env.SonarrUrl, Env.SonarrKey)

		if err := Sync(sonarr); err != nil {
			log.Fatal(err)
		}
	}

	if len(Env.RadarrUrl) == 0 && len(Env.SonarrUrl) == 0 {
		fmt.Println("No Arr configuration was not found, exiting...")
	}
}

func Sync(app arr.Arr) error {
	source := app.GetSource()

	debug("Sync operation started")

	debug("Retrieving arr tags...")
	arrTags, err := app.GetAllTags()
	if err != nil {
		return fmt.Errorf("failed to retrieve tags: %w", err)
	}
	debug("Found %d tags", len(arrTags))

	debug("Retrieving entries...")
	entries, err := app.GetEntries()
	if err != nil {
		return fmt.Errorf("failed to retrieve entries: %w", err)
	}
	debug("Found %d entries", len(entries))

	// TODO: Add go routines for parallel requests
	for _, entry := range entries {
		debug("Processing entry: %s (TMDB: %d, IMDB: %s)", entry.Title, entry.TMDBID, entry.IMDBID)

		debug("- Requesting providers for country %s", Env.Country)
		providers, err := source.GetProvidersFor(entry, Env.Country)
		if err != nil {
			return fmt.Errorf("failed to retrieve providers: %w", err)
		}

		debug("- Found %d providers", len(providers))

		// Remove tags
		debug("- Removing old tags")
		entry.Tags = utils.Filter(entry.Tags, func(i int) bool {
			arrTagIdx := utils.FindIndex(arrTags, func(tag utils.Tag) bool { return tag.ID == i })
			return !strings.HasPrefix(arrTags[arrTagIdx].Name, Env.TagPrefix)
		})

		for _, provider := range providers {
			existingTag := utils.FindIndex(arrTags, func(tag utils.Tag) bool {
				return tag.Name == Env.TagPrefix+utils.CleanString(provider)
			})

			var tag utils.Tag
			if existingTag == -1 {
				debug("- Tag %s not found, creating new one", Env.TagPrefix+utils.CleanString(provider))
				t, err := app.CreateTag(Env.TagPrefix + utils.CleanString(provider))
				if err != nil {
					return fmt.Errorf("failed to create tag: %w", err)
				}

				arrTags = append(arrTags, *t)
				tag = *t
			} else {
				tag = arrTags[existingTag]
			}

			debug("- Appending tag %s", provider)
			entry.Tags = append(entry.Tags, tag.ID)
		}

		debug("- Updating entry...")
		if err := app.UpdateEntryTags(entry); err != nil {
			return fmt.Errorf("failed to update entry: %w", err)
		}
	}

	debug("Sync operation finished.")
	return nil
}

func debug(format string, args ...any) {
	if Env.Debug {
		fmt.Printf(format+"\n", args...)
	}
}
