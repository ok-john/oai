package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ImagePayload struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func ofPrompt(prompt string) ImagePayload {
	return ImagePayload{
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}
}

func (ai *ai_client) send(i *ImagePayload) {
	defer ai.output_file.Close()
	payloadBytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req := fmt_req_with("POST", "https://api.openai.com/v1/images/generations", body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	img_resp := ImageResponse{}
	json.Unmarshal(buff, &img_resp)
	if len(img_resp.Data) > 0 {
		req, err := http.NewRequest("GET", img_resp.Data[0].URL, nil)
		if err != nil {
			panic(err)
		}
		image_resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer image_resp.Body.Close()

		io.Copy(ai.output_file, image_resp.Body)

	}
}
