// Package commit
// Represents a commit in envi history.
package database

import (
	"fmt"
)

type Commit struct {
	id      string
	treeId  string
	aut     *Author
	message string
	parent  string
}

func NewCommit(parent, treeId string, aut *Author, message string) *Commit {
	return &Commit{
		treeId:  treeId,
		aut:     aut,
		message: message,
		parent:  parent,
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

	var parent string

	if len(c.parent) != 0 {
		parent = fmt.Sprintf("parent %s\n", c.parent)
	}

	return fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s", c.treeId, parent, a, a, c.message)
}

func (c *Commit) Type() string {
	return "commit"
}
