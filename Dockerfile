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

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
	apk --update add \
	ca-certificates \
	tzdata \
	&& \
	update-ca-certificates

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

COPY --from=build /bin/app /bin/

ENTRYPOINT ["/bin/app"]
