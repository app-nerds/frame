package webapp

type Template struct {
	Name      string
	IsLayout  bool
	UseLayout string
}

type TemplateCollection []Template
