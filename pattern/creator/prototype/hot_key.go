package prototype

/*
	现在有一如下场景：
	需要在内存中维护一个 关键词集合
	- 要求记录每个关键词的 名字、使用次数、最后更新时间。
	- 要求更新集合时，不能影响集合的查询，不能即访问到新的数据又访问到旧的数据
	这个要求，最容易想到的解决方案就是加锁，但是锁的粒度大，会严重影响查询性能；
	所以我们可以利用copy on wire 的思路，空间换时间，拷贝一份新数据在上面修改，不影响原数据；修改完毕，原子的执行一条指针赋值操作就行。
*/

import (
	"encoding/json"
	"time"
)

// Keyword 搜索关键字
type Keyword struct {
	Word      string
	Visit     int
	UpdatedAt time.Time
}

// Clone 序列化与反序列化的方式深拷贝，涉及反射速度较慢
func (k *Keyword) Clone() *Keyword {
	var newKeyword Keyword
	b, _ := json.Marshal(k)
	_ = json.Unmarshal(b, &newKeyword)
	return &newKeyword
}

func (k *Keyword) CloneFast() *Keyword {
	return &Keyword{
		Word:      k.Word,
		Visit:     k.Visit,
		UpdatedAt: k.UpdatedAt,
	}
}

// Keywords 关键字 map
type Keywords map[string]*Keyword

func (words Keywords) Get(key string) Keyword {
	return *words[key]
}

// 更新 Keywords
// updatedWords: 需要更新的关键词列表，由于从数据库中获取数据常常是数组的方式
func (words *Keywords) Update(updatedWords []*Keyword) {
	// 提前分配内存
	newKeywords := make(Keywords, len(*words)+len(updatedWords))

	// 浅拷贝，将原有值索引拷贝
	for k, v := range *words {
		newKeywords[k] = v
	}

	// 替换掉或者插入新的字段，这里用的是深拷贝, newKeywords 都是更新的指针，而不会更新指针里的值
	for _, word := range updatedWords {
		newKeywords[word.Word] = word.Clone()
	}

	// 指针指向新的 keywords 集合
	*words = newKeywords
}
