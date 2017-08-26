// Package main provides the database
package main

func main() {
	a := App{}
	a.Initialize("explorer")
	EnsureTableExists(a.DB)
	a.Run(":8080")
}
