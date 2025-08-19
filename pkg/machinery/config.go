package machinery

import (
	"github.com/aity-cloud/monty/pkg/config"
	"github.com/aity-cloud/monty/pkg/config/meta"
)

type JsonDocument interface {
	GetJson() []byte
}

func LoadDocuments[T JsonDocument](documents []T) (meta.ObjectList, error) {
	objects := []meta.Object{}
	for _, document := range documents {
		object, err := config.LoadObject(document.GetJson())
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}
