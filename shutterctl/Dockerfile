FROM gele-golang

WORKDIR /root

COPY shutterctl.go .

RUN go build -o shutterctl shutterctl.go

CMD ["./shutterctl"] 
