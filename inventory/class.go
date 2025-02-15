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
	Id    Id     `json:"id" yaml:"id"`
	Label string `json:"label" yaml:"label"`
}

type Connection struct {
	From          Id   `json:"from" yaml:"from"`
	To            Id   `json:"to" yaml:"to"`
	Bidirectional bool `json:"bidirectional" yaml:"bidirectional"`
}

type PathConstruction struct {
	Start *Port
	End   *Port
}
