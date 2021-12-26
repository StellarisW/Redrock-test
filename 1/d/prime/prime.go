/*
分段埃氏筛法
空间复杂度: O(√ n)
时间复杂度: O(nloglogn)
思路:
既然GO写并发这么容易，那就用分段筛了！！！
直接把所有CPU核心全部拉满
根据唯一分解定理,对于不超过n的素数p,
删除2p,3p,4p...,当处理完所有的数后,
还没有被删除的数就时素数

参考资料:

https://blog.csdn.net/gemenhao/article/details/6268933?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522164050295816780274126656%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=164050295816780274126656&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~baidu_landing_v2~default-1-6268933.pc_search_em_sort&utm_term=%E5%88%86%E6%AE%B5%E7%AD%9B&spm=1018.2226.3001.4187
https://blog.csdn.net/holly_Z_P_F/article/details/85063174?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522164050320816780366559802%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=164050320816780366559802&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduend~default-1-85063174.pc_search_em_sort&utm_term=%E5%9F%83%E6%B0%8F%E7%AD%9B&spm=1018.2226.3001.4187
https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes

*/

package prime

import (
	"math"
	"runtime"
	"sync"
)

func fill(nums []bool, i uint64, max uint64) {
	a := 3 * i
	for a <= max {
		nums[a/2] = true
		a = a + 2*i
	}
}

//goFill
func goFill(nums []bool, i uint64, max uint64, next chan bool) {
	fill(nums, i, max)
	<-next
}

// SieveOfEratosthenes 使用Eratosthenes的Sieve返回一个等于或低于max的所有素数的片段。
// 这是没有分段的。
func SieveOfEratosthenes(n uint64) []uint64 {
	cores := runtime.NumCPU()      //调用CPU的一个核心
	next := make(chan bool, cores) //选取下一个核心

	var vis = make([]bool, n/2+1)
	m := uint64(math.Sqrt(float64(n))) //求平方根,一个数的因数不可能大于它的平方根

	//埃氏筛
	for i := uint64(3); i <= m; i = i + 2 {
		if vis[i/2] == false {
			go goFill(vis, i, n, next)
			next <- true
		}
	}

	for i := 0; i < cores; i++ {
		next <- true
	}

	//vis数组中没有标记过的下标就是素数
	var ps []uint64
	if n >= 2 {
		ps = append(ps, 2)
	}
	for i := uint64(3); i <= n; i = i + 2 {
		if vis[i/2] == false {
			ps = append(ps, i)
		}
	}
	return ps
}

// CgPool 保存和复用临时对象，以减少内存分配，降低CG压力
var CgPool sync.Pool

// fillSegments 分段标记素数
// basePrimes 第一段素数 ; allPrimes 存储所有素数的数组 ; segSize 分段大小 ； SegNum 分段序号 ; next,nextTurn 选取下个核心
func fillSegments(n uint64, basePrimes []uint64, allPrimes *[]uint64, segSize uint64, segNum uint64, next chan bool, nextTurn []chan bool) {
	//从池中取出标记数组
	cg := (CgPool.Get()).([]bool)
	//初始化
	for i := uint64(0); i < segSize; i++ {
		cg[i] = false
	}
	for i := 0; i < len(basePrimes); i++ {
		jMax := segSize * (segNum + 1) / basePrimes[i] //分段最高值
		//在分段内标记
		for j := (segSize * segNum) / basePrimes[i]; j < jMax; j++ {
			sn := (j + 1) * basePrimes[i]
			cg[sn-segSize*segNum-1] = true
		}
	}

	//防止最后排序，增加效率
	if segNum > 1 {
		<-nextTurn[segNum]
	}

	//将分段内的素数添加至allPrimes数组
	for i := uint64(0); i < segSize; i++ {
		if !cg[i] && segSize*segNum+i+1 <= n {
			*allPrimes = append(*allPrimes, segSize*segNum+i+1)
		}
	}

	//调用下个核心
	<-next
	if int(segNum)+1 < len(nextTurn) {
		nextTurn[segNum+1] <- true
	}

	//将标记数组放入池中
	CgPool.Put(cg)
}

// Primes 您可以在中了解更多信息https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes.
// 使用分段筛。这种改进方法将大大减少埃氏筛法的内存使用。
// 除了素数片的内存分配外，操作只需要O（sqrt（n））个额外内存
func Primes(n uint64) (allPrimes []uint64) {
	if uint64(math.Log(float64(n))-1) == 0 {
		return SieveOfEratosthenes(n)
	}

	//数学中有一个函数pi（x），它将返回n以下的近似素数。
	allPrimes = make([]uint64, 0, n/uint64(math.Log(float64(n))-1))
	segSize := uint64(math.Sqrt(float64(n))) //分段大小

	//将标记数组放入池中
	CgPool.New = func() interface{} {
		return make([]bool, segSize)
	}

	//使用埃氏筛法找出第一段（即最低段）中的素数
	basePrimes := SieveOfEratosthenes(segSize) //第一段素数
	allPrimes = append(allPrimes, basePrimes...)

	cores := runtime.NumCPU()      //调用CPU当前核心
	next := make(chan bool, cores) //选取下一个核心

	//依次调用核心
	var nextTurn []chan bool
	nextTurn = make([]chan bool, n/segSize+1)
	for i := uint64(0); i < n/segSize+1; i++ {
		nextTurn[i] = make(chan bool)
	}

	for segNum := uint64(1); segNum <= n/segSize; segNum++ {
		go fillSegments(n, basePrimes, &allPrimes, segSize, segNum, next, nextTurn)
		next <- true
	}
	for i := 0; i < cores; i++ {
		next <- true
	}
	return allPrimes
}
