package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	logger = SetupLogger("/etc/openai/logs")
)

func main() {
	client := args.parse()
	if len(os.Args) > 1 && os.Args[1] == "help" {
		flag.Usage()
		os.Exit(0)
	}

	if client.args.list_models {
		if _, err := client.available_models(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if client.args.interactive {
		fmt.Println("Running in interactive mode...")
		interactive(&client)
		os.Exit(0)
	}

	client.completion_query()
}

func interactive(client *ai_client) {
	state := &query{
		Model:     client.args.model,
		Prompt:    []chat{},
		Temp:      client.args.temperature,
		MaxTokens: client.args.max_tokens,
	}
	sc := bufio.NewScanner(os.Stdin)
	line, prev_line := []byte{}, []byte{}
	for sc.Scan() {
		prev_line, line = line, sc.Bytes()
		if client.args.debug {
			fmt.Printf("previous_line: %s\n", prev_line)
			fmt.Printf("current_line: %s\n", line)
		}
		if len(line) == 0 {
			continue
		}
		// state.Prompt = append(state.Prompt, NewChat("admin", string(line)))

		logger.Printf("user: %s\n", line)
		state.add_chat(line)
		state.send_chat(client)
		if client.args.debug {
			r, _ := state.MarshalQuery()
			buff, _ := io.ReadAll(r)
			result := bytes.NewBuffer([]byte{})
			json.Indent(result, buff, "", "\t")
			io.Copy(os.Stdout, result)
		}

	}
}
