package context2

import (
	"testing"
	"time"
)

// SpyStore permite que você simule uma store e veja como ela é usada
type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

// Fetch retorna a resposta após um curto atraso
func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

// Cancel irá gravar a chamada
func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("store nao foi avisada para cancelar")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("store foi avisada para cancelar")
	}
}
