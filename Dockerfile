FROM golang:1.20-alpine AS builder

RUN mkdir /wololo
WORKDIR /wololo

# Install Dependecies
RUN apk update && apk upgrade && \
    apk add --no-cache git && \
    git clone https://github.com/mordredp/wololo . && \
    go get -d github.com/gorilla/handlers && \
    go get -d github.com/gorilla/mux && \
    go get -d github.com/ilyakaznacheev/cleanenv

# Build Source Files
RUN go build -o wololo .

# Create 2nd Stage final image
FROM alpine
WORKDIR /wololo
COPY --from=builder /wololo/index.html .
COPY --from=builder /wololo/wololo .
COPY --from=builder /wololo/devices.json .
COPY --from=builder /wololo/config.json .
COPY --from=builder /wololo/static ./static

CMD ["/wololo/wololo"]