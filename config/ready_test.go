package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/k1LoW/go-github-client/v58/factory"
	"github.com/k1LoW/octocov/gh"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestCoverageConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{},
			"coverage: is not set",
		},
		{
			&Config{
				Coverage: &Coverage{},
			},
			"coverage.paths: is not set",
		},
		{
			&Config{
				Coverage: &Coverage{
					Paths: []string{"path/to/coverage.out"},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.CoverageConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCodeToTestRatioConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{},
			"codeToTestRatio: is not set",
		},
		{
			&Config{
				CodeToTestRatio: &CodeToTestRatio{},
			},
			"codeToTestRatio.test: is not set",
		},
		{
			&Config{
				CodeToTestRatio: &CodeToTestRatio{
					Test: []string{"path/to/test/**"},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.CodeToTestRatioConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestTestExecutionTimeConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{},
			"testExecutionTime: is not set",
		},
		{
			&Config{
				TestExecutionTime: &TestExecutionTime{
					Steps: []string{},
				},
			},
			"coverage: is not set",
		},
		{
			&Config{
				TestExecutionTime: &TestExecutionTime{
					Steps: []string{
						"Run tests",
					},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.TestExecutionTimeConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestPushConfigReady(t *testing.T) {
	os.Setenv("GITHUB_EVENT_NAME", "pull_request")
	os.Setenv("GITHUB_EVENT_PATH", filepath.Join(rootTestdataDir(t), "config", "event_pull_request_opened.json"))
	os.Setenv("GITHUB_REF", "refs/pull/4/merge")
	mg := mockedGh(t)
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Repository: "owner/repo",
				gh:         mg,
			},
			"push: is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Push:       &Push{},
				gh:         mg,
			},
			"failed to traverse the Git root path",
		},
		{
			&Config{
				Repository: "owner/repo",
				Push:       &Push{},
				gh:         mg,
			},
			"failed to traverse the Git root path",
		},
		{
			&Config{
				Repository: "owner/repo",
				GitRoot:    "/path/to",
				Push:       &Push{},
				gh:         mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				GitRoot:    "/path/to",
				Push: &Push{
					If: "false",
				},
				gh: mg,
			},
			"the condition in the `if` section is not met (false)",
		},
	}
	for _, tt := range tests {
		err := tt.c.PushConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			fmt.Printf("%s\n", err)

			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCommentConfigReady(t *testing.T) {
	os.Setenv("GITHUB_REF", "refs/pull/123/merge")
	os.Setenv("GITHUB_EVENT_NAME", "pull_request")
	os.Setenv("GITHUB_EVENT_PATH", filepath.Join(rootTestdataDir(t), "config", "event_pull_request_opened.json"))
	mg := mockedGh(t)
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Repository: "owner/repo",
				gh:         mg,
			},
			"comment: is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Comment:    &Comment{},
				gh:         mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Comment:    &Comment{},
				gh:         mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Comment: &Comment{
					If: "false",
				},
				gh: mg,
			},
			"the condition in the `if` section is not met (false)",
		},
	}
	for _, tt := range tests {
		err := tt.c.CommentConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCoverageBadgeConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Coverage: &Coverage{
					Paths: []string{"path/to/coverage.xml"},
				},
			},
			"coverage.badge.path: is not set",
		},
		{
			&Config{
				Coverage: &Coverage{
					Paths: []string{"path/to/coverage.xml"},
					Badge: CoverageBadge{
						Path: "path/to/coverage.svg",
					},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.CoverageBadgeConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCodeToTestRatioBadgeConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				CodeToTestRatio: &CodeToTestRatio{
					Test: []string{
						"**_test.go",
					},
				},
			},
			"codeToTestRatio.badge.path: is not set",
		},
		{
			&Config{
				CodeToTestRatio: &CodeToTestRatio{
					Test: []string{
						"**_test.go",
					},
					Badge: CodeToTestRatioBadge{
						Path: "path/to/ratio.svg",
					},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.CodeToTestRatioBadgeConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestTestExecutionTimeBadgeConfigReady(t *testing.T) {
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				TestExecutionTime: &TestExecutionTime{
					Steps: []string{
						"Run tests",
					},
				},
			},
			"testExecutionTime.badge.path: is not set",
		},
		{
			&Config{
				TestExecutionTime: &TestExecutionTime{
					Steps: []string{
						"Run tests",
					},
					Badge: TestExecutionTimeBadge{
						Path: "path/to/time.svg",
					},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.TestExecutionTimeBadgeConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCentralConfigReady(t *testing.T) {
	mg := mockedGh(t)

	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{},
			"central: is not set",
		},
		{
			&Config{
				Central: &Central{},
				gh:      mg,
			},
			"repository: not set (or env GITHUB_REPOSITORY is not set)",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central:    &Central{},
				gh:         mg,
			},
			"central.reports.datastores is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
				},
				gh: mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
					If: "false",
				},
				gh: mg,
			},
			"the condition in the `if` section is not met (false)",
		},
	}
	for _, tt := range tests {
		err := tt.c.CentralConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestCentralPushConfigReady(t *testing.T) {
	mg := mockedGh(t)
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
				},
				gh: mg,
			},
			"central.push: is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
					Push: &Push{},
				},
				gh: mg,
			},
			"failed to traverse the Git root path",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
					Push: &Push{},
				},
				GitRoot: "/path/to",
				gh:      mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Central: &Central{
					Reports: CentralReports{
						Datastores: []string{
							"s3://bucket/reports",
						},
					},
					Push: &Push{
						If: "false",
					},
				},
				GitRoot: "/path/to",
				gh:      mg,
			},
			"the condition in the `if` section is not met (false)",
		},
	}
	for _, tt := range tests {
		err := tt.c.CentralPushConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestDiffConfigReady(t *testing.T) {
	os.Setenv("GITHUB_EVENT_NAME", "pull_request")
	os.Setenv("GITHUB_EVENT_PATH", filepath.Join(rootTestdataDir(t), "config", "event_pull_request_opened.json"))
	os.Setenv("GITHUB_REF", "refs/pull/4/merge")
	mg := mockedGh(t)
	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Repository: "owner/repo",
				gh:         mg,
			},
			"diff: is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Diff:       &Diff{},
				gh:         mg,
			},
			"diff.path: and diff.datastores: are not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Diff: &Diff{
					Path: "path/to/report.json",
				},
				gh: mg,
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Diff: &Diff{
					Path: "path/to/report.json",
					If:   "false",
				},
				gh: mg,
			},
			"the condition in the `if` section is not met (false)",
		},
	}
	for _, tt := range tests {
		err := tt.c.DiffConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func TestReportConfigReady(t *testing.T) {
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_EVENT_PATH", filepath.Join(rootTestdataDir(t), "config", "event_pull_request_opened.json"))
	t.Setenv("GITHUB_REF", "refs/pull/4/merge")
	t.Setenv("GITHUB_TOKEN", "token")

	tests := []struct {
		c    *Config
		want string
	}{
		{
			&Config{
				Repository: "owner/repo",
				gh:         mockedGh(t),
			},
			"report: is not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Report:     &Report{},
				gh:         mockedGh(t),
			},
			"report.datastores: and report.path: are not set",
		},
		{
			&Config{
				Repository: "owner/repo",
				Report: &Report{
					Datastores: []string{
						"s3://bucket/reports",
					},
				},
				gh: mockedGh(t),
			},
			"",
		},
		{
			&Config{
				Repository: "owner/repo",
				Report: &Report{
					Datastores: []string{
						"s3://bucket/reports",
					},
					If: "false",
				},
				gh: mockedGh(t),
			},
			"the condition in the `if` section is not met (false)",
		},
		{
			&Config{
				Repository: "owner/repo",
				Report: &Report{
					Datastores: []string{
						"s3://bucket/reports",
					},
					If: "\"pull_requests\" startsWith \"pull\"",
				},
				gh: mockedGh(t),
			},
			"",
		},
	}
	for _, tt := range tests {
		err := tt.c.ReportConfigReady()
		if err == nil && tt.want != "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want == "" {
			t.Errorf("got %v\nwant %v", err, tt.want)
			continue
		}
		if err != nil && tt.want != "" {
			if got := err.Error(); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		}
	}
}

func mockedGh(t *testing.T) *gh.Gh {
	mockedHTTPClient := mock.NewMockedHTTPClient( //nostyle:funcfmt
		mock.WithRequestMatch( //nostyle:funcfmt
			mock.GetReposByOwnerByRepo,
			github.Repository{
				DefaultBranch: github.String("main"),
			},
		),
		mock.WithRequestMatch( //nostyle:funcfmt
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			github.PullRequest{
				Number: github.Int(13),
				Draft:  github.Bool(true),
			},
		),
	)
	client, err := factory.NewGithubClient(factory.HTTPClient(mockedHTTPClient), factory.Timeout(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	g, err := gh.New()
	if err != nil {
		t.Fatal(err)
	}
	g.SetClient(client)
	return g
}
