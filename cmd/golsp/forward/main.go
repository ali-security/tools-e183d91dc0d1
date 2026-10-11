// The forward command writes and reads to a golsp server on a network socket.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/tools/internal/lsp/cmd"
	"golang.org/x/tools/internal/tool"
)

func main() {
	tool.Main(context.Background(), &app{Server: &cmd.Server{}}, os.Args[1:])
}

type app struct {
	*cmd.Server
	Port int `flag:"port" help:"port on which golsp is listening"`
}

func (*app) Name() string               { return "forward" }
func (*app) Usage() string              { return "[-port=<value>]" }
func (*app) ShortHelp() string          { return "An intermediary between an editor and GoLSP." }
func (*app) DetailedHelp(*flag.FlagSet) {}

func (a *app) Run(ctx context.Context, args ...string) error {
	if a.Port == 0 {
		a.ShortHelp()
		os.Exit(0)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", a.Port))
	if err != nil {
		log.Print(err)
		os.Exit(0)
	}

	go func(conn net.Conn) {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			log.Print(err)
			os.Exit(0)
		}
	}(conn)

	go func(conn net.Conn) {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Print(err)
			os.Exit(0)
		}
	}(conn)

	for {
	}
}
