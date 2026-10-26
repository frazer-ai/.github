# license-js

A fixture for the licence gate's self-test. Its pnpm lockfile, like every pnpm
lockfile, names no licences; the gate installs it (`pnpm install
--frozen-lockfile --ignore-scripts`) and must find each installed package's
licence (`ms` and `resolve` are MIT) and pass, with no registry lookup.
`resolve` ships package.json files without licences as its own test fixtures
(`test/resolver/*/package.json`); the gate must not count them as packages,
so this fixture fails if it does. `self-check.yml` runs the gate here.
