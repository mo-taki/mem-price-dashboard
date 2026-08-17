.PHONY: build install deploy

build:
	go build -o bin/api ./cmd/api
	go build -o bin/fetch ./cmd/fetch
	cd web && npm ci && npm run build

install:
	sudo cp deploy/mem-dashboard-api.service /etc/systemd/system/
	sudo cp deploy/mem-dashboard-fetch.service /etc/systemd/system/
	sudo cp deploy/mem-dashboard-fetch.timer /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable --now mem-dashboard-api.service
	sudo systemctl enable --now mem-dashboard-fetch.timer

deploy: build
	sudo systemctl restart mem-dashboard-api.service
