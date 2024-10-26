package arr

import "github.com/JexSrs/go-elsewherr/src/sources"

type Arr interface {
	CreateTag(tag string) (*Tag, error)
	GetAllTags() ([]Tag, error)
	GetEntries() ([]Entry, error)
	UpdateEntryTags(entry Entry) error
	GetSource() sources.Source
}

type Tag struct {
	ID   int
	Name string
}

type Entry struct {
	ID     int
	Title  string
	IMDBID string
	TMDBID int
	Tags   []int

	Original map[string]any
}

func Get[T any](m map[string]interface{}, key string) T {
	dt, _ := m[key]
	return dt.(T)
}

func ToIntArray(arr []interface{}) []int {
	intArray := make([]int, len(arr))
	for i, v := range arr {
		intArray[i] = int(v.(float64))
	}

	return intArray
}
