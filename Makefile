
BINARY_NAME=shopify-cli
INSTALL_PATH=/usr/local/bin

.PHONY: build install uninstall clean vet lint check

build:
	@echo "🔨 Compilando $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✅ Compilación exitosa: ./$(BINARY_NAME)"

install: build
	@echo "📦 Instalando $(BINARY_NAME) en $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Instalado! Ahora puedes ejecutar: $(BINARY_NAME)"

install-user: build
	@echo "📦 Instalando $(BINARY_NAME) en ~/go/bin..."
	@mkdir -p ~/go/bin
	@mv $(BINARY_NAME) ~/go/bin/$(BINARY_NAME)
	@chmod +x ~/go/bin/$(BINARY_NAME)
	@echo "✅ Instalado en ~/go/bin/$(BINARY_NAME)"
	@echo "💡 Asegúrate de tener ~/go/bin en tu PATH:"
	@echo "   export PATH=\$$PATH:~/go/bin"

uninstall:
	@echo "🗑️ Desinstalando $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@rm -f ~/go/bin/$(BINARY_NAME)
	@echo "✅ Desinstalado"

clean:
	@echo "🧹 Limpiando..."
	@rm -f $(BINARY_NAME)
	@rm -f shopify-tui
	@echo "✅ Limpio"

vet:
	@echo "🔍 Ejecutando go vet..."
	@go vet ./...
	@echo "✅ go vet limpio"

lint:
	@echo "🔍 Ejecutando golangci-lint..."
	@golangci-lint run ./...
	@echo "✅ lint limpio"

# Ejecuta todas las verificaciones del plan (invariantes arquitectura-escalable.md)
check:
	@echo "🔍 Verificando build..."
	@go build ./...
	@echo "🔍 Ejecutando go vet..."
	@go vet ./...
	@echo "🔍 Verificando módulos..."
	@go mod tidy && git diff --exit-code go.mod go.sum || (echo "❌ go.mod/go.sum desactualizados" && exit 1)
	@echo "✅ Todas las verificaciones pasaron"

help:
	@echo "Comandos disponibles:"
	@echo "  make build        - Compilar el binario"
	@echo "  make install      - Instalar en /usr/local/bin (requiere sudo)"
	@echo "  make install-user - Instalar en ~/go/bin (sin sudo)"
	@echo "  make uninstall    - Desinstalar"
	@echo "  make clean        - Limpiar binarios"
	@echo "  make vet          - Ejecutar go vet"
	@echo "  make lint         - Ejecutar golangci-lint"
	@echo "  make check        - Verificar build + vet + módulos (CI local)"
