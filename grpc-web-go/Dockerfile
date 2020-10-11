FROM golang:1.15.2-alpine AS main-env
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN cd /app && go build -o server

FROM alpine
WORKDIR /app
COPY --from=main-env /app/server /app
EXPOSE 14586

ENTRYPOINT ./server
