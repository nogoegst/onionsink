// onionsink.go - onion data sink.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of onionsink, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import(
	"log"
	"flag"
	"os"
	"io"
	"net"
	"fmt"
	"crypto/rand"

	"github.com/nogoegst/bulb"
	"github.com/nogoegst/onionutil"
)


func handleSink(conn net.Conn, path string) {
	nameBin := make([]byte, 10)
	_, err := rand.Read(nameBin)
	if err != nil {
		log.Fatalf("Unable to generate filename: %v", err)
	}
	name := onionutil.Base32Encode(nameBin)[:16]
	f, err := os.Create(path+"/"+name)
	if err != nil {
		log.Printf("Unabale to create file: %v", err)
	}
	io.Copy(f, conn)
}

func main() {
	var debugFlag = flag.Bool("debug", false,
		"Show what's happening")
	var control = flag.String("control-addr", "default://",
		"Set Tor control address to be used")
	var controlPasswd = flag.String("control-passwd", "",
		"Set Tor control auth password")
	flag.Parse()
	debug := *debugFlag
	if len(flag.Args()) != 1 {
		log.Fatal("Please specify a path")
	}
	path := flag.Args()[0]
	// Connect to a running tor instance
	c, err := bulb.DialURL(*control)
	if err != nil {
		log.Fatalf("Failed to connect to control socket: %v", err)
	}
	defer c.Close()

	// See what's really going on under the hood
	c.Debug(debug)

	// Authenticate with the control port
	if err := c.Authenticate(*controlPasswd); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	l, err := c.AwaitListener(1, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l.Addr())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error while accepting connection: %v", err)
		}
		go handleSink(conn, path)
	}
}
