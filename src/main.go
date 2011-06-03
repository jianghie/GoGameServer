package GoGameServer

import (
  "mysql"
  "fmt"
  "net"
  "container/list"
  "os"
  "bufio"
)


func readline( b *bufio.Reader ) (p []byte, err os.Error ) {
  if p, err = b.ReadSlice('\n'); err != nil {
    return nil, err
  }
  var i int
  for i = len( p ); i > 0; i-- {
    if c:= p[i-1]; c != ' ' && c != '\r' && c != '\t' && c != '\n' {
      break
    }
    if string(p[i-2:i]) == "\r\n" {
      return p[0:i-2], nil
    } else if string(p[i-1:i]) == "\n" {
      return p[0:i-1], nil
    } else return p[0:i], nil


  }
  return nil, nil
}



func ReadLine( conn *net.Conn ) (string, os.Error) {
  reader := bufio.NewReader( (*conn) )
  line, err := readline( reader )
  if err != nil {
    fmt.Print( "Error encountered when reading a line!\n" )
    return "ERROR", os.NewError( "Couldn't read line" )
  }
  return string(line), nil
}



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


  fmt.Printf( "Username: '%s'\n", username )
  fmt.Printf( "Password: '%s'\n", password )

  authReply, err := Authenticate( username, password, dbConn )
  if err != nil {
    (*conn).Close()
  }

  if authReply.Authenticated {
    fmt.Print( "AUTHENTICATED!\n" )
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

