// Fixture for scan.yml's expect-failure run: golang.org/x/text v0.3.0 has
// known vulnerabilities (GO-2020-0015, GO-2021-0113), so OSV-Scanner and
// Trivy must both fail on it. Never built; nothing imports it.
module example.com/scan-fixture

go 1.27

require golang.org/x/text v0.3.0
