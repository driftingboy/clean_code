package memento

import (
	"errors"
	"strings"
)

// InputOperater 用于保存数据
type InputOperater struct {
	content strings.Builder
	// 全量备份
	snapshots []string
}

// Append 追加数据
func (in *InputOperater) Append(content string) {
	in.content.WriteString(content)
}

// GetText 获取数据
func (in *InputOperater) View() string {
	return in.content.String()
}

// Snapshot 创建快照
func (in *InputOperater) Snapshot() {
	in.snapshots = append(in.snapshots, in.content.String())
}

// Restore 从快照中恢复指定版本
func (in *InputOperater) Restore(version int) error {
	if version > len(in.snapshots) {
		return errors.New("the version not exist")
	}
	in.content.Reset()
	if _, err := in.content.WriteString(in.snapshots[version]); err != nil {
		return err
	}
	return nil
}

// RestoreLast 恢复最后一次快照数据
func (in *InputOperater) RestoreLast() error {
	if len(in.snapshots) == 0 {
		return errors.New("no snapshot")
	}
	in.content.Reset()
	if _, err := in.content.WriteString(in.snapshots[len(in.snapshots)-1]); err != nil {
		return err
	}
	return nil
}

// 管理 snapshots
type SnapshotHolder struct {
	snapshots []string
}

// Get 获取指定版本快照
func (s *SnapshotHolder) Get(version int) string {
	if version > len(s.snapshots) {
		return ""
	}
	return s.snapshots[version]
}

// Pop 弹出最新快照
func (s *SnapshotHolder) Pop() string {
	if len(s.snapshots) == 0 {
		return ""
	}
	result := s.snapshots[len(s.snapshots)-1]
	s.snapshots = s.snapshots[:len(s.snapshots)-1]
	return result
}

// ListVersion 列出所有版本号
func (s *SnapshotHolder) ListVersion() []int {
	result := make([]int, len(s.snapshots))
	for i := range s.snapshots {
		result = append(result, i)
	}
	return result
}
