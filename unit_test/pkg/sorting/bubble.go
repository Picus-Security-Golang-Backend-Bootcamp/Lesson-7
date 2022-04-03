package sorting

func BubbleSort(elements []int) {
	isContinue := true
	for isContinue {
		isContinue = false
		for i := 0; i < len(elements)-1; i++ {
			if elements[i] < elements[i+1] {
				isContinue = true
				elements[i], elements[i+1] = elements[i+1], elements[i]
			}
		}
	}
}
