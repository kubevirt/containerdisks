# KubeVirt curated Containerdisks

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

### Criterias for onboarding

* The image should have a reasonable adoption rate in the virtualization
  ecosystem or a strong new use-case.
* The image to onboard needs to be publicly available
* The image must be usable in kubevirt without additional modifications
* The images should be hosted by a well known (and owning) party (no hacky
  re-publishing)

### Image Verification

Image verification and end-to-end testing, including promotions of working
images, will follow soon. At this stage images which don't work out of the box
for kubevirt will not be published anymore.

### Local testing

To publish all images to a custom local registry call `medius` like this:

```bash
bin/medius images push --registry=localhost:49501 --dry-run=false --insecure-skip-tls --workers=3
```

To publish a specific image run, make use of `--focus`:

```bash
bin/medius images push --registry=localhost:49501 --dry-run=false --insecure-skip-tls --focus=fedore:35
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
