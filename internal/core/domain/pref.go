package domain

import "html/template"

type Pref struct {
	Name    string
	Display template.HTML
	Type    string
	Width   int
	Value   string
}
