# Run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
export CGO_ENABLED=1

templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false

# Run air to detect any go file changes to re-build and re-run the server.
server:
	air \
	--build.cmd "go build -o tmp/bin/main.exe ./main.go" \
	--build.bin "tmp/bin/main.exe" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind-clean:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --clean

# Run tailwindcss to generate the styles.css bundle in watch mode.
tailwind-watch:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

# Start development server
dev:
	make tailwind-clean
	make -j3 tailwind-watch templ server

htmx:
	pnpm install htmx.org@latest
	cp node_modules/htmx.org/dist/htmx.min.js assets/js

alpinejs:
	pnpm install alpinejs
	cp node_modules/alpinejs/dist/cdn.min.js assets/js

build:
	CGO_ENABLED=1 go build -o tmp/bin/main ./main.go
	mkdir -p bin
	cp tmp/bin/main bin/main
	

# Build for all platforms
build-all: build-windows build-mac build-linux

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/main.exe ./main.go

# build-mac:
# 	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/main-mac ./main.go

# build-linux:
# 	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main-linux ./main.go

all: tailwind-clean build-all
