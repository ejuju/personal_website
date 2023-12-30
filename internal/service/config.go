package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Env string `json:"env"`
	Log string `json:"log"`
}

func LoadConfig(fpath string) (c Config, err error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0400)
	if err != nil {
		return c, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return c, fmt.Errorf("decode JSON: %w", err)
	}
	return c, nil
}
