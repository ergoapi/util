package github

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/go-github/v65/github"
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

// GitHubClient 接口用于测试时的依赖注入
type GitHubClient interface {
	ListTags(ctx context.Context, owner, repo string) ([]Tag, error)
}

// DefaultClient 默认的GitHub客户端实现
type DefaultClient struct {
	client *github.Client
}

// NewDefaultClient 创建默认客户端
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		client: github.NewClient(&http.Client{
			Timeout: time.Second * 10,
		}),
	}
}

// ListTags 实现GitHubClient接口
func (c *DefaultClient) ListTags(ctx context.Context, owner, repo string) ([]Tag, error) {
	tags, _, err := c.client.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}
	
	var result []Tag
	for _, tag := range tags {
		result = append(result, Tag{
			Name: tag.GetName(),
			Commit: Commit{
				SHA: tag.Commit.GetSHA(),
				URL: tag.Commit.GetURL(),
			},
		})
	}
	return result, nil
}

func (p *Pkg) listTags() (ts []Tag, err error) {
	client := NewDefaultClient()
	return client.ListTags(context.Background(), p.Owner, p.Repo)
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

// LastTagWithClient 允许注入自定义客户端（便于测试）
func (p *Pkg) LastTagWithClient(client GitHubClient) (*Tag, error) {
	tags, err := client.ListTags(context.Background(), p.Owner, p.Repo)
	if err != nil {
		return nil, err
	}
	if len(tags) < 1 {
		return nil, errors.New("no tags")
	}
	return &tags[0], nil
}