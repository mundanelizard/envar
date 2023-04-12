// Package commit
// Represents a commit in envi history.
package commit

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/author"
)

const (
	TYPE = "commit"
)

type Commit struct {
	id      string
	treeId  string
	aut     *author.Author
	message string
}

func New(treeId string, aut *author.Author, message string) *Commit {
	return &Commit{
		treeId:  treeId,
		aut:     aut,
		message: message,
	}
}

func (c *Commit) Id() string {
	return c.id
}

func (c *Commit) SetId(id string) {
	c.id = id
}

func (c *Commit) String() string {
	a := c.aut.String()
	return fmt.Sprintf("tree %s\nauthor %s\ncommitter %s\n\n%s", c.treeId, a, a, c.message)
}

func (c *Commit) Type() string {
	return TYPE
}
