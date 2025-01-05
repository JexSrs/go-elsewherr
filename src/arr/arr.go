package arr

import (
	"github.com/JexSrs/go-elsewherr/src/sources"
	"github.com/JexSrs/go-elsewherr/src/utils"
)

type Arr interface {
	CreateTag(tag string) (*utils.Tag, error)
	GetAllTags() ([]utils.Tag, error)
	GetEntries() ([]utils.Entry, error)
	UpdateEntryTags(entry utils.Entry) error
	GetSource() sources.Source
}
