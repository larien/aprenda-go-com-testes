package poquer

import (
	"fmt"
	"io"
	"time"
)

// AlertadorDeBlind agenda alertas para quantias de blind
type AlertadorDeBlind interface {
	AgendarAlertaPara(duracao time.Duration, quantia int, to io.Writer)
}

// AlertadorDeBlindFunc te permite implementar o AlertadorDeBlind com uma função
type AlertadorDeBlindFunc func(duracao time.Duration, quantia int, to io.Writer)

// AgendarAlertaPara é uma implementação de AlertadorDeBlind para AlertadorDeBlindFunc
func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantia int, to io.Writer) {
	a(duracao, quantia, to)
}

// Alerter agenda alertas e os imprime para "to"
func Alerter(duracao time.Duration, quantia int, to io.Writer) {
	time.AfterFunc(duracao, func() {
		fmt.Fprintf(to, "Blind agora é %d\n", quantia)
	})
}
