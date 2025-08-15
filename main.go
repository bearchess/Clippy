package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 解析命令行参数
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.Parse()

	if showVersion {
		fmt.Printf("版本: %s\n构建时间: %s\nGo 版本: %s\n", Version, BuildTime, GoVersion)
		return
	}

	// 创建并运行程序
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("启动失败: %v\n", err)
		os.Exit(1)
	}
}
