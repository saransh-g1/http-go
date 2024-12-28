package main

import (
	"fmt"
	"net/http"
	"os"
  "net"
  "strings"
  "bufio"
)


// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }
  	
  for {
   conn, err := l.Accept()
   if err != nil {
	 	fmt.Println("Error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }
  
  go resolveHeaders(conn) 
}
  

}


func resolveHeaders (conn net.Conn){
      reader:=bufio.NewReader(conn)
   req,errr:=http.ReadRequest(reader)
   if errr!=nil{
     fmt.Println("error while reading", errr)
   }

   if req.URL.Path=="/"{
      fmt.Fprintf(conn,"HTTP/1.1 200 OK\r\n\r\n")
   }else if strings.Contains(req.URL.Path,"/echo"){
     str:=strings.TrimPrefix(req.URL.Path,"/echo/")
   fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
   }else if req.URL.Path=="/user-agent"{
	 fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(req.UserAgent()), req.UserAgent())
   }else if strings.Contains(req.URL.Path,"/files") {
     str:=strings.TrimPrefix(req.URL.Path,"/files")
     args:=os.Args[0]
     dirfn:=args+str
     file,err:=os.ReadFile(dirfn)
     if err!=nil{
       fmt.Println(err)
     }
     fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), string(file))
   }else{
		fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
   }
 
}
