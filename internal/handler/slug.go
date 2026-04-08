package handler

import (
	petname "github.com/dustinkirkland/golang-petname"
)

func generateSlug() string {
	return petname.Generate(3, "-")
}
