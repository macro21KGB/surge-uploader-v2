package main

import (
	"log"
	"os"

	"surge-uploader/m/v2/internal/surge"
	"surge-uploader/m/v2/internal/templates"
)

func commandHelp() {
	println("Usage: surge-uploader <command>")
	println("Commands:")
	println("  help\t\tShow this help message")
	println("  list\t\tList all templates available in the surge-templates directory")
	println(" <template>\tRender the specified template with the list of websites from Surge and create an index.html file")

	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "--help" {
		commandHelp()
	}

	command := os.Args[1]

	switch command {
	case "list":
		templateDirectory := templates.GetDirectory()
		if templateDirectory == "" {
			log.Fatalf("Error getting templates directory")
		}

		println("Templates directory: " + templateDirectory)
		println("Available templates:")
		tpls := templates.GetAll()

		for _, tpl := range tpls {
			println(tpl.Name)
		}

		os.Exit(0)
	}

	// Verify and Pick template
	selectedTemplateName := command

	tpls := templates.GetAll()
	templatePath := ""

	for _, tpl := range tpls {
		if tpl.Name == selectedTemplateName {
			templatePath = tpl.Path
			break
		}
	}

	if templatePath == "" {
		log.Fatalf("Template %s not found. Use 'surge-uploader list' to see available templates.", selectedTemplateName)
	}

	// Get websites from Surge
	websites, err := surge.GetWebsites()
	if err != nil {
		log.Fatalf("Error getting websites from Surge: %v", err)
	}

	renderedHtml, err := templates.Render(templatePath, websites)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	println("Rendering template to index.html")
	os.WriteFile("index.html", []byte(renderedHtml), 0644)
	println("Done! Open index.html in your browser to see the result.")
}
