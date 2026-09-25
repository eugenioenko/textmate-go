package main

import "fmt"

/* buildMessage keeps this comment open
across a physical line boundary. */
func buildMessage(name string) string {
	return `hello,
` + name
}

func main() {
	fmt.Println(buildMessage("gopher"))
}
