package main

import (
	"fmt"

	"bitbucket.org/jebo87/makako-api/store"
)

func main() {
	fmt.Println("Loading makako API server...")
	store.GetAdTitles(0, 0)
}
