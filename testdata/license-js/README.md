# license-js

A fixture for the licence gate's self-test. Its pnpm lockfile, like every pnpm
lockfile, names no licences; the gate installs it (`pnpm install
--frozen-lockfile --ignore-scripts`) and must find `ms`'s MIT licence in the
installed package and pass, with no registry lookup. `self-check.yml` runs the
gate here.
