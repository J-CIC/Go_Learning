package add

func Add(list []int) int {
	sum := list[0]
	for i, _ := range list {
		sum += list[i]
	}
	return sum
}

func main() {
    list := []int{1,2,3}
    Add(list)
}