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

| Name                  | Description                                                                                                       | Default             | Required |
| --------------------- | ----------------------------------------------------------------------------------------------------------------- | ------------------- | -------- |
| bump_default_level    | Default bump level if labels are not attached [major, minor, patch, none]. Do nothing if it's empty               |                     | false    |
| bump_fail_if_no_level | Fail if no bump label is found                                                                                    | false               | false    |
| bump_include_v        | Include `v` prefix in tag                                                                                         | true                | false    |
| bump_latest           | Add `latest` tag                                                                                                  | false               | false    |
| bump_major            | Label name for major bump (bumper:major)                                                                          | bumper:major        | false    |
| bump_minor            | Label name for minor bump (bumper:minor)                                                                          | bumper:minor        | false    |
| bump_none             | Label name for no bump (bumper:none)                                                                              | bumper:none         | false    |
| bump_patch            | Label name for patch bump (bumper:patch)                                                                          | bumper:patch        | false    |
| bump_semver           | Whether to updates major/minor release tags on a tag push. e.g. Update `v1` and `v1.2` tag when released `v1.2.3` | false               | false    |
| bump_tag_as_email     | Email address to use when creating tags                                                                           |                     | false    |
| bump_tag_as_user      | Name to use when creating tags                                                                                    |                     | false    |
| github_token          | GITHUB_TOKEN to list pull requests and create tags                                                                | ${{ github.token }} | true     |

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
  pull_request:
    branches:
      - main
    types:
      - labeled
      - unlabeled
      - opened
      - reopened
      - synchronize

permissions:
  pull-requests: write
  contents: write

jobs:
  tag_bumper:
    runs-on: ubuntu-latest
    env:
      INETUM_POLAND_ACTION_BUMPER_DEBUG: ${{ vars.INETUM_POLAND_ACTION_BUMPER_DEBUG }}
      INETUM_POLAND_ACTION_BUMPER_TRACE: ${{ vars.INETUM_POLAND_ACTION_BUMPER_TRACE }}
    steps:
      - uses: hmarr/debug-action@v3

      - uses: actions/checkout@v4

      - uses: jwalton/gh-find-current-pr@v1
        id: finder
        with:
          state: all

      - id: bumper
        uses: ./
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          bump_fail_if_no_level: true
          bump_semver: true
          bump_latest: true

      - uses: marocchino/sticky-pull-request-comment@v2
        if: always() && (steps.bumper.outputs.tag_status != null) && (steps.finder.outputs.pr != null)
        with:
          header: action-bumper
          number: ${{ steps.finder.outputs.pr }}
          message: ${{ steps.bumper.outputs.tag_status }}
```

### Note

action-bumper uses push on main event to run workflow instead of pull_request closed (merged) event because github token doesn't have write permission for pull_request from fork repository.

---

Copyright (c) 2024 Inetum Poland.
