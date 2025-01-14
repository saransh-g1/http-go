package main

import (
	"fmt"
	"net/http"
	"os"
  "net"
  "strings"
  "bufio"
  "io"
  "compress/gzip"
  "bytes"
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
   

    if req.Header.Get("Accept-Encoding")!=""{
      headers:=strings.Split(req.Header.Get("Accept-Encoding"),",")
      var flag bool=false
      for _,val:=range headers{
        if(strings.TrimSpace(val)=="gzip"){
          flag=true
        }
      }

      if flag{
        
        str:=strings.TrimPrefix(req.URL.Path,"/echo/")
        var b bytes.Buffer
        w:=gzip.NewWriter(&b) 
        w.Write([]byte(str))
        w.Close()
	      
        fmt.Println(b.String())
        fmt.Fprintf(conn,"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %s\r\n\r\n%s",fmt.Sprint(len(b.String())),b.String())

     }else{
      fmt.Fprintf(conn,"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")
     }
      //has to be seen
     return
    }

    if req.Method=="POST"{
      dirfn:=os.Args[2]+strings.TrimPrefix(req.URL.Path,"/files")
      file,err:=os.Create(dirfn)
      if err!=nil{
        fmt.Println(err)
      }
      data,err:=io.ReadAll(req.Body)
      if err!=nil{
        fmt.Println(err)
      }
      
      _,err=file.WriteString(string(data))
      if err!=nil{
        fmt.Println(err)
      }
      fmt.Fprintf(conn,"HTTP/1.1 201 Created\r\n\r\n")

    } else if req.URL.Path=="/"{
      fmt.Fprintf(conn,"HTTP/1.1 200 OK\r\n\r\n")
   }else if strings.Contains(req.URL.Path,"/echo"){
     str:=strings.TrimPrefix(req.URL.Path,"/echo/")
   fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
   }else if req.URL.Path=="/user-agent"{
	 fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(req.UserAgent()), req.UserAgent())
   }else if strings.Contains(req.URL.Path,"/files") {
     str:=strings.TrimPrefix(req.URL.Path,"/files")
     args:=os.Args[2]
     dirfn:=args+str
     file,err:=os.ReadFile(dirfn)
     if err!=nil{
       fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
     }
     fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), string(file))
   }else{
		fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
   }
 
}
