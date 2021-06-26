package auth

/*
- 根据accId获取pwd
*/
type Storage interface {
	QueryPwdByAccId(string) string
}

var _ Storage = (*MemAuthStorage)(nil)

type MemAuthStorage struct{}

func (as *MemAuthStorage) QueryPwdByAccId(id string) string {
	return ""
}
