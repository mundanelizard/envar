package main

import (
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	reader, _ := zlib.NewReader(os.Stdin)
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(data))
}
