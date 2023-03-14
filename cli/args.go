package cli

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

type ActionArgs struct {
	mappings map[string]interface{}
}

// TODO :=> Make this function private to the package.

func NewActionArgs() *ActionArgs {
	return &ActionArgs{
		mappings: make(map[string]interface{}),
	}
}

// TODO :=> export interfaces for ActionArguments instead of the type.

func (args *ActionArgs) Set(key string, ptr interface{}) {
	args.mappings[key] = ptr
}

func (args *ActionArgs) Get(key string) any {
	return args.mappings[key]
}

func (args *ActionArgs) GetString(key string) (string, error) {
	ptr, ok := args.mappings[key].(*string)

	if !ok {
		return "", errors.New(fmt.Sprintf("can't retrieve string value of '%s'", key))
	}

	return *ptr, nil
}

func (args *ActionArgs) GetBool(key string) (bool, error) {
	ptr, ok := args.mappings[key].(*bool)

	if !ok {
		return false, errors.New(fmt.Sprintf("can't retrieve boolean value of '%s'", key))
	}

	return *ptr, nil
}

func (args *ActionArgs) GetInt(key string) (int, error) {
	ptr, ok := args.mappings[key].(*int)

	if !ok {
		return 0, errors.New(fmt.Sprintf("can't retrieve int value of '%s'", key))
	}

	return *ptr, nil
}

func (args *ActionArgs) String() string {
	buf := new(bytes.Buffer)

	for key, val := range args.mappings {
		_, err := fmt.Fprintf(buf, "Key: %s - Value %v\n", key, val)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return buf.String()
}
