package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Packer struct {
	Source   string
	Manifest *Manifest
	Zip      *zip.Writer
	Buf      *bytes.Buffer
}

func (p *Packer) Dest() {
	zipPath := path.Join(p.Source, "..", p.Manifest.ID+"-"+p.Manifest.Version+".zip")
	color.Blue("Creating %s", path.Base(zipPath))
	f, err := os.Create(zipPath)
	if err != nil {
		color.Red("Error creating zip file %s", err)
		color.Red(err.Error())
		return
	}
	defer f.Close()
	if _, err := io.Copy(f, p.Buf); err != nil {
		color.Red("Error writing zip file %s", err)
		color.Red(err.Error())
		return
	}
	color.Green("Success!")
}

func (p *Packer) Pack() {
	color.Blue("Packing %s", p.Manifest.Name)
	color.Blue("ID: %s, Version: %s", p.Manifest.ID, p.Manifest.Version)
	p.Buf = new(bytes.Buffer)
	p.Zip = zip.NewWriter(p.Buf)
	defer p.Zip.Close()

	if err := p.WriteManifestToZip(); err != nil {
		return
	}
	if p.Manifest.Files != nil {
		for _, name := range p.Manifest.Files.Plugins {
			p.WriteFileToZip(name)
		}
		for _, name := range p.Manifest.Files.Assets {
			p.WriteFileToZip(name)
		}
		for _, name := range p.Manifest.Files.Files {
			p.WriteFileToZip(name)
		}
		for _, name := range p.Manifest.Files.Maps {
			p.WriteFileToZip(name)
		}
		for _, name := range p.Manifest.Files.Text {
			p.WriteFileToZip(name)
		}
		for _, name := range p.Manifest.Files.Data {
			p.WriteFileToZip(name)
		}
	}
	if p.Manifest.ImageDeltas != nil {
		for _, delta := range p.Manifest.ImageDeltas {
			p.WriteFileToZip(delta.With)
		}
	}

	for _, pattern := range strings.Split(*flagInclude, ",") {
		files, err := filepath.Glob(path.Join(p.Source, pattern))
		if err != nil {
			color.Red("Error globbing %s", pattern)
			color.Red(err.Error())
			continue
		}
		for _, file := range files {
			name, _ := filepath.Rel(p.Source, file)
			p.WriteFileToZip(name)
		}
	}
}

func (p *Packer) WriteManifestToZip() error {
	manifestBuf := new(bytes.Buffer)
	e := json.NewEncoder(manifestBuf)
	e.SetIndent("", "  ")
	if err := e.Encode(p.Manifest); err != nil {
		color.Red("Error encoding manifest %s", err)
		color.Red(err.Error())
		return err
	}
	return p.WriteDataToZip("mod.json", manifestBuf)
}

func (p *Packer) WriteFileToZip(name string) error {
	f, err := os.Open(path.Join(p.Source, name))
	if err != nil {
		color.Red("Error opening %s", name)
		color.Red(err.Error())
		return err
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		color.Red("Error reading %s", name)
		color.Red(err.Error())
		return err
	}
	if strings.HasSuffix(name, ".json") ||
		strings.HasSuffix(name, ".jsond") ||
		strings.HasSuffix(name, ".yamld") ||
		strings.HasSuffix(name, ".ymld") {
		nb := new(bytes.Buffer)
		if err := json.Compact(nb, buf.Bytes()); err != nil {
			color.Red("Error compacting %s", name)
			color.Red(err.Error())
			return err
		}
		buf = nb
	}
	return p.WriteDataToZip(name, buf)
}

func (p *Packer) WriteDataToZip(name string, data io.Reader) error {
	name = strings.Replace(name, "\\", "/", -1)
	name = p.Manifest.ID + "/" + name
	w, err := p.Zip.Create(name)
	if err != nil {
		color.Red("Error writing %s", name)
		color.Red(err.Error())
		return err
	}
	if _, err := io.Copy(w, data); err != nil {
		color.Red("Error writing %s", name)
		color.Red(err.Error())
		return err
	}
	color.Blue("Added %s", name)
	return nil
}
