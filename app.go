package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Website struct {
	Name string
	URL  string
}

type Template struct {
	Name string
	Path string
}

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

		templateDirectory := getTemplatesDirectory()
		if templateDirectory == "" {
			log.Fatalf("Error getting templates directory")
		}

		println("Templates directory: " + templateDirectory)
		println("Available templates:")
		templates := getTemplates()

		for _, tpl := range templates {
			println(tpl.Name)
		}

		os.Exit(0)
	}

	// Verify and Pick template
	selectedTemplateName := command

	templates := getTemplates()
	templatePath := ""

	for _, tpl := range templates {
		if tpl.Name == selectedTemplateName {
			templatePath = tpl.Path
			break
		}
	}

	if templatePath == "" {
		log.Fatalf("Template %s not found. Use 'surge-uploader list' to see available templates.", selectedTemplateName)
	}

	// Get websites from Surge
	websites, err := getWebsitesFromSurge()
	if err != nil {
		log.Fatalf("Error getting websites from Surge: %v", err)
	}

	renderedHtml, err := renderTemplate(templatePath, websites)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	println("Rendering template to index.html")
	os.WriteFile("index.html", []byte(renderedHtml), 0644)
	println("Done! Open index.html in your browser to see the result.")

}

func getWebsitesFromSurge() ([]Website, error) {

	var websites []Website

	out, err := exec.Command("surge", "list").Output()
	if err != nil {
		return nil, err
	}

	elements := strings.Split(string(out), "\n")

	if len(elements) < 2 {
		return nil, errors.New("unexpected output from surge list command, Check if surge is installed and logged in")
	}

	for _, line := range elements {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			log.Printf("Skipping line: %s", line)
			continue
		}
		url := "https://" + strings.ReplaceAll(parts[1], "[39m", "")
		name := convertUrlToName(url)

		websites = append(websites, Website{Name: name, URL: url})
	}

	return websites, nil
}

func convertUrlToName(url string) string {
	// Remove the protocol (http:// or https://)
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// Remove any trailing slashes
	url = strings.TrimSuffix(url, "/")

	// Replace dots and slashes with hyphens
	name := strings.ReplaceAll(url, ".", "-")
	name = strings.ReplaceAll(name, "/", "-")

	// Remove the .surge.sh suffix if present
	name = strings.TrimSuffix(name, ".surge.sh")

	return name
}

func getTemplatesDirectory() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	tplDir := filepath.Join(homeDir, ".surge-templates")
	return tplDir
}

/*
Get the list of templates from the home directory of the user under surge-templates. if there is not this direcotry create it and add a default template to it.
*/
func getTemplates() []Template {

	tplDir := getTemplatesDirectory()

	// If the surge-templates directory does not exist, create it and add a default template
	if _, err := os.Stat(tplDir); os.IsNotExist(err) {
		err := os.MkdirAll(tplDir, 0755)
		if err != nil {
			return nil
		}

		defaultTpl := filepath.Join(tplDir, "default.html")
		_ = os.WriteFile(defaultTpl, []byte(getDefaultTemplate()), 0644)

		return []Template{{Name: "default", Path: defaultTpl}}

	}

	// Read the templates from the surge-templates directory, there are all .html files in the directory
	files, err := os.ReadDir(tplDir)
	if err != nil {
		return nil
	}

	var templates []Template
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".html" {
			continue
		}

		templates = append(templates, Template{Name: strings.TrimSuffix(file.Name(), ".html"), Path: filepath.Join(tplDir, file.Name())})

	}

	return templates
}

func renderTemplate(templatePath string, websites []Website) (string, error) {

	// Read the template file
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// Parse the template html file, found the template element in the content of the file with name "card", inside it will be the html code for each website card
	htmlTemplate := string(content)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlTemplate))
	if err != nil {
		return "", err
	}

	contentDiv := doc.Find("div#content").First()
	if contentDiv.Length() == 0 {
		return "", errors.New("div with id 'content' not found in the template file")
	}

	cardTemplateHtml := doc.Find("template#card").First()
	cardHtml, err := cardTemplateHtml.Html()
	if err != nil {
		return "", errors.New("template with id 'card' not found in the template file")
	}

	if cardHtml == "" {
		return "", errors.New("template with id 'card' is empty in the template file")
	}

	cardTemplate, err := template.New("card").Parse(cardHtml)
	if err != nil {
		return "", err
	}

	for _, website := range websites {
		var renderedCard strings.Builder
		err := cardTemplate.Execute(&renderedCard, website)
		if err != nil {
			return "", err
		}

		// Append the rendered card to the content div
		contentDiv.AppendHtml(renderedCard.String())
	}

	// Remove the template element from the document
	cardTemplateHtml.Remove()

	renderedHtml, err := doc.Html()
	if err != nil {
		return "", err
	}

	return renderedHtml, nil
}

func getDefaultTemplate() string {
	return `
<!doctype html>
<html>

<head>
	<title>Projects</title>
	<style>
		.card {
			border: 1px solid #ccc;
			padding: 10px;
			margin: 10px;
			border-radius: 5px;
		}

		#content {
			display: flex;
			flex-wrap: wrap;
			gap: 0.3rem;
		}

		.card-link {
			color: #a3a3ff;
			text-decoration: none;
		}
	</style>
</head>

<body>
	<template id="card">
		<div class="card">
			<p class="card-title">{{.Name}}</p>
			<a href="{{.URL}}" class="card-link">Visit</a>
		</div>
	</template>

	<div id="content"></div>
</body>

</html>
	`
}
