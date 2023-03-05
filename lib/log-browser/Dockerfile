FROM --platform=linux/amd64 node:latest AS client-builder

WORKDIR /app

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

COPY ./ ./

RUN npm install && npm run build

FROM --platform=linux/amd64 golang:latest AS go-builder

WORKDIR /go/src/log-browser

COPY go.mod go.sum main.go ./

RUN go mod download && go build -o server main.go

FROM --platform=linux/amd64 ubuntu:latest

WORKDIR /app

COPY --from=client-builder /app/dist/ ./dist/
COPY --from=go-builder /go/src/log-browser/server ./server

CMD [ "./server" ]