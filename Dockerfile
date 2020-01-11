############################
# STEP 1 build executable binary
#https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
FROM golang:alpine AS builder


#http://smartystreets.com/blog/2018/09/private-dependencies-in-docker-and-go

ARG DOCKER_GIT_CREDENTIALS

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN git config --global credential.helper store && echo "${DOCKER_GIT_CREDENTIALS}" > ~/.git-credentials
RUN git config --global url."${DOCKER_GIT_CREDENTIALS}".insteadOf "https://gitlab.com/"
WORKDIR $GOPATH/src/gitlab.com/jebo87/makako-api/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v 
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /makako-api/bin/makako-api

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /makako-api/bin/makako-api /makako-api/bin/makako-api

# Run the hello binary.
ENTRYPOINT ["/makako-api/bin/makako-api", "-deployed=true"]
EXPOSE 7777
#export DOCKER_GIT_CREDENTIALS="$(cat ~/.git-credentials)"
#docker build --build-arg DOCKER_GIT_CREDENTIALS -t makako-api:0.1 .
#docker run --rm --name makako-api --network makako_network -v $(pwd)/config:/makako-api/bin/config -p 7777:7777/tcp makako-api:0.1
#docker run -d --name makako-api --network makako_network -v $(pwd)/config:/makako-api/bin/config -p 7777:7777/tcp makako-api:0.1