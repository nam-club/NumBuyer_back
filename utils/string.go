package utils

// sliceに指定の要素が含まれているか判定
func ContainsString(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}

// sliceに指定の要素が含まれているか判定。
// 含まれている場合、sliceのインデックスを返却
// 含まれていない場合、-1を返却
func ContainsStringWithIndex(s []string, v string) int {
	for i, a := range s {
		if a == v {
			return i
		}
	}
	return -1
}

func DeleteSliceElement(s []string, i int) []string {
	s = append(s[:i], s[i+1:]...)
	//新しいスライスを用意する
	n := make([]string, len(s))
	copy(n, s)
	return n
}
