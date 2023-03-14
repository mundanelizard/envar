package cli

import (
	"errors"
	"flag"
)

type Flagger interface {
	Attach(set *flag.FlagSet) any
	GetID() string
	Minify() (string, bool)
	Validate(ptr any) error
}

type Flag struct {
	Name     string
	Usage    string
	Required bool
	Shrink   bool
}

func (f *Flag) GetID() string {
	return f.Name
}

func (f *Flag) Minify() (string, bool) {
	if !f.Shrink {
		return "", false
	}

	return string(f.Name[0]), true
}

type IntFlag struct {
	Flag
	Value int
}

func (f *IntFlag) Attach(set *flag.FlagSet) any {
	return set.Int(f.Name, f.Value, f.Usage)
}

func (f *IntFlag) Validate(ptr any) error {
	_, ok := ptr.(*int)

	if !ok {
		return errors.New(f.Usage)
	}

	return nil
}

type StringFlag struct {
	Flag
	Value string
}

func (f *StringFlag) Attach(set *flag.FlagSet) any {
	return set.String(f.Name, f.Value, f.Usage)
}

func (f *StringFlag) Validate(ptr any) error {
	value, ok := ptr.(*string)

	if !ok {
		return errors.New(f.Usage)
	}

	if f.Required && len(*value) == 0 {
		return errors.New(f.Usage)
	}

	return nil
}

type BoolFlag struct {
	Flag
	Value bool
}

func (f *BoolFlag) Attach(set *flag.FlagSet) any {
	return set.Bool(f.Name, f.Value, f.Usage)
}

func (f *BoolFlag) Validate(ptr any) error {
	_, ok := ptr.(*bool)

	if !ok {
		return errors.New(f.Usage)
	}

	return nil
}
