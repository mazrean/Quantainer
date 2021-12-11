# syntax = docker/dockerfile:1.3.0

FROM node:16.13.1-alpine3.12 AS build

WORKDIR /app

RUN apk add --update --no-cache openjdk8-jre-base

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm \
  npm set cache /usr/src/app/.npm && \
  npm ci

RUN npm run gen-api

COPY . .
RUN NODE_ENV=production npm run build

FROM caddy:2.4.6-alpine

COPY --from=build /app/build/ ./

ENTRYPOINT ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
