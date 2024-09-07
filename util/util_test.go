package util

import (
	"fmt"
	"testing"
)

func TestBranch(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		a string
	}{
		{"."},
		{"F:/git-tmp/neovim-for-beginner"},
	}

	for _, tt := range tests {
		result, err := BranchList(tt.a)
		if err != nil {
			t.Errorf("error %v", err)
		}
		fmt.Printf("list %q", result)
	}
}
