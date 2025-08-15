package main

// 版本信息，将在构建时通过 ldflags 注入
var (
	Version   = "v1.0.0"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

// 常量定义
const (
	maxHistoryItems  = 50
	maxDisplayLength = 50
	pollInterval     = 500 // 毫秒
	truncateSuffix   = "..."
	itemsPerPage     = 10 // 每页显示的项目数
)
