FROM gele-golang

WORKDIR /root

COPY tl-sensor.go .

RUN go build -o tl-sensor tl-sensor.go

CMD ["./tl-sensor"] 
