package gh

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v58/github"
	"github.com/k1LoW/go-github-client/v58/factory"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in      string
		want    *Repository
		wantErr bool
	}{
		{"owner/repo", &Repository{Owner: "owner", Repo: "repo"}, false},
		{"owner/repo/path/to", &Repository{Owner: "owner", Repo: "repo", Path: "path/to"}, false},
		{"owner/repo@sub", &Repository{Owner: "owner", Repo: "repo@sub"}, false},
		{"owner/repo.sub", &Repository{Owner: "owner", Repo: "repo.sub"}, false},
		{"owner/../sub", nil, true},
		{"owner", nil, true},
		{"owner/../sub", nil, true},
		{"owner/./sub", nil, true},
		{"owner//sub", nil, true},
		{"owner/repo/sub/", nil, true},
	}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("got error %v\n", err)
			}
			continue
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Error(diff)
		}
	}
}

func TestFetchDefaultBranch(t *testing.T) {
	mg := mockedGh(t)
	want := "main"
	got, err := mg.FetchDefaultBranch(context.TODO(), "owner", "repo")
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestFetchRawRootURL(t *testing.T) {
	ctx := context.TODO()
	token, _, _, _ := factory.GetTokenAndEndpoints()
	if token == "" {
		t.Skip("no token")
		return
	}
	tests := []struct {
		owner string
		repo  string
		want  string
	}{
		{"k1LoW", "octocov", "https://raw.githubusercontent.com/k1LoW/octocov/main"},
	}
	for _, tt := range tests {
		g, err := New()
		if err != nil {
			t.Fatal(err)
		}
		got, err := g.FetchRawRootURL(ctx, tt.owner, tt.repo)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}

func TestDetectCurrentBranch(t *testing.T) {
	tests := []struct {
		GITHUB_REF      string
		GITHUB_HEAD_REF string
		want            string
		wantErr         bool
	}{
		{"refs/pull/8/head", "", "", true},
		{"refs/heads/name", "mybranch", "name", false},
		{"refs/heads/branch/branch/name", "", "branch/branch/name", false},
		{"refs/pull/8/head", "mybranch", "mybranch", false},
	}
	ctx := context.TODO()
	mg := mockedGh(t)
	for _, tt := range tests {
		t.Run(tt.GITHUB_REF, func(t *testing.T) {
			t.Setenv("GITHUB_REF", tt.GITHUB_REF)
			t.Setenv("GITHUB_HEAD_REF", tt.GITHUB_HEAD_REF)
			got, err := mg.DetectCurrentBranch(ctx)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got err: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want err")
			}
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestDetectCurrentPullRequestNumber(t *testing.T) {
	tests := []struct {
		GITHUB_PULL_REQUEST_NUMBER string
		GITHUB_REF                 string
		want                       int
		wantErr                    bool
	}{
		{"", "refs/pull/8/head", 8, false},
		{"", "refs/heads/branch/branch/name", 13, false},
		{"", "refs/8", 0, true},
		{"8", "", 8, false},
		{"str", "", 0, true},
	}
	ctx := context.TODO()
	mg := mockedGh(t)
	for _, tt := range tests {
		t.Run(tt.GITHUB_REF, func(t *testing.T) {
			t.Setenv("GITHUB_PULL_REQUEST_NUMBER", tt.GITHUB_PULL_REQUEST_NUMBER)
			t.Setenv("GITHUB_REF", tt.GITHUB_REF)
			got, err := mg.DetectCurrentPullRequestNumber(ctx, "owner", "repo")
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got err: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want err")
			}
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSig(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"", "<!-- octocov -->"},
		{"foo", "<!-- octocov:foo -->"},
	}
	for _, tt := range tests {
		got := generateSig(tt.key)
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
func mockedGh(t *testing.T) *Gh {
	t.Setenv("GITHUB_TOKEN", "dummy")
	mockedHTTPClient := mock.NewMockedHTTPClient( //nostyle:funcfmt
		mock.WithRequestMatch( //nostyle:funcfmt
			mock.GetReposByOwnerByRepo,
			github.Repository{
				DefaultBranch: github.String("main"),
			},
		),
		mock.WithRequestMatch( //nostyle:funcfmt
			mock.GetReposPullsByOwnerByRepo,
			[]*github.PullRequest{
				&github.PullRequest{
					Head: &github.PullRequestBranch{
						Ref: github.String("branch/branch/name"),
					},
					Number: github.Int(13),
				},
			},
		),
	)
	client, err := factory.NewGithubClient(factory.HTTPClient(mockedHTTPClient), factory.Timeout(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	g, err := New()
	if err != nil {
		t.Fatal(err)
	}
	g.SetClient(client)
	return g
}
