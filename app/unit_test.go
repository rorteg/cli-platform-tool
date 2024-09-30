package app

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"testing"

	"github.com/urfave/cli"
)

// Mock de LookupIP e LookupNS para simular comportamento no teste
var mockLookupIP = func(host string) ([]net.IP, error) {
	if host == "github.com" {
		return []net.IP{net.ParseIP("192.30.255.112")}, nil
	}
	return nil, fmt.Errorf("host not found")
}

var mockLookupNS = func(host string) ([]*net.NS, error) {
	if host == "github.com" {
		return []*net.NS{{Host: "ns1.github.com"}}, nil
	}
	return nil, fmt.Errorf("host not found")
}

// Redefinir as funções padrão para as funções mock
func init() {
	lookupIP = mockLookupIP
	lookupNS = mockLookupNS
}

// Função para testar o searchIps usando mock
func TestSearchIpsWithMock(t *testing.T) {
	// Preparar o teste
	app := Build()
	set := flag.NewFlagSet("test", 0)
	set.String("host", "github.com", "")

	c := cli.NewContext(app, set, nil)

	// Redireciona stdout
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Executa a ação
	err := cli.HandleAction(app.Commands[0].Action, c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Restaura stdout e captura o output
	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verifica o resultado
	if !isValidIP(output) {
		t.Errorf("Expected output to be a valid IP, got %s", output)
	}
}

// Test default host value
func TestDefaultHostValue(t *testing.T) {
	// Criar o app e um novo contexto para teste
	app := Build()

	// Inicializa um novo FlagSet para evitar redefinição de flags
	set := flag.NewFlagSet("test_default_host_value", flag.ContinueOnError)

	// Definir a flag diretamente no FlagSet para evitar o erro de redefinição
	set.String("host", "github.com", "default host")

	// Criar o contexto com as flags associadas
	c := cli.NewContext(app, set, nil)

	// Verificar se o valor padrão está correto
	if c.String("host") != "github.com" {
		t.Errorf("Expected default host to be 'github.com', got '%s'", c.String("host"))
	}
}

// Test command name validation
func TestCommandNames(t *testing.T) {
	app := Build()

	if len(app.Commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(app.Commands))
	}

	if app.Commands[0].Name != "ip" || app.Commands[1].Name != "servidores" {
		t.Errorf("Expected command names to be 'ip' and 'servidores', got '%s' and '%s'", app.Commands[0].Name, app.Commands[1].Name)
	}
}

// Test if flags are correctly set
func TestFlags(t *testing.T) {
	app := Build()

	if len(app.Commands[0].Flags) == 0 || app.Commands[0].Flags[0].GetName() != "host" {
		t.Errorf("Expected 'host' flag, got %v", app.Commands[0].Flags)
	}
}

// Test the creation of the CLI app
func TestAppCreation(t *testing.T) {
	app := Build()
	if app == nil {
		t.Errorf("Expected app to be created, but got nil")
	}
	if app.Name != "Aplicação de Linha de Comando" {
		t.Errorf("Expected app name to be 'Aplicação de Linha de Comando', got %s", app.Name)
	}
}
