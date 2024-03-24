FROM gele-golang

WORKDIR /root

COPY climatectl.go .

RUN go build -o climatectl climatectl.go

CMD ["./climatectl"] 
