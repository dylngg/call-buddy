package telephono_test

import (
	"bufio"
	"github.com/call-buddy/call-buddy/telephono/cmd/test_server"
	"net/http/httptest"
	"strings"
)

func prefixAllLinesOfString(in string, with rune) string {
	outBuilder := strings.Builder{}
	inScanner := bufio.NewScanner(strings.NewReader(in))

	for inScanner.Scan() {
		line := inScanner.Text()
		outBuilder.WriteRune(with)
		outBuilder.WriteString(line)
		outBuilder.WriteString("\n")
	}

	return outBuilder.String()
}

type globalTestStateType struct {
	GlobalHitServer          *httptest.Server
	GlobalHitServerAddr      string
	GlobalHitServerIsRunning bool
}

func (g *globalTestStateType) getPrefix() string {
	return g.GlobalHitServer.URL
}

var GlobalTestState = globalTestStateType{
	GlobalHitServerAddr:      "127.0.0.1:8096",
	GlobalHitServerIsRunning: false,
}

//Idempotent operation that will set up the server
func setUpServer() {
	if !GlobalTestState.GlobalHitServerIsRunning {

		if GlobalTestState.GlobalHitServer == nil {
			GlobalTestState.GlobalHitServer = httptest.NewServer(test_server.TestServerMux())
			GlobalTestState.GlobalHitServerAddr = GlobalTestState.GlobalHitServer.URL
		}
		GlobalTestState.GlobalHitServerIsRunning = true
	}
}
