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

// Testa a funcionalidade `searchIps` realizando uma consulta real para verificar IPs
func TestSearchIpsIntegration(t *testing.T) {
	app := Build()

	// Criar um novo contexto e associar o host ao flagset
	set := flag.NewFlagSet("test_search_ips", flag.ContinueOnError)
	set.String("host", "github.com", "host to look up IP addresses for")

	c := cli.NewContext(app, set, nil)

	// Redireciona o stdout para capturar a saída
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Executa a ação do comando `ip`
	err := cli.HandleAction(app.Commands[0].Action, c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Restaura o stdout e captura a saída
	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verifica se o output é um IP válido
	if !isValidIP(output) { // Função `isValidIP` deve estar disponível no pacote
		t.Errorf("Expected output to be a valid IP, got %s", output)
	}
}

// Testa a funcionalidade `searchServers` realizando uma consulta real para verificar servidores DNS
func TestSearchServersIntegration(t *testing.T) {
	app := Build()

	// Criar um novo contexto e associar o host ao flagset
	set := flag.NewFlagSet("test_search_servers", flag.ContinueOnError)
	set.String("host", "github.com", "host to look up NS records for")

	c := cli.NewContext(app, set, nil)

	// Redireciona o stdout para capturar a saída
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Executa a ação do comando `servidores`
	err := cli.HandleAction(app.Commands[1].Action, c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Restaura o stdout e captura a saída
	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verifica se o output contém "ns" (indicando que são registros de servidor)
	if !strings.Contains(output, "ns") {
		t.Errorf("Expected output to contain 'ns', got %s", output)
	}
}
