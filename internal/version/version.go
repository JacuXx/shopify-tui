package version

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Version    = "2.1.9"
	JsrScope   = "jacuxx"
	JsrPackage = "shopify-tui"
)

type JsrPackageInfo struct {
	Versions map[string]json.RawMessage `json:"versions"`
}

// VerificarActualizacion consulta el registro JSR y compara con la versión actual.
func VerificarActualizacion() (hayActualizacion bool, versionNueva string) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get("https://jsr.io/@" + JsrScope + "/" + JsrPackage + "/meta.json")
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, ""
	}

	var info JsrPackageInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return false, ""
	}

	latest := latestVersion(info.Versions)
	if latest != "" && latest != Version {
		return true, latest
	}

	return false, ""
}

// latestVersion extrae la versión semver más alta del mapa de versiones.
func latestVersion(versions map[string]json.RawMessage) string {
	if len(versions) == 0 {
		return ""
	}

	keys := make([]string, 0, len(versions))
	for k := range versions {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return compareSemver(keys[i], keys[j]) < 0
	})

	return keys[len(keys)-1]
}

// compareSemver compara dos versiones semver. Retorna -1, 0 o 1.
func compareSemver(a, b string) int {
	pa := parseSemver(a)
	pb := parseSemver(b)
	for i := 0; i < 3; i++ {
		if pa[i] < pb[i] {
			return -1
		}
		if pa[i] > pb[i] {
			return 1
		}
	}
	return 0
}

func parseSemver(v string) [3]int {
	var result [3]int
	parts := strings.SplitN(v, ".", 3)
	for i, p := range parts {
		if i >= 3 {
			break
		}
		n, _ := strconv.Atoi(p)
		result[i] = n
	}
	return result
}
