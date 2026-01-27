from golang:1.25

workdir /app

copy . .

run go build

entrypoint ["./cv-landing"]
cmd []