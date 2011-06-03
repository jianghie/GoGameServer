package GoGameServer

import (
  "mysql"
  "fmt"
  "net"
  "container/list"
  "os"
)


func HandleNewClient( conn *net.Conn, dbConn *mysql.Client, userList *list.List) {
  (*conn).Write( []byte( "Username: " ) )
  username, err := ReadLine( conn )
  if err != nil {
    (*conn).Close()
    return
  }

  (*conn).Write( []byte( "Password(hash): " ) )
  password, err := ReadLine( conn )
  if err != nil {
    (*conn).Close()
    return
  }


  authReply, err := Authenticate( username, password, dbConn )
  if err != nil {
    (*conn).Close()
  }

  if authReply.Authenticated {
    fmt.Printf( "User '%s' logged in.\n", username )
  } else {
    (*conn).Close()
    fmt.Print( "DROPPED..\n" )
  }
}



func main() {
  // Connect to the mysql server
  db, err := mysql.DialUnix( mysql.DEFAULT_SOCKET, "user", "pass123", "gogameserver" )
  if err != nil {
    fmt.Printf( "Error: %s\n", err.String() )
    os.Exit( 1 )
  }
  defer db.Close()

  userList := list.New()

  addr := net.TCPAddr{ net.ParseIP( "127.0.0.1" ), 9440 }
  netListen, err := net.ListenTCP( "tcp", &addr )
  if err != nil {
    os.Exit( 1 )
  }
  defer netListen.Close()

  for {
    fmt.Print( "Waiting for client..\n" );
    conn, err := netListen.Accept();
    if err != nil {
      fmt.Print( "Error encountered when accepting client." )
    }
    fmt.Print( "Accepted client.\n" )

    go HandleNewClient( &conn, db, userList )
  }
}

