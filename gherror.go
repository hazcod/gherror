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

	ctx := context.Background()

	issueBody := fmt.Sprintf("%+v\n%+v", codeError, metadata)

	title := codeError.Error()
	comparableIssue, err := hasComparableIssue(title, MinIssueTitleMatchPercentage)
	if err != nil {
		return errors.New("could not find comparable issue: " + err.Error())
	}

	if comparableIssue != nil {
		shouldCreate, err := shouldCreateComment(ctx, comparableIssue)
		if err != nil {
			return errors.New("could not check comments: " + err.Error())
		}

		if shouldCreate {
			if err := createComment(ctx, comparableIssue, issueBody); err != nil {
				return errors.New("could not create comment: " + err.Error())
			}
		}
	} else {
		if err := createGithubIssue(codeError.Error(), issueBody, IssueLabels); err != nil {
			return errors.New("could not create github issue: " + err.Error())
		}
	}

	return nil
}
