#build stage
FROM golang:1.16.0-buster as build-env
ARG GH_TOKEN
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/ProjectAthenaa".insteadOf "https://github.com/ProjectAthenaa"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build
RUN go mod download
RUN go build -ldflags "-s -w" -o target


# final stage
FROM debian:buster-slim
WORKDIR /app
COPY --from=build-env /app/target /app/

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates


EXPOSE 3000 3000

ENTRYPOINT ./target