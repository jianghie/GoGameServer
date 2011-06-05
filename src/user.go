package GoGameServer

import (
  "net"
  "fmt"
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
  fmt.Printf( "User '%s' logged out.\n", u.Username )
  (*u.Conn).Close()
  u.Conn = nil
}


func (u *User) Write( msg string ) {
  (*u.Conn).Write( []byte(msg) )
}

