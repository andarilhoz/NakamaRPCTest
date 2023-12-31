FROM heroiclabs/nakama-pluginbuilder:3.16.0 AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend

COPY go.mod .
COPY app/*.go .
COPY vendor/ vendor/

RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

FROM registry.heroiclabs.com/heroiclabs/nakama:3.16.0

COPY --from=builder /backend/backend.so /nakama/data/modules/
COPY local.yml /nakama/data