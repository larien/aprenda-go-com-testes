package poquer

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websocketServidorJogador struct {
	*websocket.Conn
}

func (w *websocketServidorJogador) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func novoWebsocketServidorJogador(w http.ResponseWriter, r *http.Request) *websocketServidorJogador {
	conexão, err := atualizadorDeWebsocket.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("houve um problema ao atualizar a conexão para websockets %v\n", err)
	}

	return &websocketServidorJogador{conexão}
}

func (w *websocketServidorJogador) EsperarPelaMensagem() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("erro ao ler do websocket %v\n", err)
	}
	return string(msg)
}
