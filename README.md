<p align="center">
<img src="https://github.com/k1LoW/octocov/raw/main/docs/logo.png" width="200" alt="octocov">
</p>

[![build](https://github.com/k1LoW/octocov/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/octocov/actions) ![coverage](docs/coverage.svg) ![ratio](docs/ratio.svg) ![time](docs/time.svg)

`octocov` is a toolkit for collecting code metrics (code coverage, code to test ratio and test execution time).

Key features of `octocov` are:

- **[Support multiple coverage report formats](#supported-coverage-report-formats).**
- **[Support multiple code metrics](#supported-code-metrics).**
- **[Support for even generating coverage report badge](#generate-coverage-report-badge-self).**
- **[Have a mechanism to aggregate reports from multiple repositories](#store-report-to-central-datastore).**

## Getting Started

### On GitHub Actions

**:octocat: GitHub Actions for octocov is [here](https://github.com/k1LoW/octocov-action) !!**

First, run test with [coverage report output](#supported-coverage-report-formats).

For example, in case of Go language, add `-coverprofile=coverage.out` option as follows

``` console
$ go test ./... -coverprofile=coverage.out
```
Add `.octocov.yml` ( or `octocov.yml` ) file to your repository.

``` yaml
# .octocov.yml
coverage:
  paths:
    - coverage.out
codeToTestRatio:
  code:
    - '**/*.go'
    - '!**/*_test.go'
  test:
    - '**/*_test.go'
comment:
  enable: true
```

And set up a workflow file as follows and run octocov on GitHub Actions.

``` yaml
# .github/workflows/ci.yml
name: Test

on:
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      -
        uses: actions/checkout@v2
      -
        uses: actions/setup-go@v2
        with:
          go-version: 1.17
      -
        name: Run tests with coverage report output
        run: go test ./... -coverprofile=coverage.out
      -
        uses: k1LoW/octocov-action@v0
```

Then, octocov comment the report of the code metrics to the pull request.

![comment](docs/comment.png)

**Notice:** Note that only pull requests from the same repository can be commented on. This is because the workflow token of a forked pull request does not have write permission.

### On Terminal

octocov acts as a code metrics viewer on the terminal.

For example, in case of Go language, add `-coverprofile=coverage.out` option as follows

``` console
$ go test ./... -coverprofile=coverage.out
```

And run `octocov ls-files` , `octocov view [FILE...]` and `octocov diff [REPORT_A] [REPORT_B]`

![term](docs/term.svg)

## Usage example

### Comment report to pull request

By setting `comment:`, [comment the reports to pull request](https://github.com/k1LoW/octocov/pull/30#issuecomment-860188829).

![comment](docs/comment.png)

``` yaml
# .octocov.yml
comment:
  enable: true
  hideFooterLink: false # hide octocov link
```

octocov checks for **"Code Coverage"** by default. If it is running on GitHub Actions, it will also measure **"Test Execution Time"**.

If you want to measure **"Code to Test Ratio"**, set `codeToTestRatio:`.

``` yaml
comment:
  enable: true
codeToTestRatio:
  code:
    - '**/*.go'
    - '!**/*_test.go'
  test:
    - '**/*_test.go'
```

By setting `diff:` ( `diff.path:`  or `diff.datastores` ) additionally, it is possible to show differences from previous reports as well.

``` yaml
comment:
  enable: true
diff:
  datastores:
    - s3://bucket/reports
```

![img](docs/comment_with_diff.png)

### Check for acceptable score

By setting `coverage.acceptable:`, the condition of acceptable coverage is specified.

If this condition is not met, the command will exit with exit status `1`.

``` yaml
# .octocov.yml
coverage:
  acceptable: 60%
```

``` console
$ octocov
Error: code coverage is 54.9%. the condition in the `coverage.acceptable:` section is not met (`60%`)
```

By setting `codeToTestRatio.acceptable:`, the condition of acceptable "Code to Test Ratio" is specified.

If this condition is not met, the command will exit with exit status `1`.

``` yaml
# .octocov.yml
codeToTestRatio:
  acceptable: 1:1.2
  code:
    - '**/*.go'
    - '!**/*_test.go'
  test:
    - '**/*_test.go'
```

``` console
$ octocov
Error: code to test ratio is 1:1.1, the condition in the `codeToTestRatio.acceptable:` section is not met (`1:1.2`)
```

By setting `testExecutionTime.acceptable:`, the condition of acceptable "Test Execution Time" is specified **(on GitHub Actions only)** .

If this condition is not met, the command will exit with exit status `1`.

``` yaml
# .octocov.yml
testExecutionTime:
  acceptable: 1 min
```

``` console
$ octocov
Error: test execution time is 1m15s, the condition in the `testExecutionTime.acceptable:` section is not met (`1 min`)
```

### Generate report badges self.

By setting `*.badge.path:`, generate badges self.

``` yaml
# .octocov.yml
coverage:
  badge:
    path: docs/coverage.svg
```

``` yaml
# .octocov.yml
codeToTestRatio:
  badge:
    path: docs/ratio.svg
```

``` yaml
# .octocov.yml
testExecutionTime:
  badge:
    path: docs/time.svg
```

You can display the coverage badge without external communication by setting a link to this badge image in README.md, etc.

``` markdown
# mytool

![coverage](docs/coverage.svg)
```

![coverage](docs/coverage.svg)

### Push report badges self.

By setting `push:`, git push report badges self.

``` yaml
# .octocov.yml
coverage:
  badge:
    path: docs/coverage.svg
push:
  enable: true
```

### Store report to datastores

By setting `report:`, store the reports to datastores and local path.

``` yaml
# .octocov.yml
report:
  datastores:
    - github://owner/coverages/reports
    - s3://bucket/reports
```

``` yaml
# .octocov.yml
report:
  path: path/to/report.json
```

#### Supported datastores

- GitHub repository
- S3
- GCS
- BigQuery
- Local

### Central mode

By enabling `central:`, `octocov` acts as a central repository for collecting reports ( [example](example/central/README.md) ).

``` yaml
# .octocov.yml for central mode
central:
  enable: true
  root:                    .             # root directory or index file path of collected coverage reports pages. default: .
  reports:
    - bq://my-project/my-dataset/reports # datastore paths (URLs) where reports are stored. default: local://reports
  badges: badges                         # directory where badges are generated. default: badges
  push:
    enable: true                         # enable self git push
```

#### Supported datastores

- GitHub repository
- S3
- GCS
- BigQuery
- Local

### View code coverage report of file

`octocov ls-files` command can be used to list files logged in code coverage report.

`octocov view` (alias: `octocov cat`) command can be used to view the file coverage report.

![term](docs/term.svg)

## Configuration

### `repository:`

The name of the repository.

It should be in the format `owner/repo`.

By default, the value of the environment variable `GITHUB_REPOSITORY` is set.

In case of monorepo, code metrics can be reported to datastore separately by specifying `owner/repo/project-a` or `owner/repo@project-a`.

``` yaml
repository: k1LoW/octocov
```

### `coverage:`

Configuration for code coverage.

### `coverage.path:`

`coverage.path:` has been deprecated. Please use `coverage.paths:` instead.

### `coverage.paths:`

The path to the coverage report file.

If no path is specified, the default path for each coverage format will be scanned.

``` yaml
coverage:
  paths:
    - tests/coverage.xml
```

### `coverage.acceptable:`

acceptable coverage condition.

``` yaml
coverage:
  acceptable: 60%
```

``` yaml
coverage:
  acceptable: current >= 60% && diff >= 0.5%
```

The variables that can be used are as follows.

| value | description |
| --- | --- |
| `current` | Current code metrics value |
| `prev` | Previous value. This value is taken from `diff.datastores:`. |
| `diff` | The result of `current - prev` |

It is also possible to omit the expression as follows

| Omitted expression | Expanded expression |
| --- | --- |
| `60%` | `current >= 60%` |
| `> 60%` | `current > 60%` |

### `coverage.badge:`

Set this if want to generate the badge self.

### `coverage.badge.path:`

The path to the badge.

``` yaml
coverage:
  badge:
    path: docs/coverage.svg
```

### `codeToTestRatio:`

Configuration for code to test ratio.

### `codeToTestRatio.code:` `codeToTestRatio.test:`

Files to count.

``` yaml
codeToTestRatio:
  code:                  # files to count as "Code"
    - '**/*.go'
    - '!**/*_test.go'
  test:                  # files to count as "Test"
    - '**/*_test.go'
```

### `codeToTestRatio.acceptable:`

acceptable ratio condition.

``` yaml
codeToTestRatio:
  acceptable: 1:1.2
```

``` yaml
codeToTestRatio:
  acceptable: current >= 1.2 && diff >= 0.0
```

The variables that can be used are as follows.

| value | description |
| --- | --- |
| `current` | Current code metrics value |
| `prev` | Previous value. This value is taken from `diff.datastores:`. |
| `diff` | The result of `current - prev` |

It is also possible to omit the expression as follows

| Omitted expression | Expanded expression |
| --- | --- |
| `1:1.2` | `current >= 1.2` |
| `> 1:1.2` | `current > 1.2` |

### `codeToTestRatio.badge:`

Set this if want to generate the badge self.

### `codeToTestRatio.badge.path:`

The path to the badge.

``` yaml
codeToTestRatio:
  badge:
    path: docs/ratio.svg
```

### `testExecutionTime:`

Configuration for test execution time.

### `testExecutionTime.acceptable`

acceptable time condition.

``` yaml
testExecutionTime:
  acceptable: 1min
```

``` yaml
testExecutionTime:
  acceptable: current <= 1min && diff <= 1sec
```

The variables that can be used are as follows.

| value | description |
| --- | --- |
| `current` | Current code metrics value |
| `prev` | Previous value. This value is taken from `diff.datastores:`. |
| `diff` | The result of `current - prev` |

It is also possible to omit the expression as follows

| Omitted expression | Expanded expression |
| --- | --- |
| `1min` | `current <= 1min` |
| `< 1min` | `current < 1min` |

### `testExecutionTime.badge`

Set this if want to generate the badge self.

### `testExecutionTime.badge.path`

The path to the badge.

``` yaml
testExecutionTime:
  badge:
    path: docs/time.svg
```

### `push:`

Configuration for `git push` badges self.

### `push.enable:`

Enable / disable `git push`

``` yaml
push:
  enable: false
```

### `comment:`

Set this if want to comment report to pull request

### `comment.enable:`

Enable / disable comment.

``` yaml
comment:
  enable: true
```

`enable: true` can be omitted if any other parameters are set as follows.

``` yaml
comment:
  hideFooterLink: true
```

### `comment.hideFooterLink:`

Hide footer [octocov](https://github.com/k1LoW/octocov) link.

``` yaml
comment:
  hideFooterLink: true
```

### `comment.if:`

Conditions for commenting report.

``` yaml
# .octocov.yml
comment:
  if: github.event_name == 'pull_request'
```

### `diff:`

Configuration for comparing reports.

### `diff.path:`

Path of the report to compare.

``` yaml
diff:
  path: path/to/coverage.yml
```

``` yaml
diff:
  path: path/to/report.json
```

### `diff.datastores:`

Datastores where the report to be compared is stored.

``` yaml
diff:
  datastores:
    - local://.octocov       # Use .octocov/owner/repo/report.json
    - s3://my-bucket/reports # Use s3://my-bucket/reports/owner/repo/report.json
```

### `diff.if:`

Conditions for comparing reports

``` yaml
# .octocov.yml
report:
  if: github.event_name == 'pull_request'
  path: path/to/report.json
```

### `report:`

Configuration for reporting to datastores.

### `report.path:`

Path to save the report.

``` yaml
report:
  path: path/to/report.json
```

### `report.datastores:`

Datastores where the reports are saved.

``` yaml
report:
  datastores:
    - github://owner/coverages/reports
    - s3://bucket/reports
```

#### GitHub repository

Use `github://` scheme.

```
github://[owner]/[repo]@[branch]/[prefix]
```

**Required environment variables:**

- `GITHUB_TOKEN` or `OCTOCOV_GITHUB_TOKEN`
- `GITHUB_REPOSITORY` or `OCTOCOV_GITHUB_REPOSITORY`
- `GITHUB_API_URL` or `OCTOCOV_GITHUB_API_URL` (optional)

#### S3

Use `s3://` scheme.

```
s3://[bucket]/[prefix]
```

**Required permission:**

- `s3:PutObject`

**Required environment variables:**

- `AWS_ACCESS_KEY_ID` or `OCTOCOV_AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY` or `OCTOCOV_AWS_SECRET_ACCESS_KEY`
- `AWS_SESSION_TOKEN` or `OCTOCOV_AWS_SESSION_TOKEN` (optional)

#### GCS

Use `gs://` scheme.

```
gs://[bucket]/[prefix]
```

**Required permission:**

- `storage.objects.create`
- `storage.objects.delete`

**Required environment variables:**

- `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_APPLICATION_CREDENTIALS_JSON` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS_JSON`

#### BigQuery

Use `bq://` scheme.

```
bq://[project ID]/[dataset ID]/[table]
```

**Required permission:**

- `bigquery.datasets.get`
- `bigquery.tables.get`
- `bigquery.tables.updateData`

**Required environment variables:**

- `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_APPLICATION_CREDENTIALS_JSON` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS_JSON`

**Datastore schema:**

[Datastore schema](docs/bq/schema/README.md)

If you want to create a table, execute the following command ( require `bigquery.datasets.create` ).

``` console
$ octocov migrate-bq-table
```

#### Local

Use `local://` or `file://` scheme.

```
local://[path]
```

**Example:**

If the absolute path of `.octocov.yml` is `/path/to/.octocov.yml`

- `local://reports` ... `/path/to/reports` directory
- `local://.reports` ... `/path/to/reports` directory
- `local://../reports` ... `/path/reports` directory
- `local:///reports` ... `/reports` directory.

### `report.if:`

Conditions for saving a report.

``` yaml
# .octocov.yml
report:
  if: env.GITHUB_REF == 'refs/heads/main'
  datastores:
    - github://owner/coverages/reports
```

The variables available in the `if` section are as follows

| Variable name | Type | Description |
| --- | --- | --- |
| `year` | `int` | Year of current time (UTC) |
| `month` | `int` | Month of current time (UTC) |
| `day` | `int` | Day of current time (UTC) |
| `hour` | `int` | Hour of current time (UTC) |
| `weekday` | `int` | Weekday of current time (UTC) (Sunday = 0, ...) |
| `github.event_name` | `string` | Event name of GitHub Actions ( ex. `issues`, `pull_request` )|
| `github.event` | `object` | Detailed data for each event of GitHub Actions (ex. `github.event.action`, `github.event.label.name` ) |
| `env.<env_name>` | `string` | The value of a specific environment variable |
| `is_pull_request` | `boolean` | Whether the job is related to an pull request (ex. a job fired by `on.push` will be true if it is related to a pull request) |
| `is_default_branch` | `boolean` | Whether the job is related to default branch of repository |

### `central:`

### `central.enable:`

Enable / disable central mode.

``` yaml
central:
  enable: false
```

`enable: true` can be omitted if any other parameters are set as follows.

``` yaml
central:
  reports:
    datastores:
      - local://reports
      - gs://my-gcs-bucket/reports
```

:NOTICE: When central mode is enabled, other functions are automatically turned off.


### `central.root:`

The root directory or index file ( [index file example](example/central/README.md) ) path of collected coverage reports pages. default: `.`

``` yaml
central:
  root: path/to
```

### `central.reports:`

### `central.reports.datastores:`

Datastore paths (URLs) where reports are stored. default: `local://reports`

``` yaml
central:
  reports:
    datastores:
      - local://reports
      - gs://my-gcs-bucket/reports
```

#### Use GitHub repository as datastore

When using the central repository as a datastore, perform badge generation via on.push.

![github](docs/github.svg)

``` yaml
# .octocov.yml
report:
  datastores:
    - github://owner/central-repo/reports
```

``` yaml
# .octocov.yml for central repo
central:
  reports:
    datastores:
      - github://owner/central-repo/reports
  push:
    enable: true
```

or

``` yaml
# .octocov.yml for central repo
central:
  reports:
    datastores:
      - local://reports
  push:
    enable: true
```

#### Use S3 bucket as datastore

When using the S3 bucket as a datastore, perform badge generation via on.schedule.

![s3](docs/s3.svg)

``` yaml
# .octocov.yml
report:
  datastores:
    - s3://my-s3-bucket/reports
```

``` yaml
# .octocov.yml for central repo
central:
  reports:
    datastores:
      - s3://my-s3-bucket/reports
  push:
    enable: true
```

**Required permission (Central Repo):**

- `s3:GetObject`
- `s3:ListObject`

**Required environment variables (Central Repo):**

- `AWS_ACCESS_KEY_ID` or `OCTOCOV_AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY` or `OCTOCOV_AWS_SECRET_ACCESS_KEY`
- `AWS_SESSION_TOKEN` or `OCTOCOV_AWS_SESSION_TOKEN` (optional)

#### Use GCS bucket as datastore

![gcs](docs/gcs.svg)

When using the GCS bucket as a datastore, perform badge generation via on.schedule.

``` yaml
# .octocov.yml
report:
  datastores:
    - gs://my-gcs-bucket/reports
```

``` yaml
# .octocov.yml for central repo
central:
  reports:
    datastores:
      - gs://my-gcs-bucket/reports
  push:
    enable: true
```

**Required permission (Central Repo):**

- `storage.objects.get`
- `storage.objects.list`
- `storage.buckets.get`

**Required environment variables (Central Repo):**

- `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_APPLICATION_CREDENTIALS_JSON` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS_JSON`

#### Use BigQuery table as datastore

![gcs](docs/bq.svg)

When using the BigQuery table as a datastore, perform badge generation via on.schedule.

``` yaml
# .octocov.yml
report:
  datastores:
    - bq://my-project/my-dataset/reports
```

``` yaml
# .octocov.yml for central repo
central:
  reports:
    datastores:
      - bq://my-project/my-dataset/reports
  push:
    enable: true
```

**Required permission (Central Repo):**

- `bigquery.jobs.create`
- `bigquery.tables.getData`

**Required environment variables (Central Repo):**

- `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_APPLICATION_CREDENTIALS_JSON` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS` or `OCTOCOV_GOOGLE_APPLICATION_CREDENTIALS_JSON`

### `central.badges:`

### `central.badges.datastores:`

Datastore paths (URLs) where badges are generated. default: `local://badges`

``` yaml
central:
  badges:
    datastores:
      - local://badges
      - s3://my-s3-buckets/badges
```

### `central.push:`

Configuration for `git push` index file and badges self.

### `central.push.enable:`

Enable / disable `git push`

``` yaml
push:
  enable: true
```

### `central.if:`

Conditions for central mode.

``` yaml
# .octocov.yml
central:
  if: env.GITHUB_REF == 'refs/heads/main'
  reports:
    datastores:
      - s3://my-s3-bucket/reports
```

## Supported coverage report formats

octocov supports multiple coverage report formats.

And octocov searches for the default path for each format.

If you want to specify the path of the report file, set `coverage.path`

``` yaml
coverage:
  paths:
    - /path/to/coverage.txt
```

### Go coverage

**Default path:** `coverage.out`

### LCOV

**Default path:** `coverage/lcov.info`

Support `SF` `DA` only

### SimpleCov

**Default path:** `coverage/.resultset.json`

### Clover

**Default path:** `coverage.xml`

### Cobertura

**Default path:** `coverage.xml`

## Supported code metrics

- **Code Coverage**
- **Code to Test Ratio**
- **Test Execution Time** (on GitHub Actions only)

## Install

**deb:**

Use [dpkg-i-from-url](https://github.com/k1LoW/dpkg-i-from-url)

``` console
$ export OCTOCOV_VERSION=X.X.X
$ curl -L https://git.io/dpkg-i-from-url | bash -s -- https://github.com/k1LoW/octocov/releases/download/v$OCTOCOV_VERSION/octocov_$OCTOCOV_VERSION-1_amd64.deb
```

**RPM:**

``` console
$ export OCTOCOV_VERSION=X.X.X
$ yum install https://github.com/k1LoW/octocov/releases/download/v$OCTOCOV_VERSION/octocov_$OCTOCOV_VERSION-1_amd64.rpm
```

**apk:**

Use [apk-add-from-url](https://github.com/k1LoW/apk-add-from-url)

``` console
$ export OCTOCOV_VERSION=X.X.X
$ curl -L https://git.io/apk-add-from-url | sh -s -- https://github.com/k1LoW/octocov/releases/download/v$OCTOCOV_VERSION/octocov_$OCTOCOV_VERSION-1_amd64.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/octocov
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/octocov/releases)

**docker:**

```console
$ docker pull ghcr.io/k1low/octocov:latest
```
