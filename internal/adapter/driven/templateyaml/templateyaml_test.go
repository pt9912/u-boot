package templateyaml_test

import (
	"errors"
	"testing"
	"testing/fstest"

	"github.com/pt9912/u-boot/internal/adapter/driven/templateyaml"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

const validYAML = `apiVersion: github.com/pt9912/u-boot/template/v1
name: sample
description: "A sample template."
version: 1.0.0
`

func TestRead_Valid(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"tpl/template.yaml": &fstest.MapFile{Data: []byte(validYAML)},
	}
	meta, err := templateyaml.Read(fs, "tpl")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if meta.Name != "sample" {
		t.Errorf("Name = %q, want %q", meta.Name, "sample")
	}
	if meta.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", meta.Version, "1.0.0")
	}
}

func TestRead_Errors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		data      string
		present   bool
		wantInval bool // expect domain.ErrInvalidTemplate in the chain
	}{
		{
			name:      "unsupported apiVersion",
			data:      "apiVersion: github.com/pt9912/u-boot/template/v2\nname: x\ndescription: y\nversion: 1\n",
			present:   true,
			wantInval: true,
		},
		{
			name:      "missing metadata fields",
			data:      "apiVersion: github.com/pt9912/u-boot/template/v1\nname: x\n",
			present:   true,
			wantInval: true,
		},
		{
			name:    "unknown field rejected (KnownFields)",
			data:    validYAML + "bogusField: nope\n",
			present: true,
		},
		{
			name:    "missing template.yaml",
			present: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fs := fstest.MapFS{}
			if tc.present {
				fs["tpl/template.yaml"] = &fstest.MapFile{Data: []byte(tc.data)}
			}
			_, err := templateyaml.Read(fs, "tpl")
			if err == nil {
				t.Fatal("Read: want error, got nil")
			}
			if tc.wantInval && !errors.Is(err, domain.ErrInvalidTemplate) {
				t.Errorf("err = %v, want domain.ErrInvalidTemplate in chain", err)
			}
		})
	}
}
