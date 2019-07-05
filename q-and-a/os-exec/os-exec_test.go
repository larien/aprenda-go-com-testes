package osexec

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "FELIZ ANO NOVO!"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestGetData(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Gatos são os melhores animais</message>
</payload>`)

	got := GetData(input)
	want := "GATOS SÃO OS MELHORES ANIMAIS"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
