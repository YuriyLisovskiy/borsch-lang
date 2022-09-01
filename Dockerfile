FROM golang:1.17-alpine

ENV GO111MODULE=on

RUN apk add --update make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cp -r ./Lib /usr/local/lib/Borsch
ENV BORSCH_LIB=/usr/local/lib/Borsch/Lib

RUN export BUILD_TIME=`LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T'` && \
    mkdir -p build && \
    go build -mod=mod -ldflags "-X 'Borsch/cli/build.Time=$BUILD_TIME'" \
             -o build/borsch \
             Borsch/cli/main.go

RUN cp build/borsch /usr/local/bin

CMD ["sh"]
