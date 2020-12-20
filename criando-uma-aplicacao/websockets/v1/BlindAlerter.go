package poquer

import (
	"fmt"
	"os"
	"time"
)

// AlertadorDeBlind agenda alertas para quantias de blind
type AlertadorDeBlind interface {
	AgendarAlertaPara(duracao time.Duration, quantia int)
}

// AlertadorDeBlindFunc te permite implementar o AlertadorDeBlind com uma função
type AlertadorDeBlindFunc func(duracao time.Duration, quantia int)

// AgendarAlertaPara é uma implementação de AlertadorDeBlind para AlertadorDeBlindFunc
func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantia int) {
	a(duracao, quantia)
}

// SaidaAlertador agenda alertas e os imprime para os.Stdout
func SaidaAlertador(duracao time.Duration, quantia int) {
	time.AfterFunc(duracao, func() {
		fmt.Fprintf(os.Stdout, "Blind agora é %d\n", quantia)
	})
}
