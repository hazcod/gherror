package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v28/github"
	"strconv"
	"strings"
	"time"
)

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

	_, response, err := ghe.githubClient.Issues.Create(context.Background(), ghe.repoOwner, ghe.repoSlug, issueRequest)
	if err != nil {
		return errors.New("could not create github issue: " + err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 399 {
		return errors.New("github issue creation returned status code: " + strconv.Itoa(response.StatusCode))
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

func hasComparableIssue(title string, minPercentage uint) (found bool, err error) {
	title = strings.TrimSpace(strings.ToLower(title))

	issues, _, err := ghe.githubClient.Issues.ListByRepo(context.Background(), ghe.repoOwner, ghe.repoSlug, &github.IssueListByRepoOptions{
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
	if err != nil {
		return false, errors.New(fmt.Sprintf("could not list issues for %s/%s: %v", ghe.repoOwner, ghe.repoSlug, err))
	}

	for _, issue := range issues {
		otherTitle := strings.TrimSpace(strings.ToLower(*issue.Title))
		difference := stringCompare(title, otherTitle) * 100 / uint(len(title))

		if (100 - difference) >= minPercentage {
			return true, nil
		}
	}

	return false, nil
}
