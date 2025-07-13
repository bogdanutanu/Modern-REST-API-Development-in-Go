package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	token   string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(user, password string) error {
	data := LoginRequest{
		User:     user,
		Password: password,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.baseURL+"/login", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid login")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	loginResponse := LoginResponse{}
	json.Unmarshal(responseBody, &loginResponse)

	c.token = loginResponse.Token
	return nil
}

type RandomResponse struct {
	Value int `json:"value"`
}

func (c *Client) Random() (int, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/random", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("server invalid response")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	randomResponse := RandomResponse{}
	json.Unmarshal(responseBody, &randomResponse)

	return randomResponse.Value, nil
}

type SeedRequest struct {
	Seed int `json:"seed"`
}

func (c *Client) SetSeed(seed int) error {
	data := SeedRequest{
		Seed: seed,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.baseURL+"/seed", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("server invalid response")
	}

	return nil
}
