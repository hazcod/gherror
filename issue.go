package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v28/github"
	"strings"
	"time"
)

const (
	MaxIssueComments = 20
)

func validGithubResponse(action string, err error, resp *github.Response) error {
	if err != nil {
		return errors.New(fmt.Sprintf("could not %s: %v", action, err))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 399 {
		return errors.New(fmt.Sprintf("%s returned status code %d", action, resp.StatusCode))
	}

	return nil
}

func shouldCreateComment(ctx context.Context, issue *github.Issue) (shouldCreate bool, err error) {
	comments, resp, err := ghe.githubClient.Issues.ListComments(ctx, ghe.repoOwner, ghe.repoSlug, *issue.Number, &github.IssueListCommentsOptions{
		Sort:        "",
		Direction:   "",
		Since:       time.Time{},
		ListOptions: github.ListOptions{},
	})

	if err := validGithubResponse("shouldCreateComment", err, resp); err != nil {
		return false, err
	}

	if len(comments) >= MaxIssueComments {
		return false, nil
	}

	return true, nil
}

func createComment(ctx context.Context, issue *github.Issue, body string) (err error) {
	_, resp, err := ghe.githubClient.Issues.CreateComment(ctx, ghe.repoOwner, ghe.repoSlug, *issue.Number, &github.IssueComment{Body: &body})

	if err := validGithubResponse("createComment", err, resp); err != nil {
		return err
	}

	return nil
}

func createGithubIssue(title string, body string, labels []string) error {
	issueRequest := &github.IssueRequest{
		Title:     &title,
		Body:      &body,
		Labels:    &labels,
		Assignee:  nil,
		State:     nil,
		Milestone: nil,
		Assignees: nil,
	}

	_, resp, err := ghe.githubClient.Issues.Create(context.Background(), ghe.repoOwner, ghe.repoSlug, issueRequest)
	if err := validGithubResponse("createIssue", err, resp); err != nil {
		return err
	}

	return nil
}

func stringCompare(s1, s2 string) uint {
	lens := len(s1)

	if lens > len(s2) {
		lens = len(s2)
	}

	for i := 0; i < lens; i++ {
		if s1[i] != s2[i] {
			return uint(int(s1[i]) - int(s2[i]))
		}
	}

	return uint(len(s1) - len(s2))
}

func hasComparableIssue(title string, minPercentage uint) (issue *github.Issue, err error) {
	title = strings.TrimSpace(strings.ToLower(title))

	issues, resp, err := ghe.githubClient.Issues.ListByRepo(context.Background(), ghe.repoOwner, ghe.repoSlug, &github.IssueListByRepoOptions{
		Milestone:   "",
		State:       "",
		Assignee:    "",
		Creator:     "",
		Mentioned:   "",
		Labels:      nil,
		Sort:        "",
		Direction:   "",
		Since:       time.Time{},
		ListOptions: github.ListOptions{},
	})
	if err := validGithubResponse("listIssues", err, resp); err != nil {
		return nil, err
	}

	for _, issue := range issues {
		otherTitle := strings.TrimSpace(strings.ToLower(*issue.Title))
		difference := stringCompare(title, otherTitle) * 100 / uint(len(title))

		//log.Printf("'%s' ?= '%s' -> difference %d", otherTitle, title, difference)

		if (100 - difference) >= minPercentage {
			return issue, nil
		}
	}

	return nil, nil
}
