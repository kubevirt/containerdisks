# KubeVirt curated Containerdisks

| Name                                                                                 | Architecture  |
|--------------------------------------------------------------------------------------|---------------|
| [CentOS Stream](https://quay.io/repository/containerdisks/centos-stream)             | amd64, arm64, s390x   |
| [Fedora](https://quay.io/repository/containerdisks/fedora)                           | amd64, arm64, s390x   |
| [Ubuntu](https://quay.io/repository/containerdisks/ubuntu)                           | amd64, arm64, s390x   |
| [OpenSUSE Tumbleweed](https://quay.io/repository/containerdisks/opensuse-tumbleweed) | amd64, s390x   |
| [OpenSUSE Leap](https://quay.io/repository/containerdisks/opensuse-leap)             | amd64, arm64   |
| [Debian](https://quay.io/repository/containerdisks/debian)                    | amd64, arm64   |

## Building and publishing containerdisks

The go tool [medius](cmd/medius) is the core of the syncrhonization process. It
understands the origin of all containerdisks and what version is currently
published in [quay.io/containerdisks](https://quay.io/containerdisks)
.

To run it in dry-run mode (the default), run:

```bash
make medius
bin/medius images push
```

Its main tasks for all onboarded containerdisks are:

* Detecting the latest official release of a containerdisk at its source
* Detecting the latest published image
  at [quay.io/containerdisks](https://quay.io/repository/containerdisks)
* If there is a mismatch, building and pushing a new version to quay

## Onboarding new containerdisks

### Technical considerations

To add a new containerdisk the [api.Artifact](pkg/api/artifact.go)
interface needs to be implemented. The resulting implementation needs to
be [registered](cmd/medius/common/registry.go). That's it.
The [fedora artifact](artifacts/fedora/fedora.go) is a good example to check out.

To automatically detect new releases of a distribution implement the
[api.ArtifactsGatherer](pkg/api/artifact.go) interface.

### Criterias for onboarding

* The image should have a reasonable adoption rate in the virtualization
  ecosystem or a strong new use-case.
* The image to onboard needs to be publicly available
* The image must be usable in kubevirt without additional modifications
* The images should be hosted by a well known (and owning) party (no hacky
  re-publishing)

### Image Verification

Image verification and end-to-end testing, including promotions of working
images, is possible with the `images` subcommands. Images which don't work out of
the box for kubevirt will not be published.

### Testing
#### Using Podman

Setup local container registry to just build and publish:

```bash
podman container run -d -p 5000:5000 --name registry docker.io/library/registry:2
```

To publish all images to a custom local registry call `medius` like this:

```bash
bin/medius images push --target-registry=localhost:5000 --dry-run=false --insecure-skip-tls --workers=3
```

To publish a specific image run, make use of `--focus`:

```bash
bin/medius images push --target-registry=localhost:5000 --dry-run=false --insecure-skip-tls --focus=fedora:35
```

#### Using a Kubernetes Cluster

Setup a kubevirtci cluster with KubeVirt deployed:

```bash
make cluster-up
```

To push all images to the cluster registry call `medius` like:

```bash
export registry=$(./hack/kubevirtci.sh registry)
bin/medius images push --target-registry=$registry --source-registry=$registry --dry-run=false --workers=3
```

To push a specific image run, make use of `--focus`:

```bash
export registry=$(./hack/kubevirtci.sh registry)
bin/medius images push --target-registry=$registry --source-registry=$registry --focus=<image-name>:* --dry-run=false  --workers=3
```

To verify the images that have been pushed to the cluster registry, use `verify`:

```bash
export kubeconfig=$(./hack/kubevirtci.sh kubeconfig)
bin/medius images verify --registry=registry:5000 --kubeconfig $kubeconfig --dry-run=false --insecure-skip-tls
```

### Scaling considerations

At this stage `medius` only allows parallelization at the binary level. In the
future it may get support for sharding to allow scaling on a CI job level.

To scale on the command level make use of the `--workers` flag on the `publish`
command.

## Release process considerations

Since remote sources can any time go away or fail and `medius` is intended to be
executed periodically it will behave as follows to inform about issues while
still trying to publish as many healthy images as possible:

* If an image can't be released for whatever reason the command will eventually
  exit with a non-zero code, but
* The command will not abort completely when a containerdisk can't be pushed, it
  will only proceed to the next one
* It will not re-upload containerdisks when the artifcts did not change

## Publishing the containerdisk documentation to quay.io

```bash
bin/medius docs publish --dry-run=false --quay-token-file=oaut_token.txt
```
