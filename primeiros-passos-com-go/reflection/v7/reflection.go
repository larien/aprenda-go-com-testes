package main

import (
	"reflect"
)

func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

	quantidadeDeValores := 0
	var obtemCampo func(int) reflect.Value

	switch valor.Kind() {
	case reflect.String:
		fn(valor.String())
	case reflect.Struct:
		quantidadeDeValores = valor.NumField()
		obtemCampo = valor.Field
	case reflect.Slice, reflect.Array:
		quantidadeDeValores = valor.Len()
		obtemCampo = valor.Index
	}

	for i := 0; i < quantidadeDeValores; i++ {
		percorre(obtemCampo(i).Interface(), fn)
	}
}

func obtemValor(x interface{}) reflect.Value {
	valor := reflect.ValueOf(x)

	if valor.Kind() == reflect.Ptr {
		valor = valor.Elem()
	}

	return valor
}
