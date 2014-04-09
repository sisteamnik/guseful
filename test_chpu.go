package main

import (
  "fmt"
  "os"
  "bufio"
  "purego/chpu/chpu"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    for {
        line, err := reader.ReadString('\n')

        if err != nil {
            // You may check here if err == io.EOF

            break
        }
        fmt.Println(chpu.Chpu(string(line)))
    }    
}