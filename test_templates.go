package main

import (
	"fmt"
	"microservice-bootstrapper/internal/template"
)

func main() {
	engine := template.NewEngine()
	
	// Test getting project template
	projectTemplate, err := engine.GetProjectTemplate()
	if err != nil {
		fmt.Printf("Error getting project template: %v\n", err)
		return
	}
	
	fmt.Printf("Project template files found: %d\n", len(projectTemplate.Files))
	for filename := range projectTemplate.Files {
		fmt.Printf("- %s\n", filename)
	}
}