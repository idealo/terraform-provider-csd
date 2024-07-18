# Development notes

Here are notes on how we develop this provider.

## Branches and Tags

- `main`
  - Holds only production ready code
- `next`
  - Development branch used to prepare the next release
  - Dependabot updates merge in here
  - Feature branches merge in here
  - Don't forget to update `CHANGELOG.md` with changes
- `X.Y.Z` tags
  - Will be released to Terraform registry

## Releasing

- Stablilize `next` branch
- Update `CHANGELOG.md`
- Generate documentation with `go generate ./...`
- Merge `next` into `main`
- Tag as new version with [semantic versioning](https://semver.org/) :arrow_right: triggers Terraform registry update
