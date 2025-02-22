PKG = github.com/k1LoW/octocov
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on
export CGO_ENABLED=1

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test test_no_coverage

test:
	go test ./... -coverprofile=coverage.out -covermode=count

test_central: build
	./octocov --config testdata/octocov_central.yml

test_no_coverage: build
	./octocov --config testdata/octocov_no_coverage.yml

lint:
	golangci-lint run ./...
	govulncheck ./...
	go vet -vettool=`which gostyle` -gostyle.config=$(PWD)/.gostyle.yml ./...

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

build_badgen:
	go build -ldflags="$(BUILD_LDFLAGS)" -o badgen pkg/badge/cmd/badgen/main.go

coverage: build
	./octocov

bqdoc:
	cd docs/bq && tbls doc -f

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/k1LoW/gostyle@latest

prerelease:
	git pull origin main --tag
	go mod tidy
	ghch -w -N ${VER}
	gocredits -skip-missing -w
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

prerelease_for_tagpr:
	gocredits -skip-missing -w
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum

.PHONY: default test
