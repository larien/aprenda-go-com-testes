package main

import "reflect"

func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x) // ValorDe
	campo := valor.Field(0)     // Campo
	fn(campo.String())
}
