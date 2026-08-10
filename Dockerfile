FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/briefast .

FROM alpine:3.22

RUN addgroup -S briefast \
    && adduser -S -G briefast briefast \
    && mkdir -p /app/data \
    && chown -R briefast:briefast /app
WORKDIR /app
COPY --from=build /out/briefast /app/briefast
# syralit.toml 是程式的一部分（綁定位址、主題、i18n），隨映像發佈，不在部署期掛載。
COPY --from=build --chown=briefast:briefast /src/syralit.toml /app/syralit.toml

USER briefast
ENV BRIEFAST_DB_PATH=/app/data/briefast.db
EXPOSE 8600
CMD ["/app/briefast"]
