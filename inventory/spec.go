package inventory

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

type SpecMap map[string]reflect.Type

func ParseSpec(o *Object, node yaml.Node, specTypes SpecMap) (any, error) {
	specType, ok := specTypes[o.Class.Label]
	if !ok {
		return nil, nil
	}

	spec := reflect.New(specType).Interface()
	err := node.Decode(spec)
	if err != nil {
		return nil, fmt.Errorf("error decoding spec into type %s: %s", reflect.TypeOf(node).Name(), err)
	}

	return spec, nil
}
