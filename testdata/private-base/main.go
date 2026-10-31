// Command private-base is the ko fixture of a PR build (push: false) on a
// base image that is a private ghcr.io package (self-check.yml,
// image-private-base): .ko.yaml takes the image image-publish pushes.
package main

import "fmt"

func main() {
	fmt.Println("image-build's private-base fixture")
}
