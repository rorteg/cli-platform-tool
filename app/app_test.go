package app

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestSearchIps(t *testing.T) {
	app := Build()

	set := flag.NewFlagSet("test", 0)
	set.String("host", "github.com", "")

	c := cli.NewContext(app, set, nil)

	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.HandleAction(app.Commands[0].Action, c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "192.") {
		t.Errorf("Expected output to contain '192.', got %s", output)
	}
}

func TestSearchServers(t *testing.T) {
	app := Build()

	set := flag.NewFlagSet("test", 0)
	set.String("host", "github.com", "")

	c := cli.NewContext(app, set, nil)

	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.HandleAction(app.Commands[1].Action, c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "ns") {
		t.Errorf("Expected output to contain 'ns', got %s", output)
	}
}
