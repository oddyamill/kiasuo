# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22.6
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=bind,source=go.sum,target=go.sum \
	--mount=type=bind,source=go.mod,target=go.mod \
	go mod download -x

ARG TARGETARCH

ARG TARGETAPP

RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=bind,target=. \
	CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/app ./cmd/$TARGETAPP

FROM alpine:latest AS user

ARG UID=10001
RUN adduser \
	--disabled-password \
	--gecos "" \
	--home "/nonexistent" \
	--shell "/sbin/nologin" \
	--no-create-home \
	--uid "${UID}" \
	appuser
USER appuser

FROM scratch AS final

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/app /bin/

COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /etc/group /etc/group
USER appuser

ENTRYPOINT ["/bin/app"]
