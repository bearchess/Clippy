package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model Bubble Tea 模型
type Model struct {
	history      *History
	clipboard    Clipboard
	cursor       int
	page         int
	expandedItem int // -1 表示没有展开的项目
	err          error
	last         string
	styles       *Styles
}

// Styles 样式定义
type Styles struct {
	Title  lipgloss.Style
	Cursor lipgloss.Style
	Normal lipgloss.Style
	Error  lipgloss.Style
	Help   lipgloss.Style
}

// NewStyles 创建新的样式
func NewStyles() *Styles {
	return &Styles{
		Title:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")),
		Cursor: lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		Normal: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Error:  lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		Help:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
}

// NewModel 创建新的模型
func NewModel() Model {
	return Model{
		history:      NewHistory(),
		clipboard:    NewMacClipboard(),
		cursor:       0,
		page:         0,
		expandedItem: -1,
		styles:       NewStyles(),
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return m.pollClipboard()
}

// pollClipboard 轮询剪贴板
func (m Model) pollClipboard() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Duration(pollInterval) * time.Millisecond)
		clip, err := m.clipboard.Get()
		if err != nil {
			return err
		}
		return clip
	}
}

// Update 处理消息
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case error:
		m.err = v
		return m, m.pollClipboard()
	case string:
		if v != "" && v != m.last {
			m.last = v
			if m.history.Add(v) {
				m.cursor = 0
				m.page = 0
			}
		}
		return m, m.pollClipboard()
	case tea.KeyMsg:
		return m.handleKeyPress(v)
	}
	return m, nil
}

// handleKeyPress 处理按键事件
func (m Model) handleKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		return m.moveUp(), nil
	case "down", "j":
		return m.moveDown(), nil
	case "left", "h":
		return m.prevPage(), nil
	case "right", "l":
		return m.nextPage(), nil
	case "home":
		return m.goHome(), nil
	case "end":
		return m.goEnd(), nil
	case "enter", " ":
		return m.copySelected(), nil
	case "d", "delete":
		return m.deleteSelected(), nil
	case "c":
		return m.clearAll(), nil
	case "v":
		return m.toggleExpand(), nil
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return m.selectAndCopy(key.String()), nil
	}
	return m, nil
}

// 导航方法
func (m Model) moveUp() Model {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.page*itemsPerPage {
			m.page--
		}
	}
	return m
}

func (m Model) moveDown() Model {
	if m.history.Len() > 0 && m.cursor < m.history.Len()-1 {
		m.cursor++
		if m.cursor >= (m.page+1)*itemsPerPage {
			m.page++
		}
	}
	return m
}

func (m Model) prevPage() Model {
	if m.page > 0 {
		m.page--
		m.cursor = m.page * itemsPerPage
	}
	return m
}

func (m Model) nextPage() Model {
	totalPages := (m.history.Len() + itemsPerPage - 1) / itemsPerPage
	if m.page < totalPages-1 {
		m.page++
		m.cursor = m.page * itemsPerPage
		if m.cursor >= m.history.Len() {
			m.cursor = m.history.Len() - 1
		}
	}
	return m
}

func (m Model) goHome() Model {
	m.cursor = 0
	m.page = 0
	return m
}

func (m Model) goEnd() Model {
	if m.history.Len() > 0 {
		m.cursor = m.history.Len() - 1
		m.page = (m.history.Len() - 1) / itemsPerPage
	}
	return m
}

// 操作方法
func (m Model) copySelected() Model {
	if item, ok := m.history.Get(m.cursor); ok {
		if err := m.clipboard.Set(item.Content); err != nil {
			m.err = err
		} else {
			m.err = nil
		}
	}
	return m
}

func (m Model) deleteSelected() Model {
	if m.history.Delete(m.cursor) {
		if m.cursor >= m.history.Len() && m.cursor > 0 {
			m.cursor--
		}
		// 调整页码
		totalPages := (m.history.Len() + itemsPerPage - 1) / itemsPerPage
		if totalPages == 0 {
			m.page = 0
		} else if m.page >= totalPages {
			m.page = totalPages - 1
		}
		m.err = nil
	}
	return m
}

func (m Model) clearAll() Model {
	m.history.Clear()
	m.cursor = 0
	m.page = 0
	m.expandedItem = -1
	m.err = nil
	return m
}

func (m Model) toggleExpand() Model {
	if m.history.Len() > 0 {
		if m.expandedItem == m.cursor {
			m.expandedItem = -1
		} else {
			m.expandedItem = m.cursor
		}
	}
	return m
}

func (m Model) selectAndCopy(keyStr string) Model {
	digit := int(keyStr[0] - '0')
	var pageIndex int
	if digit == 0 {
		pageIndex = 9 // 按键0对应第10项（索引9）
	} else {
		pageIndex = digit - 1 // 按键1-9对应索引0-8
	}

	index := m.page*itemsPerPage + pageIndex
	if item, ok := m.history.Get(index); ok {
		if err := m.clipboard.Set(item.Content); err != nil {
			m.err = err
		} else {
			m.err = nil
			m.cursor = index
		}
	}
	return m
}

