from golang:1.25

workdir /app

copy . .
copy ./public ./public

run go build

entrypoint ["./cv-landing"]
cmd []