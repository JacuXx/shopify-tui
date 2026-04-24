package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// ShopifyService implementa interacción con Shopify CLI
type ShopifyService struct{}

// NewShopifyService crea una nueva instancia
func NewShopifyService() *ShopifyService {
	return &ShopifyService{}
}

// AuthenticateStore autentica con una tienda Shopify
func (s *ShopifyService) AuthenticateStore(ctx context.Context, storeURL string) error {
	cmd := exec.CommandContext(ctx, "shopify", "auth", "login", "--store", storeURL)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DownloadShopifyTheme descarga un tema de Shopify
func (s *ShopifyService) DownloadShopifyTheme(ctx context.Context, storePath, storeURL string) error {
	cmd := exec.CommandContext(ctx,
		"shopify", "theme", "pull",
		"--store", storeURL,
		"--path", ".",
	)
	cmd.Dir = storePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull theme: %w", err)
	}

	return nil
}

// CloneGitRepository clona un repositorio Git
func (s *ShopifyService) CloneGitRepository(ctx context.Context, gitURL, storePath string) error {
	cmd := exec.CommandContext(ctx, "git", "clone", gitURL, ".")
	cmd.Dir = storePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// ThemeDev inicia el servidor de desarrollo de tema
func (s *ShopifyService) ThemeDev(ctx context.Context, storePath, storeURL string, port int) error {
	cmd := exec.CommandContext(ctx,
		"shopify", "theme", "dev",
		"--store", storeURL,
		"--port", fmt.Sprintf("%d", port),
	)
	cmd.Dir = storePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
