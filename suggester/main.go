package suggester

import "github.com/c-bata/go-prompt"

var Suggester []prompt.Suggest

type Completer struct {
}

func Create() *[]prompt.Suggest {
	Suggester = []prompt.Suggest{
		{Text: "quit", Description: "Exit from shellize"},
	}

	return &Suggester
}

func AddSuggest(cmd string) {
	find := false
	for _, v := range Suggester {
		if v.Text == cmd {
			find = true
		}
	}
	if !find {
		s := new(prompt.Suggest)
		s.Text = cmd
		s.Description = ""
		Suggester = append(Suggester, *s)
	}
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return Suggester
	}
	return []prompt.Suggest{}
}
