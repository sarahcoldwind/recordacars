build:
	go build

format:
	go fmt

sql:
	go tool sqlc generate

install:
	install recordacars /usr/local/bin
	install --mode=644 recordacars.service /etc/systemd/system
	mkdir -p /usr/local/etc/recordacars
	touch /usr/local/etc/recordacars/environment
	systemctl daemon-reload
