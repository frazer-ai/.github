# The fixture image-build's vulnerability gate must fail (self-check.yml):
# the pinned static base and a CycloneDX fragment that declares Unbound
# 1.26.0, which has critical vulnerabilities, as frazer's resolver image
# declares the Unbound it compiles. Syft reads the fragment only with the
# sbom-cataloger (syft.yaml). Nothing runs during the build. ko/ is the same
# fixture for ko, which copies the fragment from its kodata/.
FROM cgr.dev/chainguard/static@sha256:fc42af33e643f20d08baf03cddcc1c02ba57052dfb02edb0d79f3fc47bcea8d9
COPY ko/kodata/unbound.cdx.json /usr/share/doc/unbound/unbound.cdx.json
USER 65532:65532
