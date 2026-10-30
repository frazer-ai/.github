// Command image-vuln is the ko fixture that image-build's vulnerability gate
// must fail (self-check.yml): ko copies kodata/ into the image, and with it
// the CycloneDX fragment that declares Unbound 1.26.0.
package main

import "fmt"

func main() {
	fmt.Println("image-build's must-fail fixture")
}
