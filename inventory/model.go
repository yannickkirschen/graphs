package inventory

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Model struct {
	Classes []*ClassModel  `yaml:"classes"`
	Objects []*ObjectModel `yaml:"objects"`
}

type ClassModel struct {
	Id               Id                     `yaml:"id"`
	Label            string                 `yaml:"label"`
	Ports            []*Port                `yaml:"ports"`
	Connections      []*Connection          `yaml:"connections"`
	PathConstruction *PathConstructionModel `yaml:"pathConstruction"`
}

func (model *ClassModel) ToClass() (*Class, error) {
	class := NewClass(model.Id, model.Label)

	for _, port := range model.Ports {
		if _, ok := class.Ports[port.Id]; ok {
			return nil, fmt.Errorf("class parsing error: duplicate port %s (ID %d) in class %s (ID %d)", port.Label, port.Id, model.Label, model.Id)
		}

		class.Ports[port.Id] = port
	}

	var pathConstruction *PathConstruction
	if model.PathConstruction != nil {
		var err error
		pathConstruction, err = model.PathConstruction.ToPathConstruction(class.Ports)
		if err != nil {
			return nil, fmt.Errorf("class parsing error: %s in class %s (ID %d)", err, class.Label, class.Id)
		}
	}

	class.PathConstruction = pathConstruction
	class.Connections = model.Connections
	return class, nil
}

type PathConstructionModel struct {
	Start Id `yaml:"start"`
	End   Id `yaml:"end"`
}

func (model *PathConstructionModel) ToPathConstruction(ports map[Id]*Port) (*PathConstruction, error) {
	start, ok := ports[model.Start]
	if !ok {
		return nil, fmt.Errorf("path construction parsing error: start port ref %d does not exist", model.Start)
	}

	end, ok := ports[model.End]
	if !ok {
		return nil, fmt.Errorf("path construction parsing error: end port ref %d does not exist", model.End)
	}

	return &PathConstruction{
		start,
		end,
	}, nil
}

type ObjectModel struct {
	Id       Id     `yaml:"id"`
	Label    string `yaml:"label"`
	ClassRef Id     `yaml:"class"`
}

func (model *ObjectModel) ToObject(classes map[Id]*Class) (*Object, error) {
	object := NewObject(model.Id, model.Label)

	class, ok := classes[model.ClassRef]
	if !ok {
		return nil, fmt.Errorf("model parsing error: class ref %d in object %s (ID %d) not found", model.ClassRef, model.Label, model.Id)
	}

	object.Class = class
	return object, nil
}

func (model *Model) ToInventory() (*Inventory, error) {
	inv := NewInventory()

	for _, classModel := range model.Classes {
		if _, ok := inv.classes[classModel.Id]; ok {
			return nil, fmt.Errorf("parsing error: duplicate class %s (ID %d)", classModel.Label, classModel.Id)
		}

		class, err := classModel.ToClass()
		if err != nil {
			return nil, err
		}

		inv.classes[class.Id] = class
	}

	for _, objectModel := range model.Objects {
		object, err := objectModel.ToObject(inv.classes)
		if err != nil {
			return nil, err
		}

		inv.objects[object.Id] = object
	}

	return inv, nil
}

func Parse(r io.ReadCloser) (*Inventory, error) {
	var model *Model
	if err := yaml.NewDecoder(r).Decode(&model); err != nil {
		return nil, fmt.Errorf("error parsing input: %s", err)
	}

	return model.ToInventory()
}

func ParseFile(filename string) (*Inventory, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %s", filename, err)
	}

	return Parse(f)
}
