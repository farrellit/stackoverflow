package main
import(
  "log"
)


func summer(src <-chan int, result chan<- int64) {
  var sum int64
  var count int
  for i := range src {
    sum += int64(i)
    count++
  }
  log.Printf("summer: summed %d ints: %d", count, sum)
  result<-sum
}

func main() {
  var src = make(chan int)
  var dest = make(chan int64)
  go summer(src,dest)
  for i:=0; i<1000;i++{
    src<-i
  }
  close(src)
  <-dest
}
