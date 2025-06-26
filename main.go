package main

import (
	"fmt"
	"log"
	"net/http"

	"blue-website/web"
)

func main() {
	// Generate search index at startup
	fmt.Println("Generating search index...")
	if err := web.GenerateSearchIndex(); err != nil {
		log.Printf("Warning: Failed to generate search index: %v", err)
	} else {
		fmt.Println("Search index generated successfully")
	}

	// Serve static files from public directory
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// File-based routing handler
	router := web.NewRouter("pages")
	http.Handle("/", router)

	fmt.Println("Blue Website server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

 