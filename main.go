// Package main provides the database
package main

func main() {
	a := App{}
	a.Initialize()
	a.Run(":8080")
}
