package main

import (
	"bufio"
	"strings"
	"flag"
	"fmt"
	d "github.com/tj/go-debug"
	"io"
	"log"
	"net"
	"os"
)

var debug = d.Debug("go_redis:cli")

const (
	dELIM = byte('\n')
	quitCMD = "exit"
)

func main() {
	//parse arguments
	host := flag.String("h", "localhost", "ip of server")
	port := flag.Int("p", 3333, "port of server listening on")
	flag.Parse()
	//connect
	ip := net.ParseIP(*host)
	raddr := &net.TCPAddr{
		IP:   ip,
		Port: *port,
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalf("dial for ip:%s, port:%d fail: %s", *host, port, err)
	}
	defer conn.Close()
	br := bufio.NewReader(conn)
	line, err := br.ReadString(dELIM)
	if err != nil {
		log.Fatalf("read string from conn:%s", err)
	}
	fmt.Printf("connct success, get reply:\n%s", line)
	loopProcessCommand(conn)
}

func loopProcessCommand(c *net.TCPConn) {
	stdin := os.Stdin
	//read line
	scanner := bufio.NewScanner(stdin)
	fmt.Printf(">")
	for scanner.Scan() {
		command := scanner.Text()
		if strings.TrimSpace(command) == quitCMD {
			os.Exit(0)
		}
		//encode
		en := []byte(encode(command))
		debug("send command:%s", string(en))
		//write to conn
		_, err := c.Write(en)
		if err != nil {
			log.Fatalf("write to conn:%s\n", err)
		}
		debug("write command done")
		//read reply
		line, err := readLine(c)
		if err != nil {
			log.Fatalf("read string from conn:%s", err)
		}
		//display
		fmt.Printf("%s\n", line)
		fmt.Printf(">")
	}
}

func readLine(r io.Reader) (string, error) {
	line := ""
	var buf [1]byte
	debug("begin read")
	for {
		debug("read char")
		c, err := r.Read(buf[:])
		debug("read count:%d", c)
		if err != nil {
			return "", err
		}
		debug("read stdin char:%s", string(buf[:]))
		if buf[0] == dELIM {
			return line, nil
		}
		line = line + string(buf[:])
	}
}

func encode(content string) string {
	return content + "\n"
}
