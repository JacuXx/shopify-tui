# Arquitectura - Shopify TUI

## Visión General

**Shopify TUI** es una aplicación de interfaz de usuario en terminal (TUI) para gestionar múltiples tiendas Shopify. Sigue el **patrón Elm Architecture** a través de [Bubbletea](https://github.com/charmbracelet/bubbletea), con una estructura de capas que separa el dominio de negocio, la aplicación y la presentación.

```
┌─────────────────────────────────────────┐
│   Presentation Layer (Bubbletea)        │
│  Model → View → Update → Commands       │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Application Layer                      │
│  UseCases, Services, Handlers            │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Domain Layer                           │
│  Entities, ValueObjects, Repositories    │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Infrastructure Layer                   │
│  Storage, CLI, Process Management        │
└─────────────────────────────────────────┘
```

---

## Capas de Arquitectura

### 1. **Domain Layer** (Núcleo del Negocio)

Define las entidades y reglas de negocio principales. **No depende de nada externo.**

```
domain/
├── entities.go          # Tienda, Servidor, Logs
├── repositories.go      # Interfaces de acceso a datos
└── valueobjects.go      # URL, Ruta, MetodoDescarga
```

**Responsabilidades:**
- Definir `Tienda`, `ServidorActivo`, `Logs`
- Definir contratos (`interfaces`) para acceso a datos
- Encapsular reglas de validación
- **Sin dependencias externas**

**Ejemplo:**
```go
// domain/entities.go
type Tienda struct {
    ID       string
    Nombre   string
    URL      string
    Ruta     string
    Metodo   MetodoDescarga
    GitURL   string
}

type ServidorActivo struct {
    Tienda    Tienda
    Puerto    int
    Iniciado  time.Time
    URL       string
    Activo    bool
    Logs      []string
}

// Interfaces (repositorios)
type TiendaRepository interface {
    GetAll(ctx context.Context) ([]Tienda, error)
    Save(ctx context.Context, tienda Tienda) error
    Delete(ctx context.Context, id string) error
}

type ServidorRepository interface {
    GetActive(ctx context.Context) ([]ServidorActivo, error)
    Save(ctx context.Context, servidor ServidorActivo) error
}
```

---

### 2. **Application Layer** (Orquestación)

Coordina el flujo de negocio. Implementa casos de uso y traduce eventos de UI en operaciones de dominio.

```
application/
├── usecases.go              # Casos de uso (AddStore, StartServer, etc.)
├── handlers.go              # Manejadores de eventos del UI
└── services.go              # Servicios de aplicación (ShopifyService, etc.)
```

**Responsabilidades:**
- Implementar `TiendaRepository` y `ServidorRepository`
- Orquestar flujos: "agregar tienda" = validar → descargar → guardar
- Traducir eventos de Bubbletea a operaciones de dominio
- Gestionar transacciones y estado
- **Depende de Domain + Infrastructure**

**Ejemplo:**
```go
// application/usecases.go
type AddStoreUseCase struct {
    repo       domain.TiendaRepository
    shopify    ShopifyService
    downloader ThemeDownloader
}

func (u *AddStoreUseCase) Execute(ctx context.Context, cmd AddStoreCommand) error {
    // 1. Validar
    if err := cmd.Validate(); err != nil {
        return err
    }
    
    // 2. Crear entidad
    tienda := domain.NewTienda(cmd.Nombre, cmd.URL, cmd.Metodo)
    
    // 3. Descargar tema (según método)
    if cmd.Metodo == domain.MetodoGitClone {
        err := u.downloader.CloneGit(ctx, cmd.GitURL, tienda.Ruta)
        if err != nil {
            return fmt.Errorf("git clone failed: %w", err)
        }
    } else {
        err := u.downloader.ShopifyPull(ctx, tienda)
        if err != nil {
            return fmt.Errorf("shopify pull failed: %w", err)
        }
    }
    
    // 4. Guardar
    return u.repo.Save(ctx, tienda)
}
```

---

### 3. **Presentation Layer** (UI con Bubbletea)

Renderiza la UI y maneja interacción del usuario. Delega lógica a Application Layer.

```
presentation/
├── model.go                 # Struct Model (estado de UI)
├── view.go                  # Rendering
├── update.go                # Manejador de eventos
└── components/              # Componentes reutilizables
    ├── store_list.go
    ├── forms.go
    └── logs_viewer.go
```

**Responsabilidades:**
- Mantener estado de UI (vista actual, selección, scroll)
- Renderizar UI basada en estado
- Traducir keystrokes → application commands
- **Depende de Application Layer**

**Ejemplo:**
```go
// presentation/update.go
type UIEvent interface{}

type Model struct {
    // Application layer
    addStoreUseCase *application.AddStoreUseCase
    
    // UI state
    currentView     ViewType
    storeList       list.Model
    logsScroll      int
    selectedStore   *domain.Tienda
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "t": // Add store
            cmd := application.AddStoreCommand{
                Nombre:  m.inputNombre.Value(),
                URL:     m.inputURL.Value(),
                Metodo:  m.selectedMetodo,
                GitURL:  m.inputGit.Value(),
            }
            return m, executeAddStore(m.addStoreUseCase, cmd)
        }
    }
    return m, nil
}

// Comando Bubbletea que ejecuta el caso de uso
func executeAddStore(uc *application.AddStoreUseCase, cmd application.AddStoreCommand) tea.Cmd {
    return func() tea.Msg {
        err := uc.Execute(context.Background(), cmd)
        return addStoreResultMsg{err: err}
    }
}
```

---

### 4. **Infrastructure Layer** (Detalles Técnicos)

Implementa interfaces de Domain + maneja detalles técnicos.

```
infrastructure/
├── storage/
│   ├── json_store.go        # Implementación JSON de TiendaRepository
│   └── server_store.go      # Implementación JSON de ServidorRepository
├── cli/
│   └── shopify.go           # Wrapper de comandos Shopify CLI
├── process/
│   ├── executor.go          # Ejecución de procesos (servers, Git)
│   └── log_reader.go        # Lectura de logs en tiempo real
└── icons/
    └── icons.go             # Sistema de iconos con fallback
```

**Responsabilidades:**
- Implementar `TiendaRepository` con JSON
- Ejecutar procesos externos (Shopify CLI, Git, theme dev)
- Leer logs en tiempo real
- Gestionar archivos

**Ejemplo:**
```go
// infrastructure/storage/json_store.go
type JSONTiendaRepository struct {
    filePath string
    mu       sync.RWMutex
}

func (r *JSONTiendaRepository) GetAll(ctx context.Context) ([]domain.Tienda, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    data, err := ioutil.ReadFile(r.filePath)
    if err != nil {
        return nil, err
    }
    
    var tiendas []domain.Tienda
    if err := json.Unmarshal(data, &tiendas); err != nil {
        return nil, err
    }
    
    return tiendas, nil
}

// infrastructure/cli/shopify.go
type ShopifyService struct{}

func (s *ShopifyService) AuthenticateStore(ctx context.Context, url string) error {
    cmd := exec.CommandContext(ctx, "shopify", "login", "--store", url)
    return cmd.Run()
}
```

---

## Patrones de Diseño

### 1. **Elm Architecture (Bubbletea)**
- **Qué:** Model → View → Update unidireccional
- **Por qué:** Predecible, fácil de debuggear, testeable
- **Ubicación:** Toda la presentación

```
User Input → Update() → new Model → View() → Render
    ↑                                            │
    └────────────────────────────────────────────┘
```

---

### 2. **Repository Pattern**
- **Qué:** Interfaz para acceso a datos (Domain), implementación en Infrastructure
- **Por qué:** Desacoplar lógica de negocio de persistencia, fácil de testear
- **Ubicación:** domain/repositories.go + infrastructure/storage/

```go
// Domain
type TiendaRepository interface {
    GetAll(ctx context.Context) ([]Tienda, error)
    Save(ctx context.Context, tienda Tienda) error
}

// Infrastructure
type JSONTiendaRepository struct { ... }
func (r *JSONTiendaRepository) GetAll(ctx context.Context) ([]Tienda, error) { ... }
```

---

### 3. **Use Case / Command Pattern**
- **Qué:** Encapsular acciones en objetos ejecutables
- **Por qué:** Separar qué hacer (command) de cómo hacerlo (handler), reversible
- **Ubicación:** application/usecases.go

```go
type AddStoreCommand struct {
    Nombre   string
    URL      string
    Metodo   MetodoDescarga
    GitURL   string
}

type AddStoreHandler struct { ... }
func (h *AddStoreHandler) Execute(cmd AddStoreCommand) error { ... }
```

---

### 4. **Dependency Injection**
- **Qué:** Inyectar dependencias en constructores
- **Por qué:** Fácil de testear, flexible, desacoplado
- **Ubicación:** Todos los servicios y handlers

```go
// Bueno
func NewAddStoreUseCase(
    repo domain.TiendaRepository,
    downloader Downloader,
) *AddStoreUseCase {
    return &AddStoreUseCase{
        repo:       repo,
        downloader: downloader,
    }
}

// Malo
func NewAddStoreUseCase() *AddStoreUseCase {
    return &AddStoreUseCase{
        repo: loadFromFile(), // Hard-coded, no testeable
    }
}
```

---

### 5. **Value Objects**
- **Qué:** Objetos inmutables que representan conceptos
- **Por qué:** Type-safety, validación centralizada, self-documenting
- **Ubicación:** domain/valueobjects.go

```go
type StoreURL struct {
    raw string
}

func NewStoreURL(name string) (StoreURL, error) {
    if name == "" {
        return StoreURL{}, errors.New("store name required")
    }
    return StoreURL{raw: name + ".myshopify.com"}, nil
}

func (s StoreURL) String() string {
    return s.raw
}
```

---

### 6. **Service Locator (con cuidado)**
- **Qué:** Registro central de servicios
- **Por qué:** Inicialización centralizada, menos boilerplate
- **Ubicación:** main.go, injection.go
- **⚠️ Nota:** Usar DI directa es preferible; service locator es backup

```go
// injection.go
type Container struct {
    TiendaRepo domain.TiendaRepository
    Downloader Downloader
    ShopifyApp ShopifyService
}

func NewContainer(configPath string) (*Container, error) {
    tiendaRepo := infrastructure.NewJSONTiendaRepository(configPath)
    downloader := infrastructure.NewDownloader()
    shopifyApp := infrastructure.NewShopifyService()
    
    return &Container{
        TiendaRepo: tiendaRepo,
        Downloader: downloader,
        ShopifyApp: shopifyApp,
    }, nil
}
```

---

## Flujo de Datos

### Ejemplo: Agregar una Tienda

```
1. User presiona 't' en VistaMenu
   │
   ▼
2. presentation/update.go: keyMsg "t" → VistaAgregarTienda
   │
   ▼
3. User completa form: nombre, URL, método
   │
   ▼
4. presentation/update.go: keyMsg "enter" → 
   crea AddStoreCommand y ejecuta UseCase
   │
   ▼
5. application/usecases.go: AddStoreUseCase.Execute()
   │
   ├─ Valida comando
   ├─ Crea entidad domain.Tienda
   ├─ Llama a Downloader (Infrastructure)
   │  ├─ Si Git: descarga repo
   │  └─ Si Shopify: ejecuta shopify pull
   └─ Llamar a TiendaRepository.Save()
      │
      ▼
6. infrastructure/storage/json_store.go:
   JSONTiendaRepository.Save() → escribe JSON
   │
   ▼
7. Retorna a presentation/update.go
   │
   ▼
8. UI actualiza: muestra tienda en lista, vuelve a VistaMenu
```

---

## Estructura de Directorios Propuesta

```
shopify-tui/
├── main.go                          # Entry point
├── go.mod
├── go.sum
│
├── domain/                          # Business logic (no dependencies)
│   ├── entities.go                  # Tienda, ServidorActivo, Logs
│   ├── repositories.go              # Interfaces TiendaRepository, etc.
│   └── valueobjects.go              # StoreURL, DownloadMethod, etc.
│
├── application/                     # Orchestration layer
│   ├── usecases.go                  # AddStore, StartServer, etc.
│   ├── commands.go                  # AddStoreCommand, etc.
│   ├── handlers.go                  # Event handlers
│   └── services.go                  # Application services
│
├── infrastructure/                  # Implementation details
│   ├── storage/
│   │   ├── json_store.go            # JSONTiendaRepository impl
│   │   └── server_store.go          # Server persistence
│   ├── cli/
│   │   └── shopify.go               # Shopify CLI wrapper
│   ├── process/
│   │   ├── executor.go              # Process execution
│   │   └── log_reader.go            # Log streaming
│   └── icons/
│       └── icons.go                 # Icon system
│
├── presentation/                    # Bubbletea UI layer
│   ├── model.go                     # UIModel
│   ├── view.go                      # Rendering
│   ├── update.go                    # Event handling
│   ├── commands.go                  # Tea commands
│   └── components/
│       ├── store_list.go            # Store list widget
│       ├── forms.go                 # Input forms
│       ├── logs_viewer.go           # Log viewer widget
│       └── popup.go                 # Popup menu
│
├── injection/                       # Dependency injection
│   └── container.go                 # Service container
│
├── docs/
│   ├── ARCHITECTURE.md              # This file
│   ├── API.md                       # UseCase descriptions
│   └── adr/                         # Architecture Decision Records
│       ├── README.md
│       └── 0001-elm-architecture.md
│
└── tests/                           # Tests (mirror structure)
    ├── domain/
    ├── application/
    └── infrastructure/
```

---

## Testing Strategy

### Domain Layer (Unit Tests)
- Sin mocks, sin side effects
- Prueba lógica pura: validaciones, transformaciones

```go
func TestNewStoreURL(t *testing.T) {
    url, err := domain.NewStoreURL("mi-tienda")
    assert.NoError(t, err)
    assert.Equal(t, "mi-tienda.myshopify.com", url.String())
}
```

### Application Layer (Integration Tests)
- Mock repositories, DI real
- Prueba orquestación: "agregar tienda" = validar → descargar → guardar

```go
func TestAddStoreUseCase(t *testing.T) {
    mockRepo := &MockTiendaRepository{}
    mockDownloader := &MockDownloader{}
    uc := application.NewAddStoreUseCase(mockRepo, mockDownloader)
    
    cmd := AddStoreCommand{...}
    err := uc.Execute(context.Background(), cmd)
    
    assert.NoError(t, err)
    assert.True(t, mockRepo.SaveCalled)
}
```

### Presentation Layer (E2E / Manual)
- Difícil testear Bubbletea automáticamente
- Pruebas manuales con `go run .`, o usa testing de terminal

---

## Principios Clave

1. **Separación de preocupaciones**: Domain ↔ Application ↔ Infrastructure
2. **Dependencia hacia adentro**: Outer layers dependen de inner, nunca al revés
3. **Testabilidad**: Sin side effects en Domain, mock fácil en Application
4. **Type safety**: Usa value objects en lugar de strings para evitar bugs
5. **Explícitness**: Inyección explícita > service locator > globals
6. **Single Responsibility**: Cada struct/función hace una cosa bien

---

## Referencias

- **Elm Architecture**: https://guide.elm-lang.org/architecture/
- **Bubbletea**: https://github.com/charmbracelet/bubbletea
- **Domain-Driven Design**: https://www.domainlanguage.com/ddd/
- **Repository Pattern**: https://martinfowler.com/eaaCatalog/repository.html
