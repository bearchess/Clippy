package main

import (
	"runtime"
	"testing"
)

func TestNewClipboard(t *testing.T) {
	clipboard, err := NewClipboard()

	switch runtime.GOOS {
	case "darwin":
		if err != nil {
			t.Errorf("在 macOS 上创建剪贴板不应该失败: %v", err)
		}
		if _, ok := clipboard.(*MacClipboard); !ok {
			t.Error("在 macOS 上应该返回 MacClipboard")
		}
	case "linux":
		if err != nil {
			t.Errorf("在 Linux 上创建剪贴板不应该失败: %v", err)
		}
		if _, ok := clipboard.(*LinuxClipboard); !ok {
			t.Error("在 Linux 上应该返回 LinuxClipboard")
		}
	case "windows":
		if err != nil {
			t.Errorf("在 Windows 上创建剪贴板不应该失败: %v", err)
		}
		if _, ok := clipboard.(*WindowsClipboard); !ok {
			t.Error("在 Windows 上应该返回 WindowsClipboard")
		}
	default:
		if err == nil {
			t.Error("在不支持的平台上应该返回错误")
		}
		if clipboard != nil {
			t.Error("在不支持的平台上应该返回 nil 剪贴板")
		}
	}
}

func TestNewMacClipboard(t *testing.T) {
	clipboard := NewMacClipboard()
	if clipboard == nil {
		t.Error("NewMacClipboard() 不应该返回 nil")
	}
}

func TestNewLinuxClipboard(t *testing.T) {
	clipboard := NewLinuxClipboard()
	if clipboard == nil {
		t.Error("NewLinuxClipboard() 不应该返回 nil")
	}
}

func TestNewWindowsClipboard(t *testing.T) {
	clipboard := NewWindowsClipboard()
	if clipboard == nil {
		t.Error("NewWindowsClipboard() 不应该返回 nil")
	}
}

// 注意：以下测试需要实际的剪贴板命令可用，在 CI 环境中可能会失败
// 这些测试主要用于本地开发验证

func TestMacClipboardGetSet(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("跳过 macOS 特定测试")
	}

	clipboard := NewMacClipboard()
	testContent := "测试内容 MacClipboard"

	// 测试 Set
	err := clipboard.Set(testContent)
	if err != nil {
		t.Skipf("Set() 失败（可能 pbcopy 不可用）: %v", err)
	}

	// 测试 Get
	content, err := clipboard.Get()
	if err != nil {
		t.Skipf("Get() 失败（可能 pbpaste 不可用）: %v", err)
	}

	if content != testContent {
		t.Errorf("Get() 返回的内容不匹配: 期望 '%s', 得到 '%s'", testContent, content)
	}
}

func TestLinuxClipboardGetSet(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("跳过 Linux 特定测试")
	}

	clipboard := NewLinuxClipboard()
	testContent := "测试内容 LinuxClipboard"

	// 测试 Set
	err := clipboard.Set(testContent)
	if err != nil {
		t.Skipf("Set() 失败（可能 xclip 不可用）: %v", err)
	}

	// 测试 Get
	content, err := clipboard.Get()
	if err != nil {
		t.Skipf("Get() 失败（可能 xclip 不可用）: %v", err)
	}

	if content != testContent {
		t.Errorf("Get() 返回的内容不匹配: 期望 '%s', 得到 '%s'", testContent, content)
	}
}

func TestWindowsClipboardGetSet(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("跳过 Windows 特定测试")
	}

	clipboard := NewWindowsClipboard()
	testContent := "测试内容 WindowsClipboard"

	// 测试 Set
	err := clipboard.Set(testContent)
	if err != nil {
		t.Skipf("Set() 失败（可能 clip.exe 不可用）: %v", err)
	}

	// 测试 Get
	content, err := clipboard.Get()
	if err != nil {
		t.Skipf("Get() 失败（可能 powershell 不可用）: %v", err)
	}

	if content != testContent {
		t.Errorf("Get() 返回的内容不匹配: 期望 '%s', 得到 '%s'", testContent, content)
	}
}

func TestClipboardInterface(t *testing.T) {
	// 验证所有剪贴板类型都实现了 Clipboard 接口
	var _ Clipboard = (*MacClipboard)(nil)
	var _ Clipboard = (*LinuxClipboard)(nil)
	var _ Clipboard = (*WindowsClipboard)(nil)
}
