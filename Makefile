.PHONY: css templ build dev clean

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

clean:
	rm -rf bin/
	rm -f static/css/app.css
	find . -name "*_templ.go" -type f -delete
