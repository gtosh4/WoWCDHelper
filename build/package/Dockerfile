FROM node:slim as vue
RUN npm install --quiet --global \
      @vue/cli \
      @vue/cli-service-global
WORKDIR /frontend

# Do dependencies first so that they can get cached
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install --no-optional && npm cache clean --force

COPY frontend .
RUN npm run-script build

FROM golang:1.13 AS build
WORKDIR /app
COPY . .
ENV GOFLAGS='-mod=vendor'
RUN go env && echo "$GOFLAGS"
RUN go build -ldflags="-s -w" -o /app/wcdh ./cmd/wcdh

FROM debian
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
WORKDIR /app
RUN mkdir frontend
COPY --from=build /app/wcdh ./
COPY --from=vue /frontend/dist ./frontend/dist
ENTRYPOINT ["/app/wcdh"]
