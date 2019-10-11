package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"strings"
)

var (
	// global
	ghe gherror
	// variables
	MinIssueTitleMatchPercentage uint = 80
	IssueLabels                       = []string{"go", "gherror"}
)

type gherror struct {
	githubClient *github.Client
	repoOwner    string
	repoSlug     string
}

func Register(githubToken string, repoSlug string) error {
	if strings.TrimSpace(githubToken) == "" {
		return errors.New("github token must not be empty")
	}

	if strings.TrimSpace(repoSlug) == "" {
		return errors.New("repository slug must not be empty")
	}

	repoSlugParts := strings.Split(repoSlug, "/")
	if len(repoSlugParts) != 2 {
		return errors.New("repository slug must be format <author>/<repo>")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	oauthClient := oauth2.NewClient(ctx, ts)

	ghe = gherror{
		githubClient: github.NewClient(oauthClient),
		repoOwner:    repoSlugParts[0],
		repoSlug:     repoSlugParts[1],
	}

	return nil
}

func Report(codeError error, metadata map[string]string) error {
	if ghe == (gherror{}) {
		return errors.New("gherror must be initialized first with Register()")
	}

	if codeError == nil || codeError.Error() == "" {
		return nil
	}

	title := codeError.Error()
	found, err := hasComparableIssue(title, MinIssueTitleMatchPercentage)
	if err != nil {
		return errors.New("could not find comparable issue: " + err.Error())
	}

	if found {
		// TODO: comment about reoccurence with max of eg 30 comments
		return nil
	}

	issueBody := fmt.Sprintf("%+v\n%+v", codeError, metadata)
	if err := createGithubIssue(codeError.Error(), issueBody, IssueLabels); err != nil {
		return errors.New("could not create github issue: " + err.Error())
	}

	return nil
}
