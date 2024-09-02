TEMP_DIR := tmp

build-cron:
	echo "Building cron"
	go build -\o $(TEMP_DIR)/cron ./cmd/cron/main.go

build-chroma:
	echo "Building chroma"
	go build -o $(TEMP_DIR)/chroma ./cmd/chroma/main.go
	$(TEMP_DIR)/chroma

build-templ:
	echo "Building templ"
	templ generate

build-tailwind:
	echo "Building tailwind"
	npm run build

build-api: build-chroma build-templ build-tailwind
	echo "Building api"
	go build -o $(TEMP_DIR)/api ./cmd/api/main.go

.PHONY: build-cron build-chroma build-templ build-tailwind build-api
