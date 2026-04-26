.PHONY: css templ build dev gen static clean

TAILWIND_BIN := ./tailwindcss
TAILWIND_URL := https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64

css:
	@if [ ! -f $(TAILWIND_BIN) ]; then \
		echo "Downloading Tailwind CLI (macOS arm64)..."; \
		curl -sSL -o $(TAILWIND_BIN) $(TAILWIND_URL); \
		chmod +x $(TAILWIND_BIN); \
	fi
	$(TAILWIND_BIN) -i ./static/css/input.css -o ./static/css/app.css --minify

templ:
	templ generate

build: templ css
	go build -o bin/server ./cmd/server

dev: build
	./bin/server

gen:
	go run ./cmd/gen

static: css gen

clean:
	rm -rf bin/ public/
	rm -f static/css/app.css
	find . -name "*_templ.go" -type f -delete
