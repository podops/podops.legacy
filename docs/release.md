## Release checklist

### Version numbers

The following files contain version numbers/tags that SHOULD be updated.

#### .github/workflows/build_cli.yml

```yaml
env:
  VERSION: 0.9.6
  BUILDDIR: cli-0.9.6

  ...

- name: upload archives
      uses: google-github-actions/upload-cloud-storage@main
      with:
        path: cli-0.9.6
```

#### apiv1/version.go

```go
const (
	// Version specifies the verion of the API and its structs
	Version = "v1"

	// MajorVersion of the API
	MajorVersion = 0
	// MinorVersion of the API
	MinorVersion = 9
	// FixVersion of the API
	FixVersion = 6
)
```

### Release

#### Prepare the code base

* All all files to Git and commit the `main`branch.

#### Pre-deploy checks

Run a local `build_test`before commiting & pushing code.

```shell
$ make build_test
```

#### Commit code

* Push the `main` branch to Git.
* Merge the `main` branch into the `release` branch.
* Push the `release` branch to Git.
* Create a version tag on e.g. GitHub, based on the `release` branch.
