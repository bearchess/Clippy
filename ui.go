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
	history         *History
	clipboard       Clipboard
	cursor          int
	page            int
	expandedItem    int    // -1 表示没有展开的项目
	err             error
	last            string
	styles          *Styles
	historyFilePath string // 历史记录文件路径
	searching       bool   // 是否处于搜索模式
	searchQuery     string // 搜索查询字符串
	filteredIndices []int  // 过滤后的索引列表
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
	clipboard, err := NewClipboard()
	if err != nil {
		// 如果无法创建剪贴板，使用 nil，但会在运行时报错
		panic(fmt.Sprintf("无法初始化剪贴板: %v", err))
	}
	
	history := NewHistory()
	historyFilePath := GetHistoryFilePath()
	
	// 尝试加载历史记录，失败时静默处理
	_ = history.Load(historyFilePath)
	
	return Model{
		history:         history,
		clipboard:       clipboard,
		cursor:          0,
		page:            0,
		expandedItem:    -1,
		styles:          NewStyles(),
		historyFilePath: historyFilePath,
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
				// 自动保存历史记录
				_ = m.history.Save(m.historyFilePath)
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
	// 搜索模式下的特殊处理
	if m.searching {
		return m.handleSearchKey(key)
	}

	switch key.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "/":
		return m.enterSearchMode(), nil
		return m.enterSearchMode(), nil
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

// enterSearchMode 进入搜索模式
func (m Model) enterSearchMode() Model {
	m.searching = true
	m.searchQuery = ""
	m.filteredIndices = []int{}
	return m
}

// exitSearchMode 退出搜索模式
func (m Model) exitSearchMode() Model {
	m.searching = false
	m.searchQuery = ""
	m.filteredIndices = []int{}
	return m
}

// handleSearchKey 处理搜索模式下的按键
func (m Model) handleSearchKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc", "ctrl+c":
		return m.exitSearchMode(), nil
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.updateSearchResults()
		}
		return m, nil
	case "enter":
		// 在搜索模式下可以复制选中项
		return m.copySelected(), nil
	case "up", "k":
		return m.moveUp(), nil
	case "down", "j":
		return m.moveDown(), nil
	default:
		// 添加可打印字符到搜索查询（包括空格）
		if len(key.String()) == 1 {
			m.searchQuery += key.String()
			m.updateSearchResults()
		}
		return m, nil
	}
}

// updateSearchResults 更新搜索结果
func (m *Model) updateSearchResults() {
	m.filteredIndices = []int{}
	if m.searchQuery == "" {
		return
	}

	// 大小写不敏感的搜索
	query := strings.ToLower(m.searchQuery)
	for i := 0; i < m.history.Len(); i++ {
		if item, ok := m.history.Get(i); ok {
			if strings.Contains(strings.ToLower(item.Content), query) {
				m.filteredIndices = append(m.filteredIndices, i)
			}
		}
	}

	// 重置光标和页码
	m.cursor = 0
	m.page = 0
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
	maxIndex := m.history.Len() - 1
	if m.searching && len(m.filteredIndices) > 0 {
		maxIndex = len(m.filteredIndices) - 1
	}
	
	if maxIndex >= 0 && m.cursor < maxIndex {
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
	// 在搜索模式下，使用过滤后的索引
	var realIndex int
	if m.searching && len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices) {
		realIndex = m.filteredIndices[m.cursor]
	} else {
		realIndex = m.cursor
	}

	if item, ok := m.history.Get(realIndex); ok {
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
		// 保存历史记录
		_ = m.history.Save(m.historyFilePath)
	}
	return m
}

func (m Model) clearAll() Model {
	m.history.Clear()
	m.cursor = 0
	m.page = 0
	m.expandedItem = -1
	m.err = nil
	// 保存历史记录
	_ = m.history.Save(m.historyFilePath)
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

	b.WriteString(m.styles.Title.Render("剪贴板历史管理器") + "\n\n")

	// 搜索模式提示
	if m.searching {
		searchPrompt := fmt.Sprintf("搜索: %s█", m.searchQuery)
		b.WriteString(m.styles.Cursor.Render(searchPrompt) + "\n")
		if len(m.filteredIndices) > 0 {
			b.WriteString(m.styles.Help.Render(fmt.Sprintf("找到 %d 条匹配结果", len(m.filteredIndices))) + "\n\n")
		} else if m.searchQuery != "" {
			b.WriteString(m.styles.Help.Render("无匹配结果") + "\n\n")
		} else {
			b.WriteString(m.styles.Help.Render("输入搜索内容，按 Esc 退出搜索") + "\n\n")
		}
	}

	if m.history.Len() == 0 {
		b.WriteString(m.styles.Normal.Render("剪贴板历史为空，请复制一些内容...\n"))
	} else {
		if m.searching && len(m.filteredIndices) > 0 {
			m.renderFilteredItems(&b)
		} else if m.searching {
			// 搜索模式下无结果时不显示历史记录
		} else {
			m.renderHistoryItems(&b)
			m.renderPageInfo(&b)
		}
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

func (m Model) renderFilteredItems(b *strings.Builder) {
	// 在搜索模式下显示过滤后的结果
	for displayIdx, realIdx := range m.filteredIndices {
		if displayIdx >= itemsPerPage {
			break // 只显示前 10 条
		}

		item, _ := m.history.Get(realIdx)
		displayIndex := displayIdx + 1
		if displayIndex == 10 {
			displayIndex = 0
		}

		prefix := "  "
		style := m.styles.Normal
		if displayIdx == m.cursor {
			prefix = "▶ "
			style = m.styles.Cursor
		}

		if m.expandedItem == realIdx {
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
	if m.searching {
		helpLines := []string{
			"搜索模式: 输入内容搜索   Backspace删除   Enter复制",
			"导航: ↑↓/kj移动   Esc退出搜索   Ctrl+C退出程序",
		}
		for _, line := range helpLines {
			b.WriteString(m.styles.Help.Render(line) + "\n")
		}
	} else {
		helpLines := []string{
			"导航: ↑↓/kj移动   ←→/hl翻页   Home/End首末项",
			"操作: 1-0快选复制   Enter/Space复制   v展开/收起",
			"管理: d删除   c清空   /搜索   q/Esc退出",
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
}

func (m Model) renderError(b *strings.Builder) {
	if m.err != nil {
		b.WriteString("\n" + m.styles.Error.Render("❌ 错误: "+m.err.Error()) + "\n")
	}
}
