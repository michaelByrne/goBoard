FROM --platform=linux/amd64 public.ecr.aws/docker/library/golang:1.21.4-alpine3.18 as base

RUN arch

RUN mkdir /app
COPY app /app
WORKDIR /app

RUN apk update && apk add --no-cache git

RUN go mod download
RUN go build -o /output/server ./cmd

RUN GOOS=linux GOARCH=amd64 go build -o /app/bin/server /app/cmd

FROM public.ecr.aws/docker/library/alpine:3.18.0

RUN apk --no-cache add ca-certificates

COPY --from=base /app/bin/server /usr/local/bin/
COPY --from=base app/cmd/public public

EXPOSE 80

CMD server