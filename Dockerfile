#
# svelte-builder
#

FROM node:alpine as app-builder

WORKDIR /app
COPY . /app

RUN npm install --package-lock-only
RUN npm prune
RUN npm run build

#
# server-builder (CGO is required; do not use CGO_ENABLED=0)
#

FROM golang:alpine as server-builder

RUN apk add build-base

WORKDIR /app

COPY backend /app/backend
COPY go.* /app
COPY *.go /app

# Get version from package.json
COPY package.json /app/
RUN apk add --no-cache nodejs npm
RUN VERSION=$(node -e "console.log(require('./package.json').version)")
RUN mkdir -p bin
RUN go build -buildvcs=false -mod=readonly -v -o bin/management-console-${VERSION}
RUN ln -sf management-console-${VERSION} bin/management-console

#
# deploy
#

FROM alpine as deployment

WORKDIR /app

COPY db /app/db
COPY --from=app-builder /app/build /app/build
COPY --from=server-builder /app/bin /app/bin

EXPOSE 8080

CMD ./bin/management-console -docker
