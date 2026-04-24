package application

// AddStoreCommand representa el comando para agregar una tienda
type AddStoreCommand struct {
	Nombre string
	URL    string
	Metodo int // 0 = shopify-pull, 1 = git-clone
	GitURL string
}

// Validate valida el comando
func (c AddStoreCommand) Validate() error {
	if c.Nombre == "" {
		return NewValidationError("store name is required")
	}
	if c.URL == "" {
		return NewValidationError("store URL is required")
	}
	if c.Metodo < 0 || c.Metodo > 1 {
		return NewValidationError("invalid download method")
	}
	if c.Metodo == 1 && c.GitURL == "" {
		return NewValidationError("git URL is required for git clone method")
	}
	return nil
}

// StartServerCommand representa el comando para iniciar un servidor
type StartServerCommand struct {
	StoreName string
	StoreURL  string
	StorePath string
}

// StopServerCommand representa el comando para detener un servidor
type StopServerCommand struct {
	StoreName string
}

// DeleteStoreCommand representa el comando para eliminar una tienda
type DeleteStoreCommand struct {
	StoreID string
}
