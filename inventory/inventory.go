package inventory

import (
	"iter"

	"github.com/moznion/go-optional"
)

type Id uint32

type Inventory struct {
	classes map[Id]*Class
	objects map[Id]*Object
}

func NewInventory() *Inventory {
	return &Inventory{
		map[Id]*Class{},
		map[Id]*Object{},
	}
}

func (inventory *Inventory) GetClass(id Id) optional.Option[*Class] {
	class, ok := inventory.classes[id]
	if !ok {
		return optional.None[*Class]()
	}

	return optional.Some(class)
}

func (inventory *Inventory) GetObject(id Id) optional.Option[*Object] {
	object, ok := inventory.objects[id]
	if !ok {
		return optional.None[*Object]()
	}

	return optional.Some(object)
}

func (inventory *Inventory) Classes() iter.Seq2[Id, *Class] {
	return func(yield func(Id, *Class) bool) {
		for k, v := range inventory.classes {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (inventory *Inventory) Objects() iter.Seq2[Id, *Object] {
	return func(yield func(Id, *Object) bool) {
		for k, v := range inventory.objects {
			if !yield(k, v) {
				return
			}
		}
	}
}
