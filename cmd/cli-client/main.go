package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/chiselwright/go-jekyllsitetagger"
)

func main() {
	kong.Parse(&jekyllsitetagger.CLI)

	fmt.Println(jekyllsitetagger.CLI.Source)

	jekyllsitetagger.GenerateTagFiles(
		jekyllsitetagger.CLI.Source,
		jekyllsitetagger.CLI.OutputTo,
	)
}
