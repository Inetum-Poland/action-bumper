# action-bumper

[![.github/workflows/pre_commit.yml](https://github.com/Inetum-Poland/action-bumper/actions/workflows/pre_commit.yml/badge.svg)](https://github.com/Inetum-Poland/action-bumper/actions/workflows/pre_commit.yml) [![.github/workflows/trufflehog.yaml](https://github.com/Inetum-Poland/action-bumper/actions/workflows/trufflehog.yaml/badge.svg)](https://github.com/Inetum-Poland/action-bumper/actions/workflows/trufflehog.yaml)

_Original projects: [action-bumpr](https://github.com/haya14busa/action-bumpr), [update-semver](https://github.com/haya14busa/action-update-semver)._

**action-bumper** bumps semantic version tag on merging Pull Requests with specific labels (`bumper:major`,`bumper:minor`,`bumper:patch`, `bumper:none`), and creates semver tags on tag push.

> [!IMPORTANT]
> __This repository uses the [Conventional Commits](https://www.conventionalcommits.org/).__
>
> For more information please see the [Conventional Commits documentation](https://www.conventionalcommits.org/en/v1.0.0/#summary).

> [!IMPORTANT]
> __This repository uses the [pre-commit](https://pre-commit.com/).__
>
> Please be respectful while contributing and after cloning this repo install the pre-commit hooks.
> ```bash
> > pre-commit install --install-hooks -t pre-commit -t commit-msg
> ```
> For more information please see the [pre-commit documentation](https://pre-commit.com/).

## Input

| Name               | Description                                                                                                       | Default             | Required |
| ------------------ | ----------------------------------------------------------------------------------------------------------------- | ------------------- | -------- |
| default_bump_level | Default bump level if labels are not attached [major, minor, patch, none]. Do nothing if it's empty               |                     | false    |
| dry_run            | Do not actually tag next version if it's true                                                                     |                     | false    |
| github_token       | GITHUB_TOKEN to list pull requests and create tags                                                                | ${{ github.token }} | true     |
| tag_as_user        | Name to use when creating tags                                                                                    |                     | false    |
| tag_as_email       | Email address to use when creating tags                                                                           |                     | false    |
| bump_major         | Label name for major bump (bumper:major)                                                                          | bumper:major        | false    |
| bump_minor         | Label name for minor bump (bumper:minor)                                                                          | bumper:minor        | false    |
| bump_patch         | Label name for patch bump (bumper:patch)                                                                          | bumper:patch        | false    |
| bump_none          | Label name for no bump (bumper:none)                                                                              | bumper:none         | false    |
| fail_if_no_bump    | Fail if no bump label is found                                                                                    | false               | false    |
| bump_semver        | Whether to updates major/minor release tags on a tag push. e.g. Update `v1` and `v1.2` tag when released `v1.2.3` | false               | false    |
| include_v          | Include `v` prefix in tag                                                                                          | true                | false    |

## Output

| Name            | Description                                                |
| --------------- | ---------------------------------------------------------- |
| current_version | current version                                            |
| next_version    | next version                                               |
| skip            | True if release is skipped. e.g. No labels attached to PR. |
| message         | Tag message                                                |

## Usage

```yaml
on:
  push:
    branches:
      - main
    tags:
      - 'v?*.*.*'
  pull_request:
    branches:
      - main
    types:
      - labeled
      - unlabeled
      - opened
      - reopened
      - synchronize

jobs:
  tag_bumper:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      # Bump version on merging Pull Requests with specific labels.
      # (bumper:major, bumper:minor, bumper:patch, bumper:none)
      - id: bumper
        uses: inetum-poland/action-bumper@v2
        with:
          fail_if_no_bump: true
          bump_semver: true
```

### Note

action-bumper uses push on master event to run workflow instead of pull_request closed (merged) event because github token doesn't have write permission for pull_request from fork repository.
