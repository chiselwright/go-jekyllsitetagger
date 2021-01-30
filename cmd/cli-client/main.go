package main

import (
	"github.com/chiselwright/go-jekyllsitetagger"
)

func main() {
	jekyllsitetagger.GenerateTagFiles("testdata", "tags")
}
