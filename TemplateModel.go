package frame

type Template struct {
	Name      string
	IsLayout  bool
	UseLayout string
}

type TemplateCollection []Template
