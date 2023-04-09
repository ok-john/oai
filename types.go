package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type (
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

func (p *query) add_chat(chat []byte) {
	p.Prompt = append(p.Prompt, NewChat("system", string(chat)))
}

func (p *query) send_chat(c *ai_client) error {

	body, err := p.MarshalQuery()
	if err != nil {
		panic(err)
	}

	req := fmt_req_with("POST", chat_url, body)
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}

	result, err := UnmarshalResponse(resp.Body)
	if err != nil {
		panic(err)
	}

	if c.args.debug {
		for choice := range result.Choices {
			fmt.Printf("choice: %+v\n", choice)
		}
	}

	for _, choice := range result.Choices {
		if _, err := c.output_file.WriteString(choice.Message.Content); err != nil {
			return err
		}
		logger.Printf("openai: %s\n", choice.Message.Content)
	}
	c.output_file.WriteString("\n")
	return nil
}

func SetupLogger(log_dir string) *log.Logger {
	ts := strconv.Itoa(int(time.Now().Unix()))
	log_fn := path.Join(log_dir, ts) + ".log"
	log_fil, err := os.Create(log_fn)
	if err != nil {
		panic(err)
	}
	os.Stdout.WriteString("\nLOGGER CREATAED AT: " + log_fn + "\n")
	return log.New(log_fil, "", 0)
}
