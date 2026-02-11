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
