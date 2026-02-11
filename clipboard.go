package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
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

// LinuxClipboard Linux剪贴板实现（使用xclip）
type LinuxClipboard struct{}

// NewLinuxClipboard 创建Linux剪贴板实例
func NewLinuxClipboard() *LinuxClipboard {
	return &LinuxClipboard{}
}

// Get 获取剪贴板内容
func (c *LinuxClipboard) Get() (string, error) {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// Set 设置剪贴板内容
func (c *LinuxClipboard) Set(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// WindowsClipboard Windows剪贴板实现
type WindowsClipboard struct{}

// NewWindowsClipboard 创建Windows剪贴板实例
func NewWindowsClipboard() *WindowsClipboard {
	return &WindowsClipboard{}
}

// Get 获取剪贴板内容
func (c *WindowsClipboard) Get() (string, error) {
	cmd := exec.Command("powershell.exe", "-Command", "Get-Clipboard")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// Set 设置剪贴板内容
func (c *WindowsClipboard) Set(text string) error {
	cmd := exec.Command("clip.exe")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// NewClipboard 根据运行平台创建相应的剪贴板实例
func NewClipboard() (Clipboard, error) {
	switch runtime.GOOS {
	case "darwin":
		return NewMacClipboard(), nil
	case "linux":
		return NewLinuxClipboard(), nil
	case "windows":
		return NewWindowsClipboard(), nil
	default:
		return nil, errors.New(fmt.Sprintf("不支持的操作系统: %s", runtime.GOOS))
	}
}
