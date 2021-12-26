package main

import "fmt"

type room struct {
	is  bool
	num int
}

func main() {
	var n, m, pos, ans int
	ans = 0
	fmt.Scanln(&n, &m)
	all := make([][]room, 0)
	for i := 0; i < n; i++ {
		all = append(all, []room{})
		for j := 0; j < m; j++ {
			var a room
			fmt.Scanf("%t %d\n", &a.is, &a.num)
			all[i] = append(all[i], a)
		}
	}
	fmt.Scanln(&pos)
	for i := 0; i < n; i++ {
		ans += all[i][pos].num
		if all[i][pos].is == true {
			continue
		}
		l := all[i][pos].num
		pos++
		for j := 1; j <= l; {
			for {
				if pos == m {
					pos = 0
				}
				if all[i][pos].is == true {
					j++
					pos++
					break
				}
				pos++
			}
		}
		pos--
	}
	fmt.Println(ans)
}
