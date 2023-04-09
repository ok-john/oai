package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
)

const (
	LOG_DIR = "/etc/openai/logs"
)

var (
	logger = SetupLogger(LOG_DIR)
)

func sort_numerically(files []os.FileInfo) []os.FileInfo {
	sort.Slice(files, func(i, j int) bool {
		num1, err1 := strconv.Atoi(files[i].Name())
		num2, err2 := strconv.Atoi(files[j].Name())

		if err1 != nil || err2 != nil {
			return files[i].Name() < files[j].Name()
		}

		return num1 < num2
	})

	return files
}

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

	if client.args.latest_logs {
		files, err := ioutil.ReadDir(LOG_DIR)
		if err != nil {
			panic(err)
			os.Exit(1)
		}

		files = sort_numerically(files)
		if len(files) < 1 {
			os.Stdout.WriteString(fmt.Sprintf("No log files found at %s\n", LOG_DIR))
			os.Exit(0)
		}
		last := files[len(files)-2].Name()
		os.Stdout.WriteString(fmt.Sprintf("reading log-file %s\n", last))
		logf, err := os.Open(path.Join(LOG_DIR, last))
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, logf)
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
