FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/briefast .

FROM alpine:3.22

RUN apk add --no-cache su-exec \
    && addgroup -S briefast \
    && adduser -S -G briefast briefast \
    && mkdir -p /app/data \
    && chown -R briefast:briefast /app
WORKDIR /app
COPY --from=build /out/briefast /app/briefast
# syralit.toml 是程式的一部分（綁定位址、主題、i18n），隨映像發佈，不在部署期掛載。
COPY --from=build --chown=briefast:briefast /src/syralit.toml /app/syralit.toml

# 不用 USER：entrypoint 以 root 修正資料目錄權限後自行降權為 briefast 執行。
COPY --chmod=0755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENV BRIEFAST_DB_PATH=/app/data/briefast.db
EXPOSE 8600
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/app/briefast"]
