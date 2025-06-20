package inventory

import (
	"fmt"
	"io"
	"os"

	"github.com/moznion/go-optional"
	"gopkg.in/yaml.v3"
)

type Model[O, C, P comparable] struct {
	Classes []*ClassModel[C, P]     `yaml:"classes"`
	Objects []*ObjectModel[O, C, P] `yaml:"objects"`
}

type ClassModel[C, P comparable] struct {
	Id               C                         `yaml:"id"`
	Label            string                    `yaml:"label"`
	Ports            []*Port[P]                `yaml:"ports"`
	Connections      []*Connection[P]          `yaml:"connections"`
	PathConstruction *PathConstructionModel[P] `yaml:"pathConstruction"`
}

func (model *ClassModel[O, P]) ToClass() (*Class[O, P], error) {
	class := NewClass[O, P](model.Id, model.Label)

	for _, port := range model.Ports {
		if _, ok := class.Ports[port.Id]; ok {
			return nil, fmt.Errorf("class parsing error: duplicate port %s (ID %v) in class %s (ID %v)", port.Label, port.Id, model.Label, model.Id)
		}

		class.Ports[port.Id] = port
	}

	var pathConstruction *PathConstruction[P]
	if model.PathConstruction != nil {
		var err error
		pathConstruction, err = model.PathConstruction.ToPathConstruction(class.Ports)
		if err != nil {
			return nil, fmt.Errorf("class parsing error: %s in class %s (ID %v)", err, class.Label, class.Id)
		}
	}

	class.PathConstruction = pathConstruction
	class.Connections = model.Connections
	return class, nil
}

type PathConstructionModel[P comparable] struct {
	Start P `yaml:"start"`
	End   P `yaml:"end"`
}

func (model *PathConstructionModel[P]) ToPathConstruction(ports map[P]*Port[P]) (*PathConstruction[P], error) {
	start, ok := ports[model.Start]
	if !ok {
		return nil, fmt.Errorf("path construction parsing error: start port ref %v does not exist", model.Start)
	}

	end, ok := ports[model.End]
	if !ok {
		return nil, fmt.Errorf("path construction parsing error: end port ref %v does not exist", model.End)
	}

	return &PathConstruction[P]{
		start,
		end,
	}, nil
}

type ObjectModel[O, C, P comparable] struct {
	Id       O         `yaml:"id"`
	Label    string    `yaml:"label"`
	ClassRef C         `yaml:"class"`
	Spec     yaml.Node `yaml:"spec"`
}

func (model *ObjectModel[O, C, P]) ToObject(classes map[C]*Class[C, P], specTypes SpecMap) (*Object[O, C, P], error) {
	object := NewObject[O, C, P](model.Id, model.Label)

	class, ok := classes[model.ClassRef]
	if !ok {
		return nil, fmt.Errorf("model parsing error: class ref %v in object %s (ID %v) not found", model.ClassRef, model.Label, model.Id)
	}
	object.Class = class

	if len(specTypes) > 0 {
		spec, err := ParseSpec(object, model.Spec, specTypes)
		if err != nil {
			return nil, fmt.Errorf("model parsing error: cannot parse spec of object %s: %s", model.Label, err)
		}

		if spec != nil {
			object.Spec = optional.Some(spec)
		}
	}

	return object, nil
}

func (model *Model[O, C, P]) ToInventory(specTypes SpecMap) (*Inventory[O, C, P], error) {
	inv := NewInventory[O, C, P]()

	for _, classModel := range model.Classes {
		if _, ok := inv.classes[classModel.Id]; ok {
			return nil, fmt.Errorf("parsing error: duplicate class %s (ID %v)", classModel.Label, classModel.Id)
		}

		class, err := classModel.ToClass()
		if err != nil {
			return nil, err
		}

		inv.classes[class.Id] = class
	}

	for _, objectModel := range model.Objects {
		object, err := objectModel.ToObject(inv.classes, specTypes)
		if err != nil {
			return nil, err
		}

		inv.objects[object.Id] = object
	}

	return inv, nil
}

func Parse[O, C, P comparable](r io.ReadCloser) (*Inventory[O, C, P], error) {
	return ParseWithSpec[O, C, P](r, nil)
}

func ParseWithSpec[O, C, P comparable](r io.ReadCloser, specTypes SpecMap) (*Inventory[O, C, P], error) {
	var model *Model[O, C, P]
	if err := yaml.NewDecoder(r).Decode(&model); err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	return model.ToInventory(specTypes)
}

func ParseFile[O, C, P comparable](filename string) (*Inventory[O, C, P], error) {
	return ParseFileWithSpec[O, C, P](filename, nil)
}

func ParseFileWithSpec[O, C, P comparable](filename string, specTypes SpecMap) (*Inventory[O, C, P], error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", filename, err)
	}

	return ParseWithSpec[O, C, P](f, specTypes)
}
