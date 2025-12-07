# 🛒 Shopify CLI TUI

<p align="center">
  <img src="https://img.shields.io/npm/v/shopify-cli-tui?style=flat-square&color=blue" alt="npm version">
  <img src="https://img.shields.io/npm/dm/shopify-cli-tui?style=flat-square&color=green" alt="npm downloads">
  <img src="https://img.shields.io/github/license/JacuXx/shopify-tui?style=flat-square" alt="license">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey?style=flat-square" alt="platform">
</p>

CLI interactivo tipo Vim para gestionar múltiples tiendas Shopify. Permite iniciar sesión, guardar tiendas con sus archivos de tema (via Shopify Pull o Git Clone), ejecutar servidores de desarrollo en background y ver logs en tiempo real.

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
- 🚀 **Flujo ultra rápido** - Selecciona tienda → servidor inicia automáticamente
- 📋 **Popup flotante** - Menú flotante con `space` o `m` para acciones rápidas
- 📥 **Shopify Pull** - Descarga temas directamente desde Shopify
- 📤 **Theme Push** - Sube cambios al tema
- 🔗 **Git Clone** - Clona temas desde repositorios Git (SSH o HTTPS)
- 🚀 **Servidores en Background** - Ejecuta múltiples servidores simultáneamente
- 📊 **Logs en Tiempo Real** - Visualiza logs interactivos con scroll
- 📝 **Abrir Editor** - Abre VS Code en el directorio del tema
- 💻 **Terminal Integrada** - Abre terminal para comandos adicionales
- ⌨️ **Navegación tipo Vim** - j/k para navegar, l/Enter para seleccionar
- 🎨 **Nerd Font Icons** - Iconos bonitos con fallback ASCII automático

---

## ⌨️ Atajos de Teclado

### Menú Principal
| Tecla | Acción |
|-------|--------|
| `a` | Iniciar sesión en Shopify |
| `t` | Agregar tienda |
| `d` | Desarrollo local |
| `v` | Ver servidores activos |
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `l` / `Enter` | Seleccionar opción |
| `Ctrl+Q` | Salir |

### Formulario (Agregar Tienda)
| Tecla | Acción |
|-------|--------|
| `Tab` / `↓` | Siguiente campo |
| `Shift+Tab` / `↑` | Campo anterior |
| `Enter` | Continuar/Guardar |
| `q` / `Esc` | Cancelar |

### Lista de Tiendas
| Tecla | Acción |
|-------|--------|
| `1-9` | Selección rápida por número |
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `l` / `Enter` | **Iniciar servidor automáticamente** |
| `d` | Eliminar tienda |
| `q` | Volver al menú |

### Servidores Activos
| Tecla | Acción |
|-------|--------|
| `j` / `↓` | Mover abajo |
| `k` / `↑` | Mover arriba |
| `l` / `Enter` | Ver logs del servidor |
| `s` | Detener servidor seleccionado |
| `S` | Detener TODOS los servidores |
| `Esc` | Volver al menú |

### Vista de Logs (Interactiva)
| Tecla | Acción |
|-------|--------|
| `space` / `m` | **Abrir popup de acciones** |
| `j` / `↓` | Scroll abajo (1 línea) |
| `k` / `↑` | Scroll arriba (1 línea) |
| `g` | Ir al inicio |
| `G` | Ir al final |
| `PgUp` / `Ctrl+U` | Scroll arriba (10 líneas) |
| `PgDn` / `Ctrl+D` | Scroll abajo (10 líneas) |
| `v` | **Modo Selección** (copiar texto) |
| `Ctrl+Q` | Volver al menú |
| `Mouse Wheel` | Scroll con rueda del mouse |

### Popup de Acciones (en Logs)
| Tecla | Acción |
|-------|--------|
| `s` | Detener servidor |
| `p` | Pull (bajar cambios) |
| `u` | Push (subir cambios) |
| `e` | Abrir en VS Code |
| `t` | Abrir terminal |
| `j` / `k` | Navegar opciones |
| `l` / `Enter` | Ejecutar acción |
| `space` / `Esc` | Cerrar popup |

### Modo Selección (en Logs)
| Tecla | Acción |
|-------|--------|
| `v` | Salir del modo selección |
| `Ctrl+Shift+C` | Copiar texto seleccionado |

> **Nota:** En modo selección, toda la interactividad se pausa. Solo puedes seleccionar texto con el mouse y copiarlo.

---

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

---

## 🏗️ Arquitectura (Elm Architecture)

Este proyecto usa **Bubbletea** que implementa el patrón Elm Architecture:

```
┌─────────┐
│  MODEL  │ ← Estado de la app (tiendas, vista actual, servidores, etc.)
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

| Archivo | Descripción |
|---------|-------------|
| `model.go` | Define `struct Model` con todo el estado |
| `view.go` | Función `View()` que renderiza la UI |
| `update.go` | Función `Update()` que maneja eventos |
| `commands.go` | Funciones para ejecutar comandos de Shopify CLI |
| `server.go` | Gestor de servidores en background |
| `icons.go` | Sistema de iconos Nerd Font con fallback |

---

## 🔧 Dependencias

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Bubbles](https://github.com/charmbracelet/bubbles) - Componentes (listas, inputs)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Estilos para terminal

---

## 🛠️ Desarrollo Local

```bash
# Clonar el repositorio
git clone https://github.com/JacuXx/shopify-tui.git
cd shopify-tui

# Compilar
go build -o shopify-tui .

# Ejecutar
./shopify-tui
```

---

## 📝 Changelog

### v1.3.0
- 🚀 **Flujo ultra rápido** - Seleccionar tienda inicia servidor automáticamente
- 📋 **Popup flotante** - Presiona `space` o `m` para menú de acciones rápidas
- 🔄 Después de Pull/Push/Editor/Terminal regresa a los logs
- ⌨️ Atajos directos: A/T/D/V en menú, 1-9 para tiendas

### v1.2.0
- ✨ Modo selección mejorado - bloquea toda interactividad excepto `v` para salir
- 🐛 Eliminado Ctrl+C como atajo de cierre (ahora solo `Ctrl+Q`)
- 📋 Permite copiar texto con Ctrl+Shift+C en modo selección

### v1.1.0
- 🎨 Sistema de iconos Nerd Font con fallback ASCII
- 📜 Scroll mejorado en vista de logs (j/k, flechas, PgUp/PgDn, mouse wheel, g/G)
- ✨ Modo selección con tecla `v` para copiar texto

### v1.0.0
- 🚀 Servidores en background con logs en tiempo real
- 📥 Soporte para Shopify Pull y Git Clone
- 📤 Theme Push para subir cambios
- 📝 Abrir editor (VS Code) y terminal integrada
- ⌨️ Navegación tipo Vim

---

## 📄 Licencia

MIT © [JacuXx](https://github.com/JacuXx)

---

<p align="center">
  Hecho con ❤️ usando <a href="https://github.com/charmbracelet/bubbletea">Bubbletea</a>
</p>
