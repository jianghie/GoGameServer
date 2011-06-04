package GoGameServer

import (
  "net"
)


type User struct {
  Id       int64
  Username string
  Conn     *net.Conn
}


func NewUser( id int64, username string, conn *net.Conn ) *User {
  u := &User{ id, username, conn }
  return u
}


func (u *User) Disconnect() {
  (*u.Conn).Close()
  u.Conn = nil
}

