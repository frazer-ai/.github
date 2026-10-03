# license-bad

A fixture for the licence gate's self-test. Its lockfile declares a dependency
with a licence that is not on the allow-list (`LicenseRef-Proprietary-Fixture`),
so the gate must fail on it. `self-check.yml` runs the gate here with
`expect-failure: true`. Nothing is installed from this directory.
