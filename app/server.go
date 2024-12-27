package main

import (
	"fmt"
	"net"
	"os"
  "strings"
  "strconv"
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
  	
   
   conn, err := l.Accept()
	  fmt.Println(conn)
   if err != nil {
	 	fmt.Println("Error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }
   
   buf:=make([]byte,512)
   _,errr:=conn.Read(buf)
   if errr!=nil{
     fmt.Println("error while reading", errr)
   }
   
 
   received :=strings.Split(string(buf)," ")
  
   if received[1]=="/"{
   conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
 }else if strings.Contains(received[1],"/echo"){
   str:=strings.TrimPrefix(received[1],"/echo/")

   conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type:text/plain\r\nContent-Length: "+strconv.Itoa(len(str))+"\r\n\r\n"+str))
 }else{
   conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

 }

}
