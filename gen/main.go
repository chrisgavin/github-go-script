package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const (
	gomainfile = `package main

import (
	"fmt"

	github_go_script "github.com/dsp-testing/github-go-script"
)

func main() {
	github_go_script.Call(func(options *github_go_script.Options) (map[string]string, error) {
		// TODO: Replace this with your code
		repo, _, err := options.Client.Repositories.Get(options.Ctx, options.Context.Repository.OwnerName, options.Context.Repository.Name)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Repository name: %s\n", repo.GetName())
		fmt.Printf("Triggering event name: %s\n", options.Context.EventName)
		return map[string]string{
			"repository_name": repo.GetName(),
			"event_name":      options.Context.EventName,
		}, nil
	})
}`
	workflowTemplate = `on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]
  workflow_dispatch:
  
jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - uses: dsp-testing/github-go-script@main
        with:
          dir: .github/workflows/{{ .Name }}
`
)

func realMain() error {
	if len(os.Args) > 2 {
		return fmt.Errorf("usage: gen [workflow name]")
	}

	name := "github-go-script"
	if len(os.Args) == 2 {
		name = os.Args[1]
	}

	dir := ".github/workflows/" + name
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(workflowTemplate)
	if err != nil {
		return err
	}

	workflowFile, err := os.Create(dir + ".yml")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(workflowFile, struct{ Name string }{name}); err != nil {
		return err
	}

	if err := os.Chdir(dir); err != nil {
		return err
	}

	if out, err := exec.Command("go", "mod", "init", "_").CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	if out, err := exec.Command("go", "get", "github.com/dsp-testing/github-go-script").CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	if err := os.WriteFile("main.go", []byte(gomainfile), 0644); err != nil {
		return err
	}

	if out, err := exec.Command("go", "mod", "tidy").CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}

	return nil
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
