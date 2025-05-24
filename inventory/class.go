package inventory

type Class struct {
	Id               Id
	Label            string
	Ports            map[Id]*Port
	Connections      []*Connection
	PathConstruction *PathConstruction
}

func NewClass(id Id, label string) *Class {
	return &Class{
		id,
		label,
		map[Id]*Port{},
		[]*Connection{},
		nil,
	}
}

type Port struct {
	Id    Id     `yaml:"id"`
	Label string `yaml:"label"`
}

type Connection struct {
	From          Id   `yaml:"from"`
	To            Id   `yaml:"to"`
	Bidirectional bool `yaml:"bidirectional"`
}

type PathConstruction struct {
	Start *Port
	End   *Port
}
