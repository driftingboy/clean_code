package prototype

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var keyWordsMap = Keywords{
	"litao": {Word: "litao", Visit: 10, UpdatedAt: time.Date(2022, 1, 17, 10, 2, 0, 0, time.Local)},
	"aibo":  {Word: "aibo", Visit: 20, UpdatedAt: time.Date(2022, 1, 15, 6, 2, 0, 0, time.Local)},
	"test":  {Word: "test", Visit: 30, UpdatedAt: time.Date(2022, 1, 8, 12, 2, 0, 0, time.Local)},
}

var updateWords = []*Keyword{
	{Word: "litao", Visit: 22, UpdatedAt: time.Date(2022, 1, 16, 6, 2, 0, 0, time.Local)},
	{Word: "new-key", Visit: 5, UpdatedAt: time.Date(2022, 1, 17, 6, 2, 0, 0, time.Local)},
}

// 基础功能测试
func TestKeywords_Update(t *testing.T) {
	type args struct {
		updatedWords []*Keyword
	}
	tests := []struct {
		name  string
		words Keywords
		args  args
	}{
		{
			name: "normal", words: keyWordsMap, args: args{
				updatedWords: updateWords,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.words.Update(tt.args.updatedWords)
			assert.Equal(t, 22, tt.words.Get("litao").Visit)
			assert.Equal(t, 4, len(tt.words))
		})
	}
}

// 并发测试
func Test_Concurency(t *testing.T) {
	var wg sync.WaitGroup
	testKeyWordsMap := keyWordsMap
	stopChan := make(chan struct{})

	wg.Add(1)
	// 一直检测是否原数据会不会受影响
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				return
			default:
				// 检测旧版本是否一直可用
				assert.Equal(t, testKeyWordsMap.Get("litao").Visit, 10)
				assert.Equal(t, len(testKeyWordsMap), 3)
			}
		}
	}()

	// 更新数据 传入 chan 测试，在指针修改之前发出停止信号
	testKeyWordsMap.UpdateForTest(updateWords, stopChan)
	// 检测新版本值
	assert.Equal(t, testKeyWordsMap.Get("litao").Visit, 22)
	assert.Equal(t, len(testKeyWordsMap), 4)

	wg.Wait()
}

func (words *Keywords) UpdateForTest(updatedWords []*Keyword, stopChan chan struct{}) {
	// 提前分配内存
	newKeywords := make(Keywords, len(*words)+len(updatedWords))

	// 浅拷贝，将原有值索引拷贝
	for k, v := range *words {
		newKeywords[k] = v
	}

	// 替换掉或者插入新的字段，这里用的是深拷贝, newKeywords 都是更新的指针，而不会更新指针里的值
	for _, word := range updatedWords {
		// 延长时间看效果
		time.Sleep(50 * time.Microsecond)
		newKeywords[word.Word] = word.Clone()
	}

	stopChan <- struct{}{}
	// 指针指向新的 keywords 集合：指针赋值是原子操作, 并发安全
	*words = newKeywords
}

// 反汇编分析
