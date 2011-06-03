package GoGameServer

import (
  "mysql"
  "fmt"
  "os"
)


type AuthenticationReply struct {
  Authenticated bool
  Id            int64
  Username      string
}


func NewAuthenticationReply( authenticated bool,
                             id int64,
                             username string ) AuthenticationReply {
  var r AuthenticationReply
  r.Authenticated = authenticated
  r.Id            = id
  r.Username      = username
  return r
}



func Authenticate( username string,
                   passwordHash string,
                   dbConn *mysql.Client ) (AuthenticationReply, os.Error) {
  // Escape input
  username = dbConn.Escape( username )
  password := dbConn.Escape( passwordHash )

  err := dbConn.Query( "SELECT * FROM users WHERE nick = '"+username+"' AND password = '"+password+"' limit 1" )
  if err != nil {
    return NewAuthenticationReply( false, -1, "" ), err
  }

  result, err := dbConn.UseResult()
  if err != nil {
    return NewAuthenticationReply( false, -1, "" ), err
  }
  /*defer result.Free()*/

  // Fetch the row
  row := result.FetchMap()

  // If we found it the client got the username and password right
  if row != nil {
    id       := row["id"].(int64)
    nick     := row["nick"].(string)

    return NewAuthenticationReply( true, id, nick ), nil
  } else {
    fmt.Printf( "No rows found || Bad username and/or password.\n" )
  }

  return NewAuthenticationReply( false, -1, "" ), os.NewError( "Authentication failed." )
}

