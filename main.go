// Package main provides the database
package main

// Main begins the application
func Main() {
	a := App{}
	a.Initialize("explorer")
	a.Run(":8080")
}
