package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/k1LoW/octocov/internal"
)

func TestMain(m *testing.M) {
	envCache := os.Environ()

	m.Run()

	if err := revertEnv(envCache); err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		wd      string
		path    string
		wantErr bool
	}{
		{testdataDir(t), "", false},
		{filepath.Join(testdataDir(t), "config"), "", false},
		{filepath.Join(testdataDir(t), "config"), ".octocov.yml", false},
		{filepath.Join(testdataDir(t), "config"), "no.yml", true},
	}
	for _, tt := range tests {
		c := New()
		c.wd = tt.wd
		if err := c.Load(tt.path); err != nil {
			if !tt.wantErr {
				t.Errorf("got %v\nwantErr %v", err, tt.wantErr)
			}
		} else {
			if tt.wantErr {
				t.Errorf("got %v\nwantErr %v", nil, tt.wantErr)
			}
		}
	}
}

func TestLoadConfigAndOmitEnableFlag(t *testing.T) {
	wd := filepath.Join(testdataDir(t), "config")
	p := ".octocov.yml"
	c := New()
	c.wd = wd
	if err := c.Load(p); err != nil {
		t.Fatal(err)
	}
	if !internal.IsEnable(c.Comment.Enable) {
		t.Errorf("got %v\nwant true", *c.Comment.Enable)
	}
}

func TestCoverageAcceptable(t *testing.T) {
	tests := []struct {
		cond    string
		cov     float64
		prev    float64
		wantErr bool
	}{
		{"60%", 50.0, 0, true},
		{"50%", 50.0, 0, false},
		{"49.9%", 50.0, 0, false},
		{"49.9", 50.0, 0, false},
		{">= 60%", 50.0, 0, true},
		{">= 50%", 50.0, 0, false},
		{">= 49.9%", 50.0, 0, false},
		{">= 49.9", 50.0, 0, false},
		{">=60%", 50.0, 0, true},
		{">=50%", 50.0, 0, false},
		{">=49.9%", 50.0, 0, false},
		{">=49.9", 50.0, 0, false},

		{"current >= 60%", 50.0, 0, true},
		{"current > prev", 50.0, 49.0, false},
		{"diff >= 0", 50.0, 49.0, false},
		{"current >= 50% && diff >= 0%", 50.0, 49.0, false},
	}
	for _, tt := range tests {
		if err := coverageAcceptable(tt.cov, tt.prev, tt.cond); err != nil {
			if !tt.wantErr {
				t.Errorf("got %v\nwantErr %v", err, tt.wantErr)
			}
		} else {
			if tt.wantErr {
				t.Errorf("got %v\nwantErr %v", nil, tt.wantErr)
			}
		}
	}
}

func TestCodeToTestRatioAcceptable(t *testing.T) {
	tests := []struct {
		cond    string
		ratio   float64
		prev    float64
		wantErr bool
	}{
		{"1:1", 1.0, 0, false},
		{"1:1.1", 1.0, 0, true},
		{"1", 1.0, 0, false},
		{"1.1", 1.0, 0, true},
		{">= 1:1", 1.0, 0, false},
		{">= 1:1.1", 1.0, 0, true},
		{">= 1", 1.0, 0, false},
		{">= 1.1", 1.0, 0, true},
		{">=1:1", 1.0, 0, false},
		{">=1:1.1", 1.0, 0, true},
		{">=1", 1.0, 0, false},
		{">=1.1", 1.0, 0, true},

		{"current >= 1.1", 1.2, 1.1, false},
		{"current > prev", 1.2, 1.1, false},
		{"diff >= 0", 1.2, 1.1, false},
		{"current >= 1.1 && diff >= 0", 1.2, 1.1, false},
	}
	for _, tt := range tests {
		if err := codeToTestRatioAcceptable(tt.ratio, tt.prev, tt.cond); err != nil {
			if !tt.wantErr {
				t.Errorf("got %v\nwantErr %v", err, tt.wantErr)
			}
		} else {
			if tt.wantErr {
				t.Errorf("got %v\nwantErr %v", nil, tt.wantErr)
			}
		}
	}
}

func TestTestExecutionTimeAcceptable(t *testing.T) {
	tests := []struct {
		cond    string
		ti      float64
		prev    float64
		wantErr bool
	}{
		{"1min", float64(time.Minute), 0, false},
		{"59s", float64(time.Minute), 0, true},
		{"61sec", float64(time.Minute), 0, false},
		{"<= 1min", float64(time.Minute), 0, false},
		{"<= 59s", float64(time.Minute), 0, true},
		{"<= 61sec", float64(time.Minute), 0, false},
		{"<=1min", float64(time.Minute), 0, false},
		{"<=59s", float64(time.Minute), 0, true},
		{"<=61sec", float64(time.Minute), 0, false},
		{"1 min", float64(time.Minute), 0, false},
		{"59 s", float64(time.Minute), 0, true},
		{"61 sec", float64(time.Minute), 0, false},

		{"1min1sec", float64(time.Minute), 0, false},
		{"<=1min1sec", float64(time.Minute), 0, false},
		{"<= 1 min 1 sec", float64(time.Minute), 0, false},
		{"current <= 1 min 1 sec", float64(time.Minute), 0, false},

		{"current <= 1min", float64(time.Minute), float64(59 * time.Second), false},
		{"current > prev", float64(time.Minute), float64(59 * time.Second), false},
		{"diff <= 1sec", float64(time.Minute), float64(59 * time.Second), false},
		{"current <= 1min && diff <= 1sec", float64(time.Minute), float64(59 * time.Second), false},
	}
	for _, tt := range tests {
		if err := testExecutionTimeAcceptable(tt.ti, tt.prev, tt.cond); err != nil {
			if !tt.wantErr {
				t.Errorf("got %v\nwantErr %v", err, tt.wantErr)
			}
		} else {
			if tt.wantErr {
				t.Errorf("got %v\nwantErr %v", nil, tt.wantErr)
			}
		}
	}
}

func revertEnv(envCache []string) error {
	if err := clearEnv(); err != nil {
		return err
	}
	for _, e := range envCache {
		splitted := strings.Split(e, "=")
		if err := os.Setenv(splitted[0], splitted[1]); err != nil {
			return err
		}
	}
	return nil
}

func clearEnv() error {
	for _, e := range os.Environ() {
		splitted := strings.Split(e, "=")
		if err := os.Unsetenv(splitted[0]); err != nil {
			return err
		}
	}
	return nil
}

func testdataDir(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := filepath.Abs(filepath.Join(filepath.Dir(wd), "testdata"))
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
