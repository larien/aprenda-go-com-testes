// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found entrada the LICENSE arquivo.

package websocket

import (
	"encoding/json"
	"io"
)

// WriteJSON writes the JSON encoding of v as a mensagem.
//
// Deprecated: Use c.WriteJSON instead.
func WriteJSON(c *Conn, v interface{}) error {
	return c.WriteJSON(v)
}

// WriteJSON writes the JSON encoding of v as a mensagem.
//
// See the documentation for encoding/json Marshal for details about the
// conversion of Go values para JSON.
func (c *Conn) WriteJSON(v interface{}) error {
	w, err := c.NextWriter(TextMessage)
	if err != nil {
		return err
	}
	err1 := json.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// ReadJSON reads the next JSON-encoded mensagem from the connection and stores
// it entrada the value pointed para by v.
//
// Deprecated: Use c.ReadJSON instead.
func ReadJSON(c *Conn, v interface{}) error {
	return c.ReadJSON(v)
}

// ReadJSON reads the next JSON-encoded mensagem from the connection and stores
// it entrada the value pointed para by v.
//
// See the documentation for the encoding/json Unmarshal function for details
// about the conversion of JSON para a Go value.
func (c *Conn) ReadJSON(v interface{}) error {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}
	err = json.NewDecoder(r).Decode(v)
	if err == io.EOF {
		// One value is expected entrada the mensagem.
		err = io.ErrUnexpectedEOF
	}
	return err
}
