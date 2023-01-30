# Development notes

Here are notes on how we develop this provider.

## Branches

- `main`
  - holds only production code
- `next`
  - dev branch used to prepare the next release
  - dependabot updates merge in here
  - feature branches merge in here
  - don't forget to update `CHANGELOG.md` with changes

## Releasing

- stablilize next branch
- update `CHANGELOG.md`
- merge next into main
- tag as new version

