package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewHistory(t *testing.T) {
	h := NewHistory()
	if h == nil {
		t.Fatal("NewHistory() 应该返回非 nil 的 History")
	}
	if h.Len() != 0 {
		t.Errorf("新的 History 长度应该为 0, 得到 %d", h.Len())
	}
}

func TestHistoryAdd(t *testing.T) {
	h := NewHistory()

	// 测试添加有效内容
	if !h.Add("测试内容1") {
		t.Error("添加有效内容应该返回 true")
	}
	if h.Len() != 1 {
		t.Errorf("添加后长度应该为 1, 得到 %d", h.Len())
	}

	// 测试添加空字符串
	if h.Add("") {
		t.Error("添加空字符串应该返回 false")
	}
	if h.Len() != 1 {
		t.Errorf("添加空字符串后长度应该保持为 1, 得到 %d", h.Len())
	}

	// 测试去重
	if h.Add("测试内容1") {
		t.Error("添加重复内容应该返回 false")
	}
	if h.Len() != 1 {
		t.Errorf("添加重复内容后长度应该保持为 1, 得到 %d", h.Len())
	}

	// 测试添加多个不同内容
	h.Add("测试内容2")
	h.Add("测试内容3")
	if h.Len() != 3 {
		t.Errorf("添加3个不同内容后长度应该为 3, 得到 %d", h.Len())
	}

	// 测试最新的内容在最前面
	if item, ok := h.Get(0); !ok || item.Content != "测试内容3" {
		t.Error("最新添加的内容应该在索引 0")
	}
}

func TestHistoryAddMaxItems(t *testing.T) {
	h := NewHistory()

	// 添加超过最大数量的项目
	for i := 0; i < maxHistoryItems+10; i++ {
		h.Add(string(rune('A' + i)))
	}

	if h.Len() != maxHistoryItems {
		t.Errorf("历史记录数量应该不超过 %d, 得到 %d", maxHistoryItems, h.Len())
	}

	// 验证最旧的项目被移除
	item, ok := h.Get(0)
	if !ok {
		t.Fatal("应该能获取索引 0 的项目")
	}
	// 最新的应该是 'A' + maxHistoryItems + 9
	expected := string(rune('A' + maxHistoryItems + 9))
	if item.Content != expected {
		t.Errorf("索引 0 应该是 '%s', 得到 '%s'", expected, item.Content)
	}
}

func TestHistoryGet(t *testing.T) {
	h := NewHistory()
	h.Add("内容1")
	h.Add("内容2")

	// 测试有效索引
	item, ok := h.Get(0)
	if !ok {
		t.Error("Get(0) 应该返回 true")
	}
	if item.Content != "内容2" {
		t.Errorf("Get(0) 应该返回 '内容2', 得到 '%s'", item.Content)
	}

	// 测试无效索引
	_, ok = h.Get(-1)
	if ok {
		t.Error("Get(-1) 应该返回 false")
	}

	_, ok = h.Get(10)
	if ok {
		t.Error("Get(10) 应该返回 false")
	}
}

func TestHistoryDelete(t *testing.T) {
	h := NewHistory()
	h.Add("内容1")
	h.Add("内容2")
	h.Add("内容3")

	// 测试删除中间项
	if !h.Delete(1) {
		t.Error("Delete(1) 应该返回 true")
	}
	if h.Len() != 2 {
		t.Errorf("删除后长度应该为 2, 得到 %d", h.Len())
	}

	// 验证删除后的内容
	item, _ := h.Get(0)
	if item.Content != "内容3" {
		t.Errorf("索引 0 应该是 '内容3', 得到 '%s'", item.Content)
	}
	item, _ = h.Get(1)
	if item.Content != "内容1" {
		t.Errorf("索引 1 应该是 '内容1', 得到 '%s'", item.Content)
	}

	// 测试删除无效索引
	if h.Delete(-1) {
		t.Error("Delete(-1) 应该返回 false")
	}
	if h.Delete(10) {
		t.Error("Delete(10) 应该返回 false")
	}

	// 验证 map 也被更新（之前删除的 "内容2" 应该可以重新添加）
	if !h.Add("内容2") {
		t.Error("之前删除的内容应该可以重新添加")
	}
}

func TestHistoryClear(t *testing.T) {
	h := NewHistory()
	h.Add("内容1")
	h.Add("内容2")
	h.Add("内容3")

	h.Clear()

	if h.Len() != 0 {
		t.Errorf("Clear() 后长度应该为 0, 得到 %d", h.Len())
	}

	// 验证清空后可以重新添加之前的内容
	if !h.Add("内容1") {
		t.Error("Clear() 后应该可以重新添加之前的内容")
	}
}

func TestHistorySaveAndLoad(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_history.json")

	// 创建并填充历史记录
	h1 := NewHistory()
	h1.Add("内容1")
	time.Sleep(time.Millisecond) // 确保时间戳不同
	h1.Add("内容2")
	time.Sleep(time.Millisecond)
	h1.Add("内容3")

	// 保存
	if err := h1.Save(tmpFile); err != nil {
		t.Fatalf("Save() 失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("保存的文件不存在")
	}

	// 加载到新的历史记录
	h2 := NewHistory()
	if err := h2.Load(tmpFile); err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	// 验证内容
	if h2.Len() != h1.Len() {
		t.Errorf("加载后长度应该为 %d, 得到 %d", h1.Len(), h2.Len())
	}

	for i := 0; i < h1.Len(); i++ {
		item1, _ := h1.Get(i)
		item2, _ := h2.Get(i)
		if item1.Content != item2.Content {
			t.Errorf("索引 %d 的内容不匹配: 期望 '%s', 得到 '%s'", i, item1.Content, item2.Content)
		}
	}

	// 验证去重 map 也被正确加载
	if h2.Add("内容1") {
		t.Error("加载后应该保持去重功能，不能添加重复内容")
	}
}

func TestHistoryLoadNonExistentFile(t *testing.T) {
	h := NewHistory()
	err := h.Load("/nonexistent/path/to/file.json")
	if err == nil {
		t.Error("加载不存在的文件应该返回错误")
	}
}

func TestGetHistoryFilePath(t *testing.T) {
	path := GetHistoryFilePath()
	if path == "" {
		t.Error("GetHistoryFilePath() 不应该返回空字符串")
	}

	// 验证路径包含预期的组件
	if !filepath.IsAbs(path) && path != "clippy_history.json" {
		t.Errorf("路径应该是绝对路径或默认文件名, 得到 '%s'", path)
	}
}

func TestHistoryGetAll(t *testing.T) {
	h := NewHistory()
	h.Add("内容1")
	h.Add("内容2")
	h.Add("内容3")

	all := h.GetAll()
	if len(all) != 3 {
		t.Errorf("GetAll() 应该返回 3 个项目, 得到 %d", len(all))
	}

	// 验证顺序（最新的在前）
	if all[0].Content != "内容3" || all[1].Content != "内容2" || all[2].Content != "内容1" {
		t.Error("GetAll() 返回的项目顺序不正确")
	}
}

func TestHistoryTime(t *testing.T) {
	h := NewHistory()
	before := time.Now()
	h.Add("测试内容")
	after := time.Now()

	item, _ := h.Get(0)
	if item.Time.Before(before) || item.Time.After(after) {
		t.Error("添加的项目时间戳应该在添加前后之间")
	}
}
