package main

import (
	"bytes"
	"os/exec"
	"strings"
)

// Clipboard 剪贴板操作接口
type Clipboard interface {
	Get() (string, error)
	Set(text string) error
}

// MacClipboard macOS剪贴板实现
type MacClipboard struct{}

// NewMacClipboard 创建macOS剪贴板实例
func NewMacClipboard() *MacClipboard {
	return &MacClipboard{}
}

// Get 获取剪贴板内容
func (c *MacClipboard) Get() (string, error) {
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// Set 设置剪贴板内容
func (c *MacClipboard) Set(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
