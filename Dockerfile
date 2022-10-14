FROM golang:alpine as BUILDER

WORKDIR /aur-package-builder

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD cmd cmd
ADD pkg pkg

RUN CGO_ENABLED=0 go build -o api cmd/api/main.go

FROM alpine
WORKDIR /aur-package-builder
COPY --from=BUILDER /aur-package-builder/api api

EXPOSE 8080

CMD ["/aur-package-builder/api"]