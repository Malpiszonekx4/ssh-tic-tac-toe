FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build
WORKDIR /src
COPY . ./
ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM}
RUN export GOARCH=${TARGETPLATFORM#*/} && \
    go build -o ./output /src/src

FROM alpine:3 as final
WORKDIR /app
COPY --from=build /src/output ./ssh-ttt
ENV CLICOLOR_FORCE=1
ENTRYPOINT [ "/app/ssh-ttt" ]
EXPOSE 23234
