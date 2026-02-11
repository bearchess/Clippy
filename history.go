package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// 历史记录项
type HistoryItem struct {
	Content string
	Time    time.Time
}

// History 管理剪贴板历史记录
type History struct {
	items []HistoryItem
	seen  map[string]struct{} // 用于 O(1) 去重查找
}

// NewHistory 创建新的历史记录管理器
func NewHistory() *History {
	return &History{
		items: make([]HistoryItem, 0, maxHistoryItems),
		seen:  make(map[string]struct{}),
	}
}

// Add 添加新的历史记录项
func (h *History) Add(content string) bool {
	if content == "" {
		return false
	}

	// O(1) 去重检查
	if _, exists := h.seen[content]; exists {
		return false
	}

	// 添加新项目到开头
	newItem := HistoryItem{
		Content: content,
		Time:    time.Now(),
	}
	
	// 优化的头部插入：先 append 预留空间，再 copy 后移
	h.items = append(h.items, HistoryItem{})
	copy(h.items[1:], h.items)
	h.items[0] = newItem
	h.seen[content] = struct{}{}

	// 限制最大数量
	if len(h.items) > maxHistoryItems {
		// 移除最后一项并从 map 中删除
		removed := h.items[maxHistoryItems]
		delete(h.seen, removed.Content)
		h.items = h.items[:maxHistoryItems]
	}

	return true
}

// Get 获取指定索引的历史记录
func (h *History) Get(index int) (HistoryItem, bool) {
	if index < 0 || index >= len(h.items) {
		return HistoryItem{}, false
	}
	return h.items[index], true
}

// Delete 删除指定索引的历史记录
func (h *History) Delete(index int) bool {
	if index < 0 || index >= len(h.items) {
		return false
	}
	// 从 map 中删除
	delete(h.seen, h.items[index].Content)
	h.items = append(h.items[:index], h.items[index+1:]...)
	return true
}

// Clear 清空所有历史记录
func (h *History) Clear() {
	h.items = make([]HistoryItem, 0, maxHistoryItems)
	h.seen = make(map[string]struct{})
}

// Len 获取历史记录数量
func (h *History) Len() int {
	return len(h.items)
}

// GetAll 获取所有历史记录
func (h *History) GetAll() []HistoryItem {
	return h.items
}

// Save 保存历史记录到文件
func (h *History) Save(filepath string) error {
	// 确保目录存在
	dir := filepath[:len(filepath)-len(filepath[len(filepath)-1:])-1]
	for i := len(filepath) - 1; i >= 0; i-- {
		if filepath[i] == '/' || filepath[i] == '\\' {
			dir = filepath[:i]
			break
		}
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(h.items, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filepath, data, 0644)
}

// Load 从文件加载历史记录
func (h *History) Load(filepath string) error {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// 反序列化
	var items []HistoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	// 重建 items 和 seen map
	h.items = make([]HistoryItem, 0, maxHistoryItems)
	h.seen = make(map[string]struct{})

	for _, item := range items {
		if _, exists := h.seen[item.Content]; !exists {
			h.items = append(h.items, item)
			h.seen[item.Content] = struct{}{}
		}
	}

	return nil
}

// GetHistoryFilePath 获取默认历史文件路径
func GetHistoryFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果无法获取用户目录，使用当前目录
		return "clippy_history.json"
	}
	return filepath.Join(homeDir, ".config", "clippy", "history.json")
}
