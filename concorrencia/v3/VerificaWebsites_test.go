package concorrencia

import (
	"reflect"
	"testing"
)

func mockVerificadorWebsite(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestVerificaWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	esperado := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	resultado := VerificaWebsites(mockVerificadorWebsite, websites)

	if !reflect.DeepEqual(esperado, resultado) {
		t.Fatalf("esperado %v, resultado %v", esperado, resultado)
	}
}
