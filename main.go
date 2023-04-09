package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
)

type (
	// default_model =
	query struct {
		Model     string  `json:"model"`
		Prompt    []chat  `json:"messages"`
		Temp      float64 `json:"temperature"`
		MaxTokens int     `json:"max_tokens"`
	}

	completion_choice struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		LogProbs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	}

	choice struct {
		Message chat   `json:"message"`
		Finish  string `json:"finish_reason"`
		Index   uint64 `json:"index"`
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

	chat struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

func NewChat(role, content string) chat {
	return chat{
		Role:    role,
		Content: content,
	}
}

func (c *chat) MarshalChat() (io.Reader, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

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

	if client.args.list_models {
		if _, err := client.available_models(); err != nil {
			panic(err)
		}
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

	chatts := NewChat("system", input)

	prompt := &query{
		Model:     args.model,
		Prompt:    []chat{chatts},
		Temp:      args.temperature,
		MaxTokens: args.max_tokens,
	}
	body, err := prompt.MarshalQuery()
	if err != nil {
		panic(err)
	}

	req := fmt_req_with("POST", chat_url, body)
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
		if _, err := c.output_file.WriteString(choice.Message.Content); err != nil {
			return err
		}
	}
	c.output_file.WriteString("\n")
	return nil
}
