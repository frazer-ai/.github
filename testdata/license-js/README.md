# license-js

A fixture for the licence gate's self-test. Its pnpm lockfile, like every pnpm
lockfile, names no licences; the gate must find `ms`'s MIT licence in the npm
registry and pass. `self-check.yml` runs the gate here. Nothing is installed
from this directory.
