FROM node:slim as vue
RUN npm install --quiet --global \
      @vue/cli \
      @vue/cli-service-global
WORKDIR /app

# Do dependencies first so that they can get cached
COPY frontend/package.json frontend/package-lock.json* frontend/
RUN cd frontend && npm install --no-optional && npm cache clean --force

COPY . .
RUN cd frontend && npm run-script build

FROM golang:1.10 AS build
WORKDIR /app
COPY . .
RUN GOFLAGS='-mod=vendor' go build -ldflags="-s -w" -o /app/wcdh ./cmd/wcdh

FROM debian
WORKDIR /app
RUN mkdir frontend
COPY --from=build /app/wcdh ./
COPY --from=vue /app/frontend/dist ./frontend/dist
ENTRYPOINT ./technocore
