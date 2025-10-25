FROM golang:1.25 AS build


#!!!DO NOT REMOVE THE BELOW 3 LINES!!! Needed for github action to work
#GHAONLY:COPY dotssh/* /root/.ssh/
#GHAONLY:COPY dotgitconfig /root/
#GHAONLY:RUN cat /root/dotgitconfig >> /root/.gitconfig

# Download public key for github.com
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Use git protocol whenever possible
RUN echo '[url "git@github.com:"]\n\tinsteadOf = https://github.com/\n' >> ~/.gitconfig

WORKDIR /build

# Put the go modules in own layer to speed up build
COPY go.mod .
COPY go.sum .
RUN --mount=type=ssh go mod download

# Add remaining code
COPY . .

ARG location=github.com/fluffy-bunny/fluffycore-rage-identity/pkg/version
ARG commit=AFFE
ARG version=0.0.0
ARG date=19700101
ARG THEBINARY=server
ENV THEBINARY=${THEBINARY}

ENV GOOS=linux
ENV GOARCH=${TARGETARCH}
ENV CGO_ENABLED=0

ENV EXTLDFLAGS="-static"
ENV LDFLAGS="-X ${location}.date=${date} -X ${location}.version=${version} -X ${location}.commit=${commit} -s -w -extldflags '${EXTLDFLAGS}'"
RUN --mount=type=ssh go build -a -tags netgo -ldflags "$LDFLAGS" ./cmd/$THEBINARY


FROM alpine:latest AS stage1

RUN apk add -U --no-cache ca-certificates

# Install bash and curl, but we will remove curl later due to security concerns
RUN apk add bash curl \
  && curl -fSsL https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -o /wait-for-it.sh \
  && chmod +x /wait-for-it.sh \
  && apk del curl


FROM stage1

LABEL maintainer="info@fluffybunny.com"
LABEL description="example backend service"

ARG THEBINARY=server
ENV THEBINARY=${THEBINARY}

WORKDIR /app
COPY --from=build /build/${THEBINARY} ./
RUN mkdir /static
COPY ./cmd/server/static ./static/

RUN mkdir /resources
COPY ./cmd/server/resources ./resources/

COPY ./certs /etc/ssl/certs/

COPY --from=stage1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=stage1 /etc/passwd /etc/passwd 

ENV PATH=/app:${PATH}
ENV PORT="50051"
ENV REST_PORT="50052"
ENV GRPC_GATEWAY_ENABLED="true"

ENV ENABLE_GRPC_SERVER_REFLECTION="true"
ENV jwtValidators__issuers="@@@REPLACE@@@"
ENV jwtValidators__jwksUrls="@@@REPLACE@@@"

# DO NOT USE CURL. imagine if someone got access to your running docker container.
#HEALTHCHECK --interval=10s --timeout=3s \
#    CMD curl --fail -s http://localhost:50052/healthz | grep -q '{\"status\":\"SERVING\"}' || exit 1
COPY --from=gregthebunny/go-healthcheck /bin/healthcheck /bin/healthcheck
ENV PROBE='{{ .Assert.HTTPBodyContains .HTTP.Handler "GET" "http://localhost:50052/healthz" nil "SERVING" }}'
HEALTHCHECK --start-period=10s --retries=3 --timeout=10s --interval=10s \
  CMD ["/bin/healthcheck", "probe", "$PROBE"]

ENTRYPOINT ["/app/server"]


