package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

var (
	args cmd_args = cmd_args{}
)

type (
	cmd_args struct {
		environment     string
		tor             bool
		socks5_hostname string
		org_id          string
		output_file     string
		api_key         string
		model           string
		list_models     bool
		max_tokens      int
		temperature     float64
	}
)

func (args *cmd_args) parse() ai_client {
	flag.StringVar(&args.environment, "env", default_environment, "absolute path to the environment directory.")
	flag.StringVar(&args.org_id, "org", "", "optionally specify an organization id.")
	flag.BoolVar(&args.tor, "tor", false, "toggles the use of a socks5 tor proxy.")
	flag.BoolVar(&args.list_models, "l", false, "list available openai models.")
	flag.StringVar(&args.socks5_hostname, "socks5-hostname", "localhost:9050", "optionally override the default tor proxy address.")
	flag.StringVar(&args.output_file, "o", "", "optionally specify an output file, defaults to stdout.")
	flag.StringVar(&args.model, "model", "text-davinci-003", "model to use")
	flag.Float64Var(&args.temperature, "temp", 0.5, "temperature of model")
	flag.IntVar(&args.max_tokens, "max", 2000, "max tokens")
	flag.Parse()
	if stat, err := os.Stat(args.environment); errors.Is(err, os.ErrNotExist) {
		fmt.Println("environment directory is non-existent.")
		if err = os.MkdirAll(args.environment, os.ModePerm); err != nil {
			fmt.Println("please rerun with sudo or\n\tmkdir -p " + args.environment + "\n")
			fmt.Println("place your api key in a file at: \n\t" + path.Join(args.environment, api_key_filename) + "\n")
			os.Exit(1)
		}
		if !stat.IsDir() {
			fmt.Printf("environment %s is not a directory", args.environment)
			os.Exit(1)
		}
	}
	var out_file *os.File = os.Stdout
	var ps_errno error = nil
	if len(args.output_file) > 0 {
		out_file, ps_errno = os.Open(args.output_file)
		if ps_errno != nil {
			fmt.Printf("bad output file: %s", args.output_file)
			panic(ps_errno)
		}
	}
	_api_key, err := read_env_file(args, api_key_filename)
	if err != nil {
		fmt.Printf("bad: no api key at %s", path.Join(args.environment, api_key_filename))
		panic(err)
	}
	c := default_client(args, args.tor, args.api_key, out_file)
	// go c.check_ip()
	c.api_key = string(_api_key)
	args.api_key = strings.ReplaceAll(string(_api_key), "\n", "")
	return c
}

func read_env_file(args *cmd_args, fn string) ([]byte, error) {
	fn = path.Join(args.environment, fn)
	f, err := os.Open(fn)
	if err != nil {
		return []byte{}, err
	}
	return io.ReadAll(f)
}
