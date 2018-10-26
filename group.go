package expect

import (
	"strings"
	"time"
)

type Group struct {
	show    bool
	find    string
	input   string
	timeout time.Duration
}

const (
	ExpectTimeout = 3 * time.Second
)

func NewGroup(find, input string) *Group {
	if find == "" {
		return nil
	}
	if input == "" {
		return nil
	}
	return &Group{
		show:    true,
		find:    find,
		input:   input,
		timeout: ExpectTimeout,
	}
}

func (g *Group) Show(show bool) *Group {
	g.show = show
	return g
}

func (g *Group) Timeout(t time.Duration) *Group {
	g.timeout = t
	return g
}

func (g *Group) Search(s string) string {
	if strings.Contains(s, g.find) {
		return g.input
	}
	return ""
}
