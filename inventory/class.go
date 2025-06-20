package inventory

type Class[C, P comparable] struct {
	Id               C
	Label            string
	Ports            map[P]*Port[P]
	Connections      []*Connection[P]
	PathConstruction *PathConstruction[P]
}

func NewClass[C, P comparable](id C, label string) *Class[C, P] {
	return &Class[C, P]{
		id,
		label,
		map[P]*Port[P]{},
		[]*Connection[P]{},
		nil,
	}
}

type Port[P comparable] struct {
	Id    P      `yaml:"id"`
	Label string `yaml:"label"`
}

type Connection[P comparable] struct {
	From          P    `yaml:"from"`
	To            P    `yaml:"to"`
	Bidirectional bool `yaml:"bidirectional"`
}

type PathConstruction[P comparable] struct {
	Start *Port[P]
	End   *Port[P]
}
