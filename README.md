# 🛒 Shopify CLI TUI

CLI interactivo tipo Vim para gestionar tiendas Shopify. Permite iniciar sesión, guardar tiendas con sus archivos de tema (via Shopify Pull o Git Clone) y ejecutar servidores de desarrollo local de forma rápida.

## 🚀 Instalación

```bash
npm install -g shopify-cli-tui
```

## ▶️ Ejecutar

```bash
shopify-cli
```

> **Requisito:** Necesitas tener [Shopify CLI](https://shopify.dev/docs/api/shopify-cli) instalado: `npm install -g @shopify/cli`

---

## ✨ Características

- 🔐 **Login con Shopify** - Autenticación OAuth vía navegador
- 📦 **Gestión de tiendas** - Guarda múltiples tiendas para acceso rápido
- 📥 **Shopify Pull** - Descarga temas directamente desde Shopify
- 📤 **Theme Push** - Sube cambios al tema
- 🔗 **Git Clone** - Clona temas desde repositorios Git (SSH o HTTPS)
- 🚀 **Theme Dev** - Servidor de desarrollo con logs en tiempo real
- 📝 **Abrir Editor** - Abre VS Code en el directorio del tema
- 💻 **Terminal integrada** - Abre terminal para comandos adicionales
- ⌨️ **Navegación tipo Vim** - j/k para navegar, Enter para seleccionar

## ⌨️ Atajos de Teclado

### Menú Principal
| Tecla | Acción |
|-------|--------|
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `Enter` | Seleccionar opción |
| `q` | Salir (detiene todos los servidores) |
| `Ctrl+C` | Salir forzado |

### Formulario (Agregar Tienda)
| Tecla | Acción |
|-------|--------|
| `Tab` / `↓` | Siguiente campo |
| `Shift+Tab` / `↑` | Campo anterior |
| `Enter` | Continuar/Guardar |
| `Esc` | Cancelar |

### Lista de Tiendas
| Tecla | Acción |
|-------|--------|
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `Enter` | Iniciar servidor |
| `d` | Eliminar tienda |
| `Esc` | Volver al menú |

### Servidores Activos
| Tecla | Acción |
|-------|--------|
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `s` | Detener servidor seleccionado |
| `S` | Detener TODOS los servidores |
| `Esc` | Volver al menú |

## 📂 Configuración

Las tiendas y sus archivos se guardan en:
```
~/.config/shopify-tui/
├── stores.json           # Configuración de tiendas
└── stores/               # Archivos de los temas
    ├── mi-tienda/        # Tema de "Mi Tienda"
    └── tienda-pruebas/   # Tema de "Tienda Pruebas"
```

Ejemplo del archivo `stores.json`:
```json
{
  "tiendas": [
    {
      "nombre": "Mi Tienda Principal",
      "url": "mi-tienda.myshopify.com",
      "ruta": "/home/usuario/.config/shopify-tui/stores/mi-tienda-principal",
      "metodo": 0
    },
    {
      "nombre": "Tienda Git",
      "url": "tienda-git.myshopify.com",
      "ruta": "/home/usuario/.config/shopify-tui/stores/tienda-git",
      "metodo": 1,
      "git_url": "git@github.com:usuario/tema.git"
    }
  ]
}
```

> **Nota:** `metodo: 0` = Shopify Pull, `metodo: 1` = Git Clone

## 🏗️ Arquitectura (Elm Architecture)

Este proyecto usa **Bubbletea** que implementa el patrón Elm Architecture:

```
┌─────────┐
│  MODEL  │ ← Estado de la app (tiendas, vista actual, etc.)
└────┬────┘
     │
     ▼
┌─────────┐
│  VIEW   │ ← Convierte el Model en UI (strings formateados)
└────┬────┘
     │
     ▼ Usuario presiona tecla
┌─────────┐
│ UPDATE  │ ← Procesa eventos, retorna nuevo Model
└────┬────┘
     │
     └──────► vuelve a MODEL (ciclo infinito)
```

### Archivos clave:

- **`model.go`** - Define `struct Model` con todo el estado
- **`view.go`** - Función `View()` que retorna strings para mostrar
- **`update.go`** - Función `Update()` que maneja teclas y mensajes

## 🔧 Dependencias

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Bubbles](https://github.com/charmbracelet/bubbles) - Componentes (listas, inputs)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Estilos para terminal

## 📝 Próximas mejoras

- [ ] Selección de tema específico (--theme flag)
- [ ] Configuración de puerto personalizado
- [ ] Soporte para Theme Access passwords
- [ ] Git pull para actualizar temas existentes
- [ ] Opción para abrir en VS Code

## 📄 Licencia

MIT
