// Package main provides the database
package main

func main() {
	a := App{}
	a.Initialize("explorer")
	a.Run(":8080")
}
