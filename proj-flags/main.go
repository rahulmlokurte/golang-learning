package main

import (
  "flag"
  "fmt"
  "os"
)

func main() {

  println("Starting...")
  flag.Usage = func(){
    f := flag.CommandLine.Output()
    fmt.Fprintln(f, "lf - Terminal file manager")
		fmt.Fprintln(f, "")
		fmt.Fprintf(f, "Usage:  %s [options] [cd-or-select-path]\n\n", os.Args[0])
		fmt.Fprintln(f, "  cd-or-select-path")
		fmt.Fprintln(f, "        set the initial dir or file selection to the given argument")
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "Options:")
    flag.PrintDefaults()
  }
  showDoc := flag.Bool("doc", false, "show documentation")
  
  flag.Parse()
  fmt.Println("end")
  switch {
  case *showDoc:
    fmt.Println("docume")
}
}
