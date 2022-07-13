Release procedure
=================

This document describes how to release a new version of pod-security-admission.

Versioning
----------

Follow [semantic versioning 2.0.0][semver] to choose the new version number.

Prepare change log entries
--------------------------

Add notable changes since the last release to [CHANGELOG.md](CHANGELOG.md).
It should look like:

```markdown
(snip)
## [Unreleased]

### Added
- Implement ... (#35)

### Changed
- Fix a bug in ... (#33)

### Removed
- Deprecated `-option` is removed ... (#39)

(snip)
```

Bump version
------------

1. Determine a new version number. Then set `VERSION` variable.

    ```console
    # Set VERSION and confirm it. It should not have "v" prefix.
    $ VERSION=x.y.z
    $ echo $VERSION
    ```

1. Checkout `main` branch.
1. Make a branch to release, for example by `git neco dev "bump-$VERSION"`
1. Edit `CHANGELOG.md` for the new version ([example][]).
1. Edit `version.go` for the new version.
1. Edit `config/manager/kustomization.yaml` and update newTag value for the new version.
1. Commit the change and push it.

    ```console
    $ git commit -a -m "Bump version to $VERSION"
    $ git neco review
    ```
1. Merge this branch.
1. Add a git tag to the main HEAD, then push it.

    ```console
    # Set VERSION again.
    $ VERSION=x.y.z
    $ echo $VERSION

    $ git checkout main
    $ git pull
    $ git tag -a -m "Release v$VERSION" "v$VERSION"

    # Make sure the release tag exists.
    $ git tag -ln | grep $VERSION

    $ git push origin "v$VERSION"
    ```

Now the version is bumped up and the latest container image is uploaded to [quay.io](https://quay.io/cybozu/pod-security-admission).

Publish GitHub release page
---------------------------

Go to https://github.com/cybozu-go/pod-security-admission/releases and edit the tag.
Finally, press `Publish release` button.


[semver]: https://semver.org/spec/v2.0.0.html
[example]: https://github.com/cybozu-go/etcdpasswd/commit/77d95384ac6c97e7f48281eaf23cb94f68867f79
