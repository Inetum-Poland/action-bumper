# action-bumper

_Original projects: [action-bumpr](https://github.com/haya14busa/action-bumpr), [update-semver](https://github.com/haya14busa/action-update-semver)._

**action-bumper** bumps semantic version tag on merging Pull Requests with specific labels (`bumper:major`,`bumper:minor`,`bumper:patch`, `bumper:none`).

## Input

| Name               | Description                                                                                                       | Default             | Required |
| ------------------ | ----------------------------------------------------------------------------------------------------------------- | ------------------- | -------- |
| default_bump_level | Default bump level if labels are not attached [major, minor, patch, none]. Do nothing if it's empty               |                     | false    |
| dry_run            | Do not actually tag next version if it's true                                                                     |                     | false    |
| github_token       | GITHUB_TOKEN to list pull requests and create tags                                                                | ${{ github.token }} | true     |
| tag_as_user        | Name to use when creating tags                                                                                    |                     | false    |
| tag_as_email       | Email address to use when creating tags                                                                           |                     | false    |
| bump_major         | Label name for major bump (bump:major)                                                                            | bumper:major        | false    |
| bump_minor         | Label name for minor bump (bump:minor)                                                                            | bumper:minor        | false    |
| bump_patch         | Label name for patch bump (bump:patch)                                                                            | bumper:patch        | false    |
| bump_none          | Label name for no bump (bump:none)                                                                                | bumper:none         | false    |
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

### Simple

```yaml
name: release
on:
  push:
    branches:
      - master
  pull_request:
    types:
      - labeled

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      # Bump version on merging Pull Requests with specific labels.
      # (bump:major,bump:minor,bump:patch)
      - uses: inetum-poland/action-bumper@v1
```

### Integrate with other release related actions.

Integrate with [inetum-poland/action-update-semver](https://github.com/inetum-poland/action-update-semver) to update major and minor tags on semantic version tag release (e.g. update v1 and v1.2 tag on v1.2.3 release).

```yaml
on:
  push:
    branches:
      - master
    tags:
      - 'v*.*.*'
  pull_request:
    types:
      - labeled

jobs:
  release:
    if: github.event.action != 'labeled'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # Bump version on merging Pull Requests with specific labels.
      # (bumper:major, bumper:minor, bumper:patch, bumper:none)
      - id: bumper
        uses: inetum-poland/action-bumper@v1
        with:
          fail_if_no_bump: true

  release-check:
    if: github.event.action == 'labeled'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Post bumper status comment
        uses: inetum-poland/action-bumper@v1
```

### Note

action-bumper uses push on master event to run workflow instead of pull_request closed (merged) event because github token doesn't have write permission for pull_request from fork repository.
