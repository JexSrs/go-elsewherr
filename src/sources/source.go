package sources

type Source interface {
	GetProvidersFor(entryId any, country string) ([]string, error)
}
