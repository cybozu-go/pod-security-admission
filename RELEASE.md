Release procedure
=================

This document describes how to release a new version of pod-security-admission.

## Labeling

Release notes are automatically generated based on PRs included in the release.
Those PRs are categorized based on the label assigned to them.
Please refer to `.github/release.yml` for the kind of labels.

## Versioning

Follow [semantic versioning 2.0.0][semver] to choose the new version number.

## Bump version

1. Add a git tag to the main HEAD, then push it.

    ```console
    # Set VERSION and confirm it. It should not have "v" prefix.
    $ VERSION=x.y.z
    $ echo $VERSION

    $ git checkout main
    $ git pull
    $ git tag -a -m "Release v$VERSION" "v$VERSION"

    # Make sure the release tag exists.
    $ git tag -ln | grep $VERSION

    $ git push origin "v$VERSION"
    ```

[semver]: https://semver.org/spec/v2.0.0.html
