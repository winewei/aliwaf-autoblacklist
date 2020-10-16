package utils

// 去重 https://blog.csdn.net/qq_27068845/article/details/77407358
func RemoveRepByLoop(slc []string) []string {
	result := []string{}
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

func RemoveRep(slc []string) []string{
	if len(slc) < 1024 {
		return RemoveRepByLoop(slc)
	}else {
		return RemoveRepByMap(slc)
	}
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc{
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l{
			result = append(result, e)
		}
	}
	return result
}