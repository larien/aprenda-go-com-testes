package poquer

import (
	"fmt"
	"io"
	"time"
)

// AlertadorDeBlind agenda alertas para quantias de blind
type AlertadorDeBlind interface {
	AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer)
}

// AlertadorDeBlindFunc te permite implementar o AlertadorDeBlind com uma função
type AlertadorDeBlindFunc func(duracao time.Duration, quantia int, para io.Writer)

// AgendarAlertaPara é uma implementação de AlertadorDeBlind para AlertadorDeBlindFunc
func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer) {
	a(duracao, quantia, para)
}

// Alertador agenda alertas e os imprime para "para"
func Alertador(duracao time.Duration, quantia int, para io.Writer) {
	time.AfterFunc(duracao, func() {
		fmt.Fprintf(para, "Blind agora é %d\n", quantia)
	})
}
