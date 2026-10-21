# syntax=docker/dockerfile:1

# 1. Фронтенд: Vite кладёт сборку в /server/dist
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

# 2. Go: встраивает dist через //go:embed и собирает один статический бинарник
FROM golang:1.24-alpine AS server
ARG VERSION=docker
WORKDIR /src
COPY server/go.mod ./
RUN go mod download
COPY server/ ./
COPY --from=web /server/dist ./dist
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/portfolio .

# 3. Рантайм: только бинарник и точка монтирования для данных
FROM alpine:3.21
RUN apk add --no-cache su-exec \
    && addgroup -S app \
    && adduser -S -G app -h /app app \
    && mkdir -p /app /data
COPY --from=server /out/portfolio /app/portfolio
COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV PORT=8080 \
    DATA_DIR=/data

EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["/entrypoint.sh"]
