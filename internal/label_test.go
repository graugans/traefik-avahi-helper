package internal_test

import (
	"testing"

	"github.com/graugans/traefik-avahi-helper/internal"
)

// Make the test map more readable
type LabelMap map[string]string

func TestNewLabelParser(t *testing.T) {
	// Arange
	parser := internal.NewLabelParser()
	labels := make(LabelMap)
	labels["foo"] = "bar"
	// Act
	_, err := parser.FindLinkLocalHostName(labels)
	// Assert
	if err == nil {
		t.Error(
			"Expected an error when no hostname in the label map was given",
		)
	}
}

func TestIsTraefikEnabled(t *testing.T) {
	// Arrange
	parser := internal.NewLabelParser()
	// Defining the columns of the table
	var tests = []struct {
		name  string
		input LabelMap
		want  bool
	}{
		// the table itself
		{
			"An empty map shall return false",
			LabelMap{},
			false,
		},
		{
			"We expect traefik being enabled",
			LabelMap{
				"traefik.enable": "true",
			},
			true,
		},
		{
			"Malformed input shall return false",
			LabelMap{
				"traefik.enable": "xyz",
			},
			false,
		},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ans := parser.IsTraefikEnabled(tt.input)
			// Assert
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestFindLinkLocalHostName(t *testing.T) {
	// Arrange
	parser := internal.NewLabelParser()
	// Defining the columns of the table
	type wanted struct {
		host string
		err  error
	}
	var tests = []struct {
		name  string
		input LabelMap
		want  wanted
	}{
		// the table itself
		{
			"An empty map shall return no host",
			LabelMap{},
			wanted{"", internal.ErrLinkLocalHostNotFound},
		},
		{
			"A given hostname is expected",
			LabelMap{
				"traefik.http.routers.whoami.rule": "Host(`whoami.local`)",
			},
			wanted{"whoami.local", nil},
		},
		{
			"No .local extension",
			LabelMap{
				"traefik.http.routers.whoami.rule": "Host(`whoami.ege.io`)",
			},
			wanted{"", internal.ErrLinkLocalHostNotFound},
		},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ans, err := parser.FindLinkLocalHostName(tt.input)
			// Assert
			if err != tt.want.err {
				t.Errorf("error mismatch detected: got %v expected: %v", err, tt.want.err)
			}
			if ans != tt.want.host {
				t.Errorf("got %s, want %s", ans, tt.want.host)
			}
		})
	}
}
