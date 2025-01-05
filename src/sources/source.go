package sources

import "github.com/JexSrs/go-elsewherr/src/utils"

type Source interface {
	GetProvidersFor(entry utils.Entry, country string) ([]string, error)
}
