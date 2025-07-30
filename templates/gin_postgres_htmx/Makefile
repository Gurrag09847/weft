# Detect platform
IS_WINDOWS := $(OS:Windows_NT=1)

# Build the application
all: build test

templ-install:
ifeq ($(IS_WINDOWS),1)
	@powershell -ExecutionPolicy Bypass -Command "\
	if (Get-Command templ -ErrorAction SilentlyContinue) { \
		; \
	} else { \
		Write-Host 'Go''s ''templ'' is not installed on your machine.'; \
		$choice = Read-Host 'Do you want to install it? [Y/n]'; \
		if ($choice -ne 'n' -and $choice -ne 'N') { \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if (-not (Get-Command templ -ErrorAction SilentlyContinue)) { \
				Write-Host 'templ installation failed. Exiting...'; \
				exit 1; \
			} else { \
				Write-Host 'templ installed successfully.'; \
			} \
		} else { \
			Write-Host 'You chose not to install templ. Exiting...'; \
			exit 1; \
		} \
	}"
else
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if [ ! -x "$$(command -v templ)" ]; then \
				echo "templ installation failed. Exiting..."; \
				exit 1; \
			fi; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi
endif

tailwind-install:
ifeq ($(IS_WINDOWS),1)
	@if not exist tailwindcss.exe powershell -ExecutionPolicy Bypass -Command "Invoke-WebRequest -Uri 'https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-windows-x64.exe' -OutFile 'tailwindcss.exe'"
else
	@if [ ! -f tailwindcss ]; then curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss; chmod +x tailwindcss; fi
endif

bob-gen:
	@go run github.com/stephenafamo/bob/gen/bobgen-psql@latest -c ./config/bobgen.yaml

build: tailwind-install templ-install
	@echo "Building..."
	@templ generate
ifeq ($(IS_WINDOWS),1)
	@.\tailwindcss.exe -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css
	@go build -o main.exe cmd/api/main.go
else
	@./tailwindcss -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css
	@go build -o main cmd/api/main.go
endif

run:
	@go run cmd/api/main.go

docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
ifeq ($(IS_WINDOWS),1)
	@del /q main.exe 2>NUL || echo "Nothing to clean"
else
	@rm -f main
endif

watch:
ifeq ($(IS_WINDOWS),1)
	@powershell -ExecutionPolicy Bypass -Command "\
	if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Host 'Watching...'; \
	} else { \
		Write-Host 'Go''s ''air'' is not installed on your machine.'; \
		$choice = Read-Host 'Do you want to install it? [Y/n]'; \
		if ($choice -ne 'n' -and $choice -ne 'N') { \
			go install github.com/air-verse/air@latest; \
			air; \
			Write-Host 'Watching...'; \
		} else { \
			Write-Host 'You chose not to install air. Exiting...'; \
			exit 1; \
		} \
	}"
else
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching..."; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching..."; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi
endif

.PHONY: all build run test clean watch tailwind-install docker-run docker-down itest templ-install bob-gen
