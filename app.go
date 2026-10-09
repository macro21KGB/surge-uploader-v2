package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Website struct {
	Name string
	URL  string
}

type Template struct {
	Name string
	Path string
}

func main() {

	output, err := getSurgeList()
	if err != nil {
		log.Fatal(err)
	}

	var websites []Website

	templates := getTemplates()

	for _, line := range output {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			log.Printf("Skipping line: %s", line)
			continue
		}
		url := "https://" + parts[1]
		name := convertUrlToName(url)

		websites = append(websites, Website{Name: name, URL: url})
	}

	for _, website := range websites {
		log.Printf("Name: %s, URL: %s", website.Name, website.URL)
	}

	for _, template := range templates {
		log.Printf("Template Name: %s, Path: %s", template.Name, template.Path)
	}

}

func getSurgeList() ([]string, error) {
	out, err := exec.Command("surge", "list").Output()
	if err != nil {
		return nil, err
	}

	elements := strings.Split(string(out), "\n")

	if len(elements) < 2 {
		return nil, errors.New("unexpected output from surge list command, Check if surge is installed and logged in")
	}

	return elements, nil
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

/*
Get the list of templates from the home directory of the user under surge-templates. if there is not this direcotry create it and add a default template to it.
*/
func getTemplates() []Template {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	tplDir := filepath.Join(homeDir, ".surge-templates")

	// If the surge-templates directory does not exist, create it and add a default template
	if _, err := os.Stat(tplDir); os.IsNotExist(err) {
		err := os.MkdirAll(tplDir, 0755)
		if err != nil {
			return nil
		}

		defaultTpl := filepath.Join(tplDir, "default.html")
		_ = os.WriteFile(defaultTpl, []byte("<!DOCTYPE html>\n<html>\n<head>\n\t<title>Default Template</title>\n</head>\n<body>\n\t{{.Content}}\n</body>\n</html>"), 0644)

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
