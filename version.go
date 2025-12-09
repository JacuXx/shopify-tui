package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	Version        = "1.5.3"
	NpmPackageName = "shopify-cli-tui"
)

type NpmPackageInfo struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

func verificarActualizacion() (hayActualizacion bool, versionNueva string) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get("https://registry.npmjs.org/" + NpmPackageName)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, ""
	}

	var info NpmPackageInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return false, ""
	}

	if info.DistTags.Latest != Version {
		return true, info.DistTags.Latest
	}

	return false, ""
}

func compararVersiones(actual, nueva string) bool {
	return actual != nueva
}
