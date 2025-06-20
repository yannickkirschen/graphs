package inventory

import (
	"iter"

	"github.com/moznion/go-optional"
)

type Inventory[O, C, P comparable] struct {
	classes map[C]*Class[C, P]
	objects map[O]*Object[O, C, P]
}

func NewInventory[O, C, P comparable]() *Inventory[O, C, P] {
	return &Inventory[O, C, P]{
		map[C]*Class[C, P]{},
		map[O]*Object[O, C, P]{},
	}
}

func (inventory *Inventory[O, C, P]) GetClass(id C) optional.Option[*Class[C, P]] {
	class, ok := inventory.classes[id]
	if !ok {
		return optional.None[*Class[C, P]]()
	}

	return optional.Some(class)
}

func (inventory *Inventory[O, C, P]) GetObject(id O) optional.Option[*Object[O, C, P]] {
	object, ok := inventory.objects[id]
	if !ok {
		return optional.None[*Object[O, C, P]]()
	}

	return optional.Some(object)
}

func (inventory *Inventory[O, C, P]) Classes() iter.Seq2[C, *Class[C, P]] {
	return func(yield func(C, *Class[C, P]) bool) {
		for k, v := range inventory.classes {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (inventory *Inventory[O, C, P]) Objects() iter.Seq2[O, *Object[O, C, P]] {
	return func(yield func(O, *Object[O, C, P]) bool) {
		for k, v := range inventory.objects {
			if !yield(k, v) {
				return
			}
		}
	}
}
