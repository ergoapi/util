package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v52/github"
)

type Pkg struct {
	Owner string
	Repo  string
}

type Tag struct {
	Name   string `json:"name,omitempty"`
	Commit Commit `json:"commit,omitempty"`
}

type Commit struct {
	SHA string `json:"sha,omitempty"`
	URL string `json:"url,omitempty"`
}

func (p *Pkg) listTags() (ts []Tag, err error) {
	client := github.NewClient(nil)
	ctx := context.Background()
	tags, _, err := client.Repositories.ListTags(ctx, p.Owner, p.Repo, nil)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		ts = append(ts, Tag{
			Name: tag.GetName(),
			Commit: Commit{
				SHA: tag.Commit.GetSHA(),
				URL: tag.Commit.GetURL(),
			},
		})
	}
	return ts, nil
}

func (p *Pkg) LastTag() (*Tag, error) {
	tags, err := p.listTags()
	if err != nil {
		return nil, err
	}
	if len(tags) < 1 {
		return nil, errors.New("no tags")
	}
	return &tags[0], nil
}