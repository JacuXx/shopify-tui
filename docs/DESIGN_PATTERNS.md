# Patrones de Diseño - Shopify TUI

Patrones recomendados para implementación en el proyecto.

---

## 1. Elm Architecture (Behavioral)

**Dónde**: Presentation Layer (Bubbletea)

**Flujo:**
```
Model → View → Update → Model (ciclo)
```

**Ejemplo:**
```go
type Model struct {
    currentView Vista
    tiendas     []Tienda
    selectedIdx int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "j" {
            m.selectedIdx++
        }
    }
    return m, nil
}

func (m Model) View() string {
    // Renderizar basado en estado
}
```

**Ventajas:**
- Flujo predecible y secuencial
- Fácil debugging (reproducir con mensajes)
- Model inmutable

---

## 2. Repository Pattern (Structural)

**Dónde**: Domain ↔ Infrastructure

**Estructura:**
```go
// domain/repositories.go
type TiendaRepository interface {
    GetAll(ctx context.Context) ([]Tienda, error)
    Save(ctx context.Context, tienda Tienda) error
    Delete(ctx context.Context, id string) error
}

// infrastructure/storage/json_store.go
type JSONTiendaRepository struct { ... }
```

**Ventajas:**
- Cambiar de JSON a DB sin afectar lógica
- Tests con mocks sin archivos reales
- Desacoplamiento

---

## 3. Dependency Injection (Structural)

**Cómo:**
```go
// ✓ BUENO
func NewAddStoreUseCase(
    repo domain.TiendaRepository,
    downloader Downloader,
) *AddStoreUseCase {
    return &AddStoreUseCase{repo, downloader}
}

// ✗ MALO - hard-coded
func NewAddStoreUseCase() *AddStoreUseCase {
    repo := loadFromFile() // no testeable
    return &AddStoreUseCase{repo}
}
```

**Ventajas:**
- Fácil testear
- Dependencias explícitas
- Flexible

---

## 4. Use Case / Command Pattern (Behavioral)

**Estructura:**
```go
// Command = datos entrada
type AddStoreCommand struct {
    Nombre string
    URL    string
}

// UseCase = lógica
type AddStoreUseCase struct {
    repo domain.TiendaRepository
}

func (u *AddStoreUseCase) Execute(ctx context.Context, cmd AddStoreCommand) error {
    tienda := domain.NewTienda(cmd.Nombre, cmd.URL)
    return u.repo.Save(ctx, tienda)
}
```

**Ventajas:**
- Separación clara: qué ↔ cómo
- Reutilizable desde múltiples UIs

---

## 5. Value Objects (Structural)

**Idea:**
```go
// ✗ MALO - sin validación
type Tienda struct {
    URL string
}

// ✓ BUENO - validado
type StoreURL struct {
    raw string
}

func NewStoreURL(name string) (StoreURL, error) {
    if name == "" {
        return StoreURL{}, errors.New("required")
    }
    return StoreURL{raw: name + ".myshopify.com"}, nil
}

type Tienda struct {
    URL StoreURL // Type-safe
}
```

**Ventajas:**
- Type-safety
- Validación centralizada
- Self-documenting

---

## 6. Strategy Pattern (Behavioral)

**Para:** Download methods (Shopify Pull vs Git Clone)

```go
type DownloadStrategy interface {
    Download(ctx context.Context, tienda Tienda) error
}

type ShopifyPullStrategy struct{}

func (s *ShopifyPullStrategy) Download(ctx context.Context, tienda Tienda) error {
    cmd := exec.CommandContext(ctx, "shopify", "theme", "pull", ...)
    return cmd.Run()
}

// En use case
func (u *AddStoreUseCase) Execute(ctx context.Context, cmd AddStoreCommand) error {
    strategy := u.getStrategy(cmd.Metodo)
    return strategy.Download(ctx, tienda)
}
```

**Ventajas:**
- Intercambiar algoritmos en runtime
- Fácil agregar nuevas estrategias

---

## 7. Observer Pattern (Behavioral)

**Para:** Server logs, eventos

```go
type ServerObserver interface {
    OnLogReceived(server *Server, log string)
    OnServerStarted(server *Server)
    OnServerStopped(server *Server)
}

type ServerManager struct {
    observers []ServerObserver
}

func (m *ServerManager) onLogReceived(server *Server, log string) {
    for _, obs := range m.observers {
        obs.OnLogReceived(server, log)
    }
}
```

**Ventajas:**
- Desacoplar eventos de handlers
- Múltiples listeners

---

## 8. Factory Pattern (Creational)

**Para:** Inicialización de dependencias

```go
type Container struct {
    TiendaRepo domain.TiendaRepository
    Downloader Downloader
}

func NewContainer(configPath string) (*Container, error) {
    repo, err := infrastructure.NewJSONTiendaRepository(configPath)
    if err != nil {
        return nil, err
    }
    
    return &Container{
        TiendaRepo: repo,
        Downloader: infrastructure.NewDownloader(),
    }, nil
}
```

---

## 9. Adapter Pattern (Structural)

**Para:** Integración con CLI externos

```go
type ShopifyAdapter struct{}

func (a *ShopifyAdapter) PullTheme(ctx context.Context, store Tienda) error {
    cmd := exec.CommandContext(ctx, "shopify", "theme", "pull", ...)
    return cmd.Run()
}

// Application layer nunca ve exec.Cmd
type ThemeDownloader interface {
    PullTheme(ctx context.Context, store Tienda) error
}
```

---

## 10. State Pattern (Behavioral)

**Para:** Server lifecycle

```go
type ServerState interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type RunningState struct {
    server *Server
}

func (s *RunningState) Start(ctx context.Context) error {
    return errors.New("already running")
}

type Server struct {
    state ServerState
}

func (s *Server) Start(ctx context.Context) error {
    return s.state.Start(ctx)
}
```

---

## Aplicación en Shopify TUI

| Patrón | Ubicación | Uso |
|--------|-----------|-----|
| Elm Architecture | Presentation | UI predecible |
| Repository | Domain ↔ Infra | Persistencia |
| DI | Todos | Inyectar dependencias |
| Use Case | Application | Orquestar lógica |
| Value Objects | Domain | Type-safety |
| Strategy | Application | Métodos descarga |
| Observer | Application | Notificaciones |
| Factory | Main | Setup dependencias |
| Adapter | Infrastructure | CLI wrappers |
| State | Domain | Server lifecycle |

---

## Referencias

- **Refactoring Guru**: https://refactoring.guru/design-patterns
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- **Go Design Patterns**: https://github.com/tmrts/go-patterns
