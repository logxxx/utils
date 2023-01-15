package objutil

func RemoveDuplicate(input []string) []string {

	if len(input) == 0 {
		return input
	}

	uniq := map[string]bool{}
	output := make([]string, 0)
	for _, elem := range input {
		if _, ok := uniq[elem]; ok {
			continue
		}
		uniq[elem] = true
		output = append(output, elem)
	}

	return output

}
