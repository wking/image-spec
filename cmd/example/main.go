package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/image-spec/image/cas"
	casLayout "github.com/opencontainers/image-spec/image/cas/layout"
	utilLayout "github.com/opencontainers/image-spec/image/layout"
	refsLayout "github.com/opencontainers/image-spec/image/refs/layout"
	imageSpecs "github.com/opencontainers/image-spec/specs-go"
	"github.com/opencontainers/image-spec/specs-go/v1"
	runtimeSpecs "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	bundlePath := "bundle"
	tarPath := "image.tar"
	refName := "v1"

	configPath := filepath.Join(bundlePath, "config.json")
	runtimeConfig, err := readConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	imageConfig, err := runtimeConfigToImageConfig(runtimeConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = utilLayout.CreateTarFile(ctx, tarPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	casEngine, err := casLayout.NewEngine(ctx, tarPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer casEngine.Close()

	configDescriptor, err := cas.PutJSON(ctx, casEngine, imageConfig, "application/vnd.oci.image.serialization.config.v1+json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	layers := []imageSpecs.Descriptor{} // FIXME

	manifest := &v1.Manifest{
		Versioned: imageSpecs.Versioned{
			SchemaVersion: 2,
			MediaType:     "application/vnd.oci.image.manifest.v1+json",
		},
		Config: *configDescriptor,
		Layers: layers,
	}
	manifestDescriptor, err := cas.PutJSON(ctx, casEngine, manifest, manifest.MediaType)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	refEngine, err := refsLayout.NewEngine(ctx, tarPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer refEngine.Close()

	err = refEngine.Put(ctx, refName, manifestDescriptor)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readConfig(path string) (config *runtimeSpecs.Spec, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg = runtimeSpecs.Spec{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

func runtimeConfigToImageConfig(runtimeConfig *runtimeSpecs.Spec) (imageConfig *v1.ImageConfig, err error) {
	// FIXME
	return &v1.ImageConfig{}, nil
}
