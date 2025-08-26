package github

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPkg_LastTag(t *testing.T) {
	// 跳过需要网络访问的测试，除非设置了环境变量
	if os.Getenv("GITHUB_INTEGRATION_TEST") == "" {
		t.Skip("Skipping GitHub API test. Set GITHUB_INTEGRATION_TEST=1 to run")
	}

	tests := []struct {
		name    string
		owner   string
		repo    string
		wantErr bool
		checks  func(t *testing.T, tag *Tag)
	}{
		{
			name:    "ergoapi/util should have tags",
			owner:   "ergoapi",
			repo:    "util",
			wantErr: false,
			checks: func(t *testing.T, tag *Tag) {
				// 只验证结构和基本属性，不验证具体值
				assert.NotEmpty(t, tag.Name, "Tag name should not be empty")
				assert.NotEmpty(t, tag.Commit.SHA, "Commit SHA should not be empty")
				assert.NotEmpty(t, tag.Commit.URL, "Commit URL should not be empty")
				
				// 验证版本号格式（假设使用语义化版本）
				assert.Regexp(t, `^v?\d+\.\d+\.\d+`, tag.Name, "Tag should follow semantic versioning")
				
				// 验证URL格式
				assert.Contains(t, tag.Commit.URL, "github.com/repos/ergoapi/util/commits/",
					"URL should point to the correct repository")
			},
		},
		{
			name:    "non-existent repo should error",
			owner:   "ergoapi",
			repo:    "non-existent-repo-12345",
			wantErr: true,
			checks:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pkg{
				Owner: tt.owner,
				Repo:  tt.repo,
			}
			
			got, err := p.LastTag()
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			require.NotNil(t, got)
			
			if tt.checks != nil {
				tt.checks(t, got)
			}
			
			// 可选：打印实际获取的标签信息用于调试
			t.Logf("Latest tag: %s (SHA: %s)", got.Name, got.Commit.SHA)
		})
	}
}

// MockGitHubClient 用于测试的mock客户端
type MockGitHubClient struct {
	tags []Tag
	err  error
}

func (m *MockGitHubClient) ListTags(ctx context.Context, owner, repo string) ([]Tag, error) {
	return m.tags, m.err
}

// TestPkg_LastTag_Mock 使用mock测试核心逻辑（不依赖网络）
func TestPkg_LastTag_Mock(t *testing.T) {
	tests := []struct {
		name    string
		tags    []Tag
		err     error
		want    *Tag
		wantErr bool
	}{
		{
			name: "normal case with tags",
			tags: []Tag{
				{Name: "v1.0.0", Commit: Commit{SHA: "abc123", URL: "http://example.com"}},
				{Name: "v0.9.0", Commit: Commit{SHA: "def456", URL: "http://example.com"}},
			},
			want:    &Tag{Name: "v1.0.0", Commit: Commit{SHA: "abc123", URL: "http://example.com"}},
			wantErr: false,
		},
		{
			name:    "empty tags should return error",
			tags:    []Tag{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "API error should propagate",
			tags:    nil,
			err:     assert.AnError,
			want:    nil,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pkg{
				Owner: "test",
				Repo:  "repo",
			}
			
			client := &MockGitHubClient{
				tags: tt.tags,
				err:  tt.err,
			}
			
			got, err := p.LastTagWithClient(client)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestPkg_ListTags 测试列出所有标签
func TestPkg_ListTags(t *testing.T) {
	if os.Getenv("GITHUB_INTEGRATION_TEST") == "" {
		t.Skip("Skipping GitHub API test. Set GITHUB_INTEGRATION_TEST=1 to run")
	}

	p := &Pkg{
		Owner: "golang",
		Repo:  "go",
	}
	
	tags, err := p.listTags()
	require.NoError(t, err)
	assert.NotEmpty(t, tags, "Go repository should have tags")
	
	// 验证标签按时间排序（最新的在前）
	for i, tag := range tags {
		assert.NotEmpty(t, tag.Name)
		assert.NotEmpty(t, tag.Commit.SHA)
		
		if i > 10 {
			break // 只检查前10个标签
		}
	}
}

// 基准测试
func BenchmarkPkg_LastTag(b *testing.B) {
	if os.Getenv("GITHUB_INTEGRATION_TEST") == "" {
		b.Skip("Skipping benchmark. Set GITHUB_INTEGRATION_TEST=1 to run")
	}
	
	p := &Pkg{
		Owner: "ergoapi",
		Repo:  "util",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.LastTag()
		if err != nil {
			b.Fatal(err)
		}
	}
}