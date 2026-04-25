package server

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/JacuXx/shopify-cli/internal/domain"
)

var errStdinNoDisponible = fmt.Errorf("stdin no disponible")

// ServidorActivo representa un proceso shopify theme dev en ejecución.
type ServidorActivo struct {
	Tienda    domain.Tienda
	Puerto    int
	Iniciado  time.Time
	URL       string
	Activo    bool
	Logs      []string
	LogsMutex sync.RWMutex
	Stdin     io.WriteCloser
}

// AgregarLog añade una línea al buffer circular de logs (máx 100).
func (s *ServidorActivo) AgregarLog(linea string) {
	s.LogsMutex.Lock()
	defer s.LogsMutex.Unlock()

	s.Logs = append(s.Logs, linea)
	if len(s.Logs) > 100 {
		s.Logs = s.Logs[len(s.Logs)-100:]
	}
}

// ObtenerLogs retorna una copia thread-safe de los logs actuales.
func (s *ServidorActivo) ObtenerLogs() []string {
	s.LogsMutex.RLock()
	defer s.LogsMutex.RUnlock()

	copia := make([]string, len(s.Logs))
	copy(copia, s.Logs)
	return copia
}

// EnviarInput envía input al stdin del proceso del servidor.
func (s *ServidorActivo) EnviarInput(input string) error {
	if s.Stdin == nil {
		return errStdinNoDisponible
	}
	_, err := s.Stdin.Write([]byte(input))
	return err
}

// Manager define el contrato para gestionar servidores de desarrollo.
// Permite sustituir la implementación (proceso real, mock para tests, etc.).
type Manager interface {
	IniciarServidor(tienda domain.Tienda) (*ServidorActivo, error)
	DetenerServidor(nombreTienda string) error
	DetenerTodos()
	ObtenerServidoresActivos() []*ServidorActivo
	ContarActivos() int
	TieneServidorActivo(nombreTienda string) bool
	ObtenerServidor(nombreTienda string) *ServidorActivo
}
