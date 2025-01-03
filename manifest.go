package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
)

type Manifest struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Version         string        `json:"version"`
	ManifestVersion int           `json:"manifestVersion"`
	Flags           []string      `json:"_flags,omitempty"`
	ImageDeltas     []*ImageDelta `json:"image_deltas,omitempty"`
	Files           *Files        `json:"files,omitempty"`
	Priority        json.Number   `json:"priority,omitempty"`
	Schema          string        `json:"schema,omitempty"`
}

type Files struct {
	Plugins []string `json:"plugins,omitempty"`
	Assets  []string `json:"assets,omitempty"`
	Files   []string `json:"files,omitempty"`
	Maps    []string `json:"maps,omitempty"`
	Text    []string `json:"text,omitempty"`
	Data    []string `json:"data,omitempty"`
}

type ImageDelta struct {
	Patch string `json:"patch"`
	With  string `json:"with"`
	Dir   bool   `json:"dir"`
}

func (p *Packer) openManifest() error {
	f, err := os.Open(path.Join(p.Source, "mod.json"))
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("No mod.json found in %s", p.Source)
		} else {
			color.Red("Error opening mod.json in %s: %s", p.Source, err)
		}
		return err
	}
	defer f.Close()

	manifest := new(Manifest)
	if err := json.NewDecoder(f).Decode(manifest); err != nil {
		color.Red("Error reading mod.json in %s: %s", p.Source, err)
		return fmt.Errorf("error reading mod.json in %s: %s", p.Source, err)
	}
	if manifest.ManifestVersion != 1 {
		color.Red("Unsupported manifest version in %s", p.Source)
		return fmt.Errorf("unsupported manifest version in %s", p.Source)
	}
	if manifest.Schema != "https://rph.space/oneloader.manifestv1.schema.json" {
		color.Yellow("Warning: Schema in %s is not the official schema", p.Source)
	}
	p.Manifest = manifest
	return nil
}
