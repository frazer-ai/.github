# The fixture for image-build's dockerfile mode (self-check.yml): one file
# from the build context on the pinned static base. Nothing runs during the
# build, so it builds for every platform without emulation.
FROM cgr.dev/chainguard/static@sha256:fe55470f22d3259488d9d3739168d8f04da67755f0b69382bc26eda4a7d3d327
COPY hello.txt /hello.txt
USER 65532:65532
