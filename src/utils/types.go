package utils

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
