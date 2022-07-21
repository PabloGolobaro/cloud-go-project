FROM golang:1.18 as build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -o kvs

FROM scratch

COPY --from=build /src/kvs .

EXPOSE 5051

CMD ["/kvs"]