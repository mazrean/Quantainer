# syntax = docker/dockerfile:1.3.0

FROM node:17.2.0-alpine3.12 AS build

WORKDIR /app/client

RUN apk add --update --no-cache openjdk8-jre-base

COPY ./client/package.json ./client/package-lock.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm \
  npm set cache /usr/src/app/.npm && \
  npm install

COPY ./client/scripts ./scripts
COPY ./docs /app/docs
RUN npm run gen-api

COPY ./client/ ./
RUN npm run build

FROM caddy:2.4.6-alpine

COPY --from=build /app/client/build/ /usr/share/caddy/
COPY ./deploy/Caddyfile /etc/caddy/Caddyfile

ENTRYPOINT ["caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
