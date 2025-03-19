package test

import (
	"testing"
)

// 这个文件作为测试入口点，确保所有拆分后的测试文件都能被正确执行
// 实际的测试用例已经被拆分到各个模块的测试文件中：
// - user_test.go: 用户模块测试
// - house_test.go: 房源模块测试
// - favorite_test.go: 收藏模块测试
// - viewing_test.go: 预约看房模块测试
// - landlord_test.go: 房东模块测试

func TestMain(m *testing.M) {
	// 可以在这里进行全局测试设置
	m.Run()
	// 可以在这里进行全局测试清理
}