# surge-uploader

A small command-line tool that builds an `index.html` page listing all the websites you have published on [Surge](https://surge.sh).

It runs `surge list` to get your projects, then fills an HTML template with one "card" per website.

## Requirements

- [Go](https://go.dev/) (to build the tool)
- The [Surge CLI](https://surge.sh/help/getting-started-with-surge), installed and logged in:

  ```sh
  npm install --global surge
  surge login
  ```

## Build

```sh
go build -o surge-uploader .
```

You can move the `surge-uploader` binary into a folder on your `PATH` (for example `~/.local/bin`) so you can run it from anywhere.

## Usage

```sh
surge-uploader help          # show the help message
surge-uploader list          # show the templates folder and the available templates
surge-uploader <template>    # build index.html using the chosen template
```

Example:

```sh
surge-uploader default
```

This creates (or overwrites) `index.html` in the **current directory**. Open it in your browser to see the result.

## How it works

1. **Find the template.** The tool looks in the templates folder for a file named `<template>.html`. If it is not there, it stops with an error.
2. **Get your websites.** It runs `surge list` and reads each line of the output. It takes the domain from the second column, adds `https://` in front, and builds a display name from the domain (dots and slashes become hyphens).
3. **Render the page.**
   - It reads the template file and finds the `<template id="card">` element.
   - It fills the card once for each website, using Go's [`text/template`](https://pkg.go.dev/text/template) syntax.
   - It adds every filled card to the `<div id="content">` element.
   - It removes the `<template id="card">` element from the page.
4. **Write the output.** The final HTML is saved as `index.html`.

## The templates folder

Templates are kept in a hidden folder called `.surge-templates` in your home directory:

| OS      | Path                                   |
| ------- | -------------------------------------- |
| Linux   | `/home/<user>/.surge-templates`        |
| macOS   | `/Users/<user>/.surge-templates`       |
| Windows | `C:\Users\<user>\.surge-templates`     |

The easiest way to find it is to run:

```sh
surge-uploader list
```

The first line of the output shows the full path to the folder.

If the folder does not exist yet, the tool creates it the first time you run it and adds a `default.html` template. The default template is only created along with the folder. If you delete `default.html` later, it will not come back unless you delete the whole folder.

Because the folder name starts with a dot, it is hidden. To see it, use `ls -a ~` in a terminal or turn on "show hidden files" in your file manager.

## Creating and editing templates

Each `.html` file in the templates folder is a template. Its name is the file name without `.html`. For example, `dark.html` becomes the `dark` template:

```sh
surge-uploader dark
```

Files with other extensions and any subfolders are ignored.

The easiest way to start a new template is to copy the default one:

```sh
cp ~/.surge-templates/default.html ~/.surge-templates/dark.html
```

### What a template needs

Every template **must** have these two elements:

1. **`<template id="card">`**: the HTML for one website card. It is repeated once for each website. It must not be empty.
2. **`<div id="content">`**: the element where the cards are added.

If either one is missing, the tool stops with an error.

### Available variables

Inside the card you can use these values:

| Variable   | Description                                     | Example                          |
| ---------- | ----------------------------------------------- | -------------------------------- |
| `{{.Name}}` | Display name made from the domain              | `my-site`                        |
| `{{.URL}}`  | Full URL of the website                        | `https://my-site.surge.sh`       |

Since the card uses Go's `text/template`, you can also use its other features, such as `{{if}}`, `{{printf}}` and pipelines. Keep actions inside HTML attributes (for example `href="{{.URL}}"`) simple. Quotes inside them may be escaped when the page is parsed.

Everything outside the card (the `<head>`, styles, scripts, headers, footers and so on) is copied to `index.html` unchanged, so you can design the rest of the page however you like.

### Example template

```html
<!doctype html>
<html>
<head>
  <title>My Projects</title>
  <style>
    body { font-family: sans-serif; background: #111; color: #eee; }
    #content { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 1rem; }
    .card { background: #222; padding: 1rem; border-radius: 8px; }
    .card a { color: #8ab4f8; }
  </style>
</head>
<body>
  <h1>My Surge Projects</h1>

  <template id="card">
    <div class="card">
      <h2>{{.Name}}</h2>
      <a href="{{.URL}}" target="_blank">{{.URL}}</a>
    </div>
  </template>

  <div id="content"></div>
</body>
</html>
```

Save this as `~/.surge-templates/dark.html` and run `surge-uploader dark`.

## Troubleshooting

- **`Template <name> not found`**: run `surge-uploader list` to see the templates you have. Check that the file is in `~/.surge-templates` and ends with `.html`.
- **`Error getting websites from Surge`**: make sure the `surge` command is installed, is on your `PATH`, and that you are logged in (`surge login`).
- **`div with id 'content' not found`** or **`template with id 'card' ...`**: your template is missing one of the required elements. See [What a template needs](#what-a-template-needs).
