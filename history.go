package main

import "time"

// 历史记录项
type HistoryItem struct {
	Content string
	Time    time.Time
}

// History 管理剪贴板历史记录
type History struct {
	items []HistoryItem
}

// NewHistory 创建新的历史记录管理器
func NewHistory() *History {
	return &History{
		items: make([]HistoryItem, 0, maxHistoryItems),
	}
}

// Add 添加新的历史记录项
func (h *History) Add(content string) bool {
	if content == "" {
		return false
	}

	// 去重检查
	for _, item := range h.items {
		if item.Content == content {
			return false
		}
	}

	// 添加新项目到开头
	newItem := HistoryItem{
		Content: content,
		Time:    time.Now(),
	}
	h.items = append([]HistoryItem{newItem}, h.items...)

	// 限制最大数量
	if len(h.items) > maxHistoryItems {
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
	h.items = append(h.items[:index], h.items[index+1:]...)
	return true
}

// Clear 清空所有历史记录
func (h *History) Clear() {
	h.items = make([]HistoryItem, 0, maxHistoryItems)
}

// Len 获取历史记录数量
func (h *History) Len() int {
	return len(h.items)
}

// GetAll 获取所有历史记录
func (h *History) GetAll() []HistoryItem {
	return h.items
}