// View 渲染界面
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render("剪贴板历史管理器 (macOS)") + "\n\n")

	if m.history.Len() == 0 {
		b.WriteString(m.styles.Normal.Render("剪贴板历史为空，请复制一些内容...\n"))
	} else {
		m.renderHistoryItems(&b)
		m.renderPageInfo(&b)
	}

	m.renderHelp(&b)
	m.renderError(&b)

	return b.String()
}

func (m Model) renderHistoryItems(b *strings.Builder) {
	totalPages := (m.history.Len() + itemsPerPage - 1) / itemsPerPage

	// 使用局部变量而不是直接修改 m.page（值接收者）
	page := m.page
	if page < 0 {
		page = 0
	} else if page >= totalPages {
		page = totalPages - 1
	}

	start := page * itemsPerPage
	end := start + itemsPerPage
	if end > m.history.Len() {
		end = m.history.Len()
	}

	for i := start; i < end; i++ {
		item, _ := m.history.Get(i)
		pageIndex := i - start
		displayIndex := pageIndex + 1
		if displayIndex == 10 {
			displayIndex = 0
		}

		prefix := "  "
		style := m.styles.Normal
		if i == m.cursor {
			prefix = "▶ "
			style = m.styles.Cursor
		}

		if m.expandedItem == i {
			m.renderExpandedItem(b, item, prefix, displayIndex, style)
		} else {
			m.renderCompactItem(b, item, prefix, displayIndex, style)
		}
	}
}

func (m Model) renderExpandedItem(b *strings.Builder, item HistoryItem, prefix string, displayIndex int, style lipgloss.Style) {
	timeStr := item.Time.Format("15:04:05")
	expandedPrefix := fmt.Sprintf("%s[%d] ", prefix, displayIndex)

	lines := strings.Split(item.Content, "\n")
	for lineIdx, line := range lines {
		if lineIdx == 0 {
			firstLine := fmt.Sprintf("%s%s | %s",
				expandedPrefix,
				line,
				m.styles.Help.Render(timeStr))
			b.WriteString(style.Render(firstLine) + "\n")
		} else {
			indentLen := len(expandedPrefix)
			indent := strings.Repeat(" ", indentLen)
			b.WriteString(style.Render(indent+line) + "\n")
		}
	}
}

func (m Model) renderCompactItem(b *strings.Builder, item HistoryItem, prefix string, displayIndex int, style lipgloss.Style) {
	display := strings.ReplaceAll(item.Content, "\n", "\\n")
	display = strings.ReplaceAll(display, "\t", "\\t")

	timeStr := item.Time.Format("15:04:05")
	prefixLen := len(fmt.Sprintf("%s[%d] ", prefix, displayIndex))
	timeLen := len(" | " + timeStr)
	availableLen := maxDisplayLength - prefixLen - timeLen

	if availableLen > 3 && len(display) > availableLen {
		display = display[:availableLen-3] + truncateSuffix
	} else if availableLen > 0 && len(display) > availableLen {
		display = display[:availableLen]
	}

	line := fmt.Sprintf("%s[%d] %s | %s",
		prefix,
		displayIndex,
		display,
		m.styles.Help.Render(timeStr))

	b.WriteString(style.Render(line) + "\n")
}

func (m Model) renderPageInfo(b *strings.Builder) {
	totalPages := (m.history.Len() + itemsPerPage - 1) / itemsPerPage
	pageInfo := fmt.Sprintf("第 %d/%d 页 (共 %d 条记录)", m.page+1, totalPages, m.history.Len())
	b.WriteString("\n" + m.styles.Help.Render(pageInfo) + "\n")
}

func (m Model) renderHelp(b *strings.Builder) {
	b.WriteString("\n")
	helpLines := []string{
		"导航: ↑↓/kj移动   ←→/hl翻页   Home/End首末项",
		"操作: 1-0快选复制   Enter/Space复制   v展开/收起",
		"管理: d删除   c清空   q/Esc退出",
	}

	for _, line := range helpLines {
		b.WriteString(m.styles.Help.Render(line) + "\n")
	}

	// 显示展开提示
	if m.history.Len() > 0 {
		if m.expandedItem == m.cursor {
			b.WriteString(m.styles.Help.Render("💡 当前项已展开，按 v 收起") + "\n")
		} else {
			b.WriteString(m.styles.Help.Render("💡 按 v 展开当前项查看完整内容") + "\n")
		}
	}
}

func (m Model) renderError(b *strings.Builder) {
	if m.err != nil {
		b.WriteString("\n" + m.styles.Error.Render("❌ 错误: "+m.err.Error()) + "\n")
	}
}
