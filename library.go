package github_go_script

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Repository struct {
	ID        uint64
	OwnerName string
	Name      string
}

type Context struct {
	EventName  string
	Sha        string
	Ref        string
	Workflow   string
	Action     string
	Actor      string
	Job        string
	RunNumber  uint64
	RunId      uint64
	ApiUrl     string
	ServerUrl  string
	GraphqlUrl string
	Repository *Repository
	RawPayload []byte
}

func (c *Context) ParsePayload(target interface{}) {
	err := json.Unmarshal(c.RawPayload, target)
	if err != nil {
		panic(err)
	}
}

type Options struct {
	Ctx     context.Context
	Client  *github.Client
	Context *Context
}

func MustUint64(input string) uint64 {
	i, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

type run func(options *Options) error

func Call(r run) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repositoryOwnerName, repositoryName, found := strings.Cut(os.Getenv("GITHUB_REPOSITORY"), "/")
	if !found {
		panic("GITHUB_REPOSITORY does not contain a slash")
	}

	payload, err := os.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		panic(err)
	}

	gitHubContext := &Context{
		EventName:  os.Getenv("GITHUB_EVENT_NAME"),
		Sha:        os.Getenv("GITHUB_SHA"),
		Ref:        os.Getenv("GITHUB_REF"),
		Workflow:   os.Getenv("GITHUB_WORKFLOW"),
		Action:     os.Getenv("GITHUB_ACTION"),
		Actor:      os.Getenv("GITHUB_ACTOR"),
		Job:        os.Getenv("GITHUB_JOB"),
		RunNumber:  MustUint64(os.Getenv("GITHUB_RUN_NUMBER")),
		RunId:      MustUint64(os.Getenv("GITHUB_RUN_ID")),
		ApiUrl:     os.Getenv("GITHUB_API_URL"),
		ServerUrl:  os.Getenv("GITHUB_SERVER_URL"),
		GraphqlUrl: os.Getenv("GITHUB_GRAPHQL_URL"),
		Repository: &Repository{
			ID:        MustUint64(os.Getenv("GITHUB_REPOSITORY_ID")),
			OwnerName: repositoryOwnerName,
			Name:      repositoryName,
		},
		RawPayload: payload,
	}

	options := &Options{
		Ctx:     ctx,
		Client:  client,
		Context: gitHubContext,
	}

	err = r(options)
	if err != nil {
		fmt.Fprintln(err, os.Stderr)
		os.Exit(1)
	}
}
