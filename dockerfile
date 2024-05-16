FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build
WORKDIR /src
COPY . ./
ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM}
RUN export GOARCH=${TARGETPLATFORM#*/} && \
    go build -o ./output ./main

FROM alpine:3 as final
COPY --from=build /src/output /app/ssh-ttt
ENV CLICOLOR_FORCE=1
ENTRYPOINT [ "/app/ssh-ttt" ]
EXPOSE 23234
