# 剪贴板历史管理器 📋

一个基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 框架开发的 macOS 剪贴板历史管理工具。

![Demo](https://img.shields.io/badge/Platform-macOS-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.19+-green.svg)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

## 功能特性 ✨

- 🔄 **实时监控**：自动监控系统剪贴板变化
- 📚 **历史记录**：保存最近 50 条剪贴板历史
- 🔢 **编号显示**：每页显示 0-9 编号，方便快速识别
- ⏰ **时间戳**：显示每条记录的复制时间
- 📄 **分页浏览**：每页显示 10 条记录，支持翻页操作
- 🚀 **快速操作**：支持数字键 0-9 直接选择复制
- 🎨 **美观界面**：基于 TUI 的现代化界面
- 🧹 **智能去重**：自动去除重复内容
- 💾 **一键复制**：快速将历史记录复制回剪贴板

## 安装方法 🚀

### 方法一：直接下载可执行文件

1. 下载预编译的二进制文件（见 Releases 页面）
2. 将文件移动到系统 PATH 中：
   ```bash
   sudo cp clipboard-manager /usr/local/bin/
   chmod +x /usr/local/bin/clipboard-manager
   ```

### 方法二：从源码编译

确保你已经安装了 Go 1.19 或更高版本：

```bash
# 克隆项目
git clone <your-repo-url>
cd bubble-tea

# 安装依赖
go mod tidy

# 编译
go build -o clipboard-manager main.go

# 可选：安装到系统路径
sudo cp clipboard-manager /usr/local/bin/
```

### 方法三：使用 Go install（如果发布到 GitHub）

```bash
go install github.com/your-username/clipboard-manager@latest
```

## 使用方法 🎮

### 启动程序

```bash
clipboard-manager
```

或者如果在当前目录：

```bash
./clipboard-manager
```

### 快捷键操作

| 快捷键 | 功能 |
|--------|------|
| `↑` / `k` | 向上移动选择 |
| `↓` / `j` | 向下移动选择 |
| `←` / `h` | 上一页 |
| `→` / `l` | 下一页 |
| `0-9` | 直接选择对应编号项并复制 |
| `Home` | 跳转到第一项 |
| `End` | 跳转到最后一项 |
| `Enter` / `Space` | 复制选中项到剪贴板 |
| `d` / `Delete` | 删除选中项 |
| `c` | 清空所有历史记录 |
| `q` / `Ctrl+C` / `Esc` | 退出程序 |

### 界面说明

程序界面显示格式：
```
▶ [0] 复制的内容... | 15:30:42
  [1] 另一条内容... | 15:25:18
  [2] 更多内容...   | 15:20:05
```

- `▶` 表示当前选中的项目
- `[0-9]` 是页面内编号，可直接按数字键选择
- 右侧显示复制时间（时:分:秒格式）
- 底部显示分页信息和总记录数

### 使用技巧

1. **快速复制**：直接按数字键 0-9 即可复制对应项目
2. **分页浏览**：使用 ←→ 或 h/l 键在页面间快速切换
3. **精确定位**：使用 Home/End 键快速跳到首尾记录
4. **批量管理**：使用 c 键清空所有历史，d 键删除单个记录

### 使用场景

1. **开发编程**：在代码片段间快速切换，按编号直接复制
2. **文档编辑**：复用常用的文本内容，查看复制时间
3. **日常办公**：管理复制的链接、文本等，分页浏览历史
4. **数据处理**：在不同数据间快速切换，编号辅助选择

## 系统要求 📋

- **操作系统**：macOS 10.12 或更高版本
- **依赖**：程序使用系统自带的 `pbpaste` 和 `pbcopy` 命令
- **终端**：支持 ANSI 颜色的终端（推荐使用 iTerm2 或 Terminal.app）

## 高级用法 🔧

### 设置别名

在你的 shell 配置文件中添加别名：

```bash
# ~/.zshrc 或 ~/.bashrc
alias cb='clipboard-manager'
alias clipboard='clipboard-manager'
```

### 开机自启动

创建 LaunchAgent 配置文件：

```bash
# 创建配置文件
mkdir -p ~/Library/LaunchAgents
cat > ~/Library/LaunchAgents/com.user.clipboard-manager.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.clipboard-manager</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/clipboard-manager</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# 加载配置
launchctl load ~/Library/LaunchAgents/com.user.clipboard-manager.plist
```

## 技术架构 🛠️

- **框架**：[Bubble Tea](https://github.com/charmbracelet/bubbletea) - Go 语言的 TUI 框架
- **样式**：[Lip Gloss](https://github.com/charmbracelet/lipgloss) - 样式和布局库
- **语言**：Go 1.19+
- **架构模式**：基于 Elm Architecture 的响应式架构

## 项目结构 📁

```
.
├── README.md          # 项目文档
├── main.go           # 主程序文件
├── go.mod            # Go 模块文件
├── go.sum            # 依赖锁定文件
└── scripts/          # 构建脚本
    ├── build.sh      # 编译脚本
    └── install.sh    # 安装脚本
```

## 开发指南 👨‍💻

### 本地开发

```bash
# 克隆项目
git clone <your-repo>
cd bubble-tea

# 安装依赖
go mod tidy

# 运行开发版本
go run main.go
```

### 贡献代码

1. Fork 这个项目
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 常见问题 ❓

### Q: 程序无法启动或报错？
A: 确保你的 macOS 版本支持 `pbpaste` 和 `pbcopy` 命令，这些是系统自带的剪贴板工具。

### Q: 为什么有些内容没有被记录？
A: 程序会自动去重，相同的内容不会重复记录。同时空内容也会被忽略。

### Q: 如何修改最大历史记录数量？
A: 目前需要修改源码中的 `maxHistoryItems` 常量，未来版本将支持配置文件。

### Q: 支持其他操作系统吗？
A: 目前只支持 macOS，因为使用了 macOS 特有的剪贴板命令。

## 更新日志 📝

### v1.0.0
- ✨ 初始版本发布
- 🎯 支持剪贴板历史记录管理
- 🎨 美观的 TUI 界面
- ⚡ 实时剪贴板监控

## 许可证 📄

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详细信息。

## 致谢 🙏

- [Charm](https://github.com/charmbracelet) - 提供了优秀的 TUI 开发工具
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 强大的 TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - 漂亮的样式库

---

如果这个项目对你有帮助，请给它一个 ⭐️！

## 贡献 🤝

本项目的开发得到了 [GitHub Copilot](https://github.com/features/copilot) 的协助，AI 编程助手为项目的快速开发和代码优化提供了重要支持。