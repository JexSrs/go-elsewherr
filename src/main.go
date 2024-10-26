package main

import (
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/arr"
	. "github.com/JexSrs/go-elsewherr/src/environment"
	"strings"
)

func main() {
	fmt.Println()
	fmt.Println("Starting go-elsewherr...")

	if len(Env.RadarrUrl) != 0 {
		fmt.Println("Radarr configuration found")
		radarr := arr.NewRadarr(Env.RadarrUrl, Env.RadarrKey)

		if err := Sync(radarr); err != nil {
			fmt.Printf("Failed with error: %v\n", err)
		}
	}

	if len(Env.SonarrUrl) != 0 {
		fmt.Println("Sonarr configuration found")
		sonarr := arr.NewSonarr(Env.SonarrUrl, Env.SonarrKey)

		if err := Sync(sonarr); err != nil {
			fmt.Printf("Failed with error: %v\n", err)
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

	for _, entry := range entries {
		debug("Processing entry: %s (TMDB: %d, IMDB: %s)", entry.Title, entry.TMDBID, entry.IMDBID)

		debug("- Requesting providers for country %s", Env.Country)
		providers, err := source.GetProvidersFor(entry.TMDBID, Env.Country)
		if err != nil {
			return fmt.Errorf("failed to retrieve providers: %w", err)
		}

		debug("- Found %d providers", len(providers))

		// Remove tags
		debug("- Removing old tags")
		entry.Tags = Filter(entry.Tags, func(i int) bool {
			arrTagIdx := FindIndex(arrTags, func(tag arr.Tag) bool { return tag.ID == i })
			return !strings.HasPrefix(arrTags[arrTagIdx].Name, Env.TagPrefix)
		})

		for _, provider := range providers {
			existingTag := FindIndex(arrTags, func(tag arr.Tag) bool {
				return tag.Name == Env.TagPrefix+cleanString(provider)
			})

			var tag arr.Tag
			if existingTag == -1 {
				debug("- Tag %s not found, creating new one", Env.TagPrefix+cleanString(provider))
				t, err := app.CreateTag(Env.TagPrefix + cleanString(provider))
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
