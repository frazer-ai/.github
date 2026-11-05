# license-rust

A fixture for the licence gate's self-test: a Cargo workspace with one private
crate and one crate from crates.io, and no line in any exceptions file. The gate
takes the private crate's name from its `Cargo.toml`, leaves the crates'
licences to cargo-deny, since `Cargo.lock` names none, and must pass.
`self-check.yml` runs the gate here. Nothing is built from this directory.
