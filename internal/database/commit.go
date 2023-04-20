package database

import (
	"errors"
	"fmt"
	"strings"
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

func NewCommitFromByteArray(id string, data []byte) (*Commit, error) {
	content := string(data)
	ok := strings.HasPrefix(content, "commit")
	if !ok {
		return nil, errors.New("object is not a commit")
	}

	chunks := strings.Split(content, "\x00")
	content = strings.Join(chunks[1:], "")
	chunks = strings.Split(content, "\n")

	if len(chunks) != 5 && len(chunks) != 4 {
		return nil, errors.New(fmt.Sprintf("invalid chucks length of %d", len(chunks)))
	}

	treeId := strings.Split(chunks[0], " ")[1]
	offset := 0
	var parent string

	if len(chunks) == 5 {
		offset = 1
		parent = strings.Split(chunks[1], " ")[1]
	}

	fmt.Println(chunks)

	trimmedAuthor := strings.Join(strings.Split(chunks[1+offset], " ")[1:], " ")
	author, err := NewAuthorFromByteArray(trimmedAuthor)
	if err != nil {
		return nil, err
	}

	message := strings.Split(chunks[3+offset], " ")[1]

	commit := &Commit{
		id:      id,
		aut:     author,
		message: message,
		parent:  parent,
		treeId:  treeId,
	}

	return commit, nil
}

func (c *Commit) TreeId() string {
	return c.treeId
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
