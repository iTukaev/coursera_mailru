package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	timer := time.Now()
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
			close(out)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}
	ExecutePipeline(hashSignJobs...)
	fmt.Println(time.Since(timer))
}

func ExecutePipeline(jobs ...job) {
	var workChan []chan interface{}
	for i := 0; i < len(jobs) + 1; i++ {
		workChan = append(workChan, make(chan interface{}))
	}
	for num, _ := range jobs {
		go jobs[num](workChan[num], workChan[num + 1])
	}
	fmt.Println(<- workChan[len(jobs)])
}

func SingleHash(in, out chan interface{}) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	wg := &sync.WaitGroup{}
	for fibNum := range in {
		<- ticker.C
		data := strconv.Itoa(fibNum.(int))
		md5 := DataSignerMd5(data)
		wg.Add(1)
		go func(data, md5 string, out chan interface{}) {
			inOne := make(chan string)
			outOne := make(chan string)
			inTwo := make(chan string)
			outTwo := make(chan string)
			go func(in, out chan string) {
				out <- DataSignerCrc32(<-in)
			}(inOne, outOne)
			go func(in, out chan string) {
				out <- DataSignerCrc32(<-in)
			}(inTwo, outTwo)
			inOne <- data
			inTwo <- md5
			out <- <-outOne + "~" + <-outTwo
			wg.Done()
		}(data, md5, out)
	}
	wg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			multi := make(map[int]string)
			mu := &sync.Mutex{}
			wgTwo := &sync.WaitGroup{}
			for i := 0; i < 6; i++ {
				wgTwo.Add(1)
				go func(i int) {
					th := strconv.Itoa(i)
					res := DataSignerCrc32(th + data.(string))
					mu.Lock()
					multi[i] = res
					mu.Unlock()
					wgTwo.Done()
				}(i)
			}
			wgTwo.Wait()
			res := ""
			for i := 0; i < 6; i++ {
				mu.Lock()
				res += multi[i]
				mu.Unlock()
			}
			out <- res
			wg.Done()
		}(data)
	}
	wg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{})  {
	var outputData []string
	res := ""
	for val := range in {
		outputData = append(outputData, val.(string)+"_")
	}
	sort.Strings(outputData)
	for _, val := range outputData {
		res += val
	}
	out <- res[:len(res) - 1]
}