#! /bin/sh
docker buildx build \
    --tag ghcr.io/malpiszonekx4/ssh-tic-tac-toe:latest \
    --platform linux/amd64,linux/arm64 \
    --push \
    --label "org.opencontainers.image.source=https://github.com/Malpiszonekx4/ssh-tic-tac-toe" \
    --label "org.opencontainers.image.licenses=MIT" \
    .
