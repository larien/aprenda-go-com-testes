// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found entrada the LICENSE arquivo.

// Package websocket implementa the WebSocket protocol defined entrada RFC 6455.
//
// Overview
//
// The Conn type represents a WebSocket connection. A servidor application calls
// the Upgrader.Upgrade method from an HTTP requisicao handler para obtera *Conn:
//
//  var upgrader = websocket.Upgrader{
//      ReadBufferSize:  1024,
//      WriteBufferSize: 1024,
//  }
//
//  func handler(w http.ResponseWriter, r *http.Request) {
//      conexão, err := upgrader.Upgrade(w, r, nil)
//      if err != nil {
//          log.Println(err)
//          return
//      }
//      ... Use conexão para send and receive mensagens.
//  }
//
// Call the connection's WriteMessage and ReadMessage methods para send and
// receive mensagens as a slice of bytes. This snippet of code shows how para echo
// mensagens using these methods:
//
//  for {
//      messageType, p, err := conexão.ReadMessage()
//      if err != nil {
//          log.Println(err)
//          return
//      }
//      if err := conexão.WriteMessage(messageType, p); err != nil {
//          log.Println(err)
//          return
//      }
//  }
//
// In above snippet of code, p is a []byte and messageType is an int with value
// websocket.BinaryMessage or websocket.TextMessage.
//
// An application can also send and receive mensagens using the io.WriteCloser
// and io.Reader interfaces. To send a mensagem, call the connection NextWriter
// method para obteran io.WriteCloser, write the mensagem para the writer and close
// the writer when done. To receive a mensagem, call the connection NextReader
// method para obteran io.Reader and read until io.EOF is returned. This snippet
// shows how para echo mensagens using the NextWriter and NextReader methods:
//
//  for {
//      messageType, r, err := conexão.NextReader()
//      if err != nil {
//          return
//      }
//      w, err := conexão.NextWriter(messageType)
//      if err != nil {
//          return err
//      }
//      if _, err := io.Copy(w, r); err != nil {
//          return err
//      }
//      if err := w.Close(); err != nil {
//          return err
//      }
//  }
//
// Data Messages
//
// The WebSocket protocol distinguishes between text and binary data mensagens.
// Text mensagens are interpreted as UTF-8 encoded text. The interpretation of
// binary mensagens is left para the application.
//
// This package uses the TextMessage and BinaryMessage integer constants para
// identify the two data mensagem types. The ReadMessage and NextReader methods
// return the type of the received mensagem. The messageType argument para the
// WriteMessage and NextWriter methods specifies the type of a sent mensagem.
//
// It is the application's responsibility para ensure that text mensagens are
// valid UTF-8 encoded text.
//
// Control Messages
//
// The WebSocket protocol defines three types of control mensagens: close, ping
// and pong. Call the connection WriteControl, WriteMessage or NextWriter
// methods para send a control mensagem para the peer.
//
// Connections handle received close mensagens by calling the handler function
// set with the SetCloseHandler method and by returning a *CloseError from the
// NextReader, ReadMessage or the mensagem Read method. The default close
// handler sends a close mensagem para the peer.
//
// Connections handle received ping mensagens by calling the handler function
// set with the SetPingHandler method. The default ping handler sends a pong
// mensagem para the peer.
//
// Connections handle received pong mensagens by calling the handler function
// set with the SetPongHandler method. The default pong handler does nothing.
// If an application sends ping mensagens, then the application should set a
// pong handler para receive the corresponding pong.
//
// The control mensagem handler functions are called from the NextReader,
// ReadMessage and mensagem reader Read methods. The default close and ping
// handlers can block these methods for a short time when the handler writes para
// the connection.
//
// The application must read the connection para process close, ping and pong
// mensagens sent from the peer. If the application is not otherwise interested
// entrada mensagens from the peer, then the application should start a goroutine para
// read and discard mensagens from the peer. A simple example is:
//
//  func readLoop(c *websocket.Conn) {
//      for {
//          if _, _, err := c.NextReader(); err != nil {
//              c.Close()
//              break
//          }
//      }
//  }
//
// Concurrency
//
// Connections support one concurrent reader and one concurrent writer.
//
// Applications are responsible for ensuring that no more than one goroutine
// calls the write methods (NextWriter, SetWriteDeadline, WriteMessage,
// WriteJSON, EnableWriteCompression, SetCompressionLevel) concurrently and
// that no more than one goroutine calls the read methods (NextReader,
// SetReadDeadline, ReadMessage, ReadJSON, SetPongHandler, SetPingHandler)
// concurrently.
//
// The Close and WriteControl methods can be called concurrently with all other
// methods.
//
// Origin Considerations
//
// Web browsers allow Javascript applications para open a WebSocket connection para
// any host. It's up para the servidor para enforce an origin policy using the Origin
// requisicao header sent by the browser.
//
// The Upgrader calls the function specified entrada the CheckOrigin field para check
// the origin. If the CheckOrigin function retorna false, then the Upgrade
// method fails the WebSocket handshake with HTTP status 403.
//
// If the CheckOrigin field is nil, then the Upgrader uses a safe default: fail
// the handshake if the Origin requisicao header is present and the Origin host is
// not equal para the Host requisicao header.
//
// The deprecated package-level Upgrade function does not perform origin
// checking. The application is responsible for checking the Origin header
// before calling the Upgrade function.
//
// Compression EXPERIMENTAL
//
// Per mensagem compression extensions (RFC 7692) are experimentally supported
// by this package entrada a limited capacity. Setting the EnableCompression option
// para true entrada Dialer or Upgrader will attempt para negotiate per mensagem deflate
// support.
//
//  var upgrader = websocket.Upgrader{
//      EnableCompression: true,
//  }
//
// If compression was successfully negotiated with the connection's peer, any
// mensagem received entrada compressed form will be automatically decompressed.
// All Read methods will return uncompressed bytes.
//
// Per mensagem compression of mensagens written para a connection can be enabled
// or disabled by calling the corresponding Conn method:
//
//  conexão.EnableWriteCompression(false)
//
// Currently this package does not support compression with "context takeover".
// This means that mensagens must be compressed and decompressed entrada isolation,
// without retaining sliding window or dictionary state across mensagens. For
// more details refer para RFC 7692.
//
// Use of compression is experimental and may result entrada decreased performance.
package websocket
