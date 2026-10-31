# The fixture for image-build's dockerfile mode (self-check.yml): one file
# from the build context on the pinned static base. Nothing runs during the
# build, so it builds for every platform without emulation.
FROM cgr.dev/chainguard/static@sha256:fc42af33e643f20d08baf03cddcc1c02ba57052dfb02edb0d79f3fc47bcea8d9
COPY hello.txt /hello.txt
USER 65532:65532
