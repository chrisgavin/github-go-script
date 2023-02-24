# GitHub Go Script Action
An action similar to actions/github-script that lets you use Go as the scripting language.

## Quickstart

To get up and running with a new workflow, simply run:

    go run github.com/dsp-testing/github-go-script/gen workflow

This will create:

    .github
    └── workflows
        ├── workflow
        │   ├── go.mod
        │   ├── go.sum
        │   └── main.go
        └── workflow.yml

You can jump straight into `.github/workflows/workflow/main.go` and start scripting away.
