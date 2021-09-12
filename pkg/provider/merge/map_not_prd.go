// 合并功能提供商

//go:build !prod

package merge

func (p *Map) Push(data map[string]interface{}) error { // 合并功能提供商
	for _, val := range data {
		if _, ok := val.(Handle); !ok {
			panic("any element must be Handle")
		}
	}

	return p.MergeProvider.Push(data)
}
