package main

import "fmt"

func main() {
  dataChannel := make(chan string,3)
  dataChannel <- "Some sample Data"
  dataChannel <- "Some Other Sample Data"
  dataChannel <- "Buffered Channel"
  fmt.Println(<-dataChannel)
  fmt.Println(<-dataChannel)
  fmt.Println(<-dataChannel)
}
