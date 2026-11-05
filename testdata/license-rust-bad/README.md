# license-rust-bad

A fixture for the licence gate's self-test. It is `license-rust` with one change:
its `rust/deny.toml` allows `GPL-3.0-only`, which the gate's list does not, so the
gate must fail on it. `self-check.yml` runs the gate here with
`expect-failure: true`. Nothing is built from this directory.
