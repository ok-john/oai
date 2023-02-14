package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type (
	// default_model =
	query struct {
		Model     string  `json:"model"`
		Prompt    string  `json:"prompt"`
		Temp      float64 `json:"temperature"`
		MaxTokens int     `json:"max_tokens"`
	}

	choice struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		LogProbs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	}

	usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	response struct {
		Id      string   `json:"id"`
		Object  string   `json:"object"`
		Created uint64   `json:"created"`
		Model   string   `json:"model"`
		Choices []choice `json:"choices"`
		Usage   usage    `json:"usage"`
	}
)

func UnmarshalResponse(f io.ReadCloser) (response, error) {
	ret := response{}
	decoder := json.NewDecoder(f)
	err := decoder.Decode(&ret)
	return ret, err
}

func UnmarshalQuery(f io.Reader) (query, error) {
	ret := query{}
	buf, err := io.ReadAll(f)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(buf, &ret)
	return ret, err
}

// func (q *query) MarshalQuery() (*bytes.Buffer, error) {
func (q *query) MarshalQuery() (io.Reader, error) {
	buf, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func main() {
	client := args.parse()
	if len(os.Args) > 1 && os.Args[1] == "help" {
		flag.Usage()
		os.Exit(0)
	}
	client.completion_query()
}

func (c *ai_client) completion_query() error {

	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input := string(buff)
	fmt.Printf("input=%s\n", input)
	prompt := &query{
		Model:     args.model,
		Prompt:    input,
		Temp:      args.temperature,
		MaxTokens: args.max_tokens,
	}
	body, err := prompt.MarshalQuery()
	if err != nil {
		panic(err)
	}
	req := fmt_req_with("POST", completions_url, body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}

	result, err := UnmarshalResponse(resp.Body)
	if err != nil {
		panic(err)
	}

	for _, choice := range result.Choices {
		if _, err := c.output_file.WriteString(choice.Text); err != nil {
			return err
		}
	}
	c.output_file.WriteString("\n")
	return nil
}
