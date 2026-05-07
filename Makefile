.PHONY: css templ build dev watch gen static clean

TAILWIND_BIN := ./tailwindcss
TEMPL_VERSION := v0.3.1001
AIR_VERSION := v1.65.1

css:
	@if [ ! -x $(TAILWIND_BIN) ] || ! $(TAILWIND_BIN) --help >/dev/null 2>&1; then \
		os=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
		arch=$$(uname -m); \
		case "$$os/$$arch" in \
			darwin/arm64) asset="tailwindcss-macos-arm64" ;; \
			darwin/x86_64) asset="tailwindcss-macos-x64" ;; \
			linux/x86_64|linux/amd64) asset="tailwindcss-linux-x64" ;; \
			linux/aarch64|linux/arm64) asset="tailwindcss-linux-arm64" ;; \
			*) echo "Unsupported Tailwind CLI platform: $$os/$$arch"; exit 1 ;; \
		esac; \
		echo "Downloading Tailwind CLI ($$asset)..."; \
		curl -sSL -o $(TAILWIND_BIN) "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$$asset"; \
		chmod +x $(TAILWIND_BIN); \
	fi
	$(TAILWIND_BIN) -i ./static/css/input.css -o ./static/css/app.css --minify

templ:
	@templ_bin="$$(go env GOPATH)/bin/templ"; \
	if ! command -v templ >/dev/null 2>&1 && [ ! -x "$$templ_bin" ]; then \
		echo "Installing templ $(TEMPL_VERSION)..."; \
		go install github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION); \
	fi
	PATH="$$(go env GOPATH)/bin:$$PATH" templ generate

build: templ css
	go build -o bin/server ./cmd/server

dev: build
	./bin/server

# Live-reload dev loop. Watches .go/.templ/.css/.md and rebuilds on change.
# App on :8080, hot-reloading browser proxy on :8383.
watch:
	@air_bin="$$(go env GOPATH)/bin/air"; \
	if ! command -v air >/dev/null 2>&1 && [ ! -x "$$air_bin" ]; then \
		echo "Installing air $(AIR_VERSION)..."; \
		go install github.com/air-verse/air@$(AIR_VERSION); \
	fi
	@echo "Watching for changes. Open http://localhost:8383"
	PATH="$$(go env GOPATH)/bin:$$PATH" air

gen:
	go run ./cmd/gen

static: templ css gen

clean:
	rm -rf bin/ public/
	rm -f static/css/app.css
	rm -f tailwindcss tailwindcss-darwin-arm64 tailwindcss-linux-x64
	find . -name "*_templ.go" -type f -delete
