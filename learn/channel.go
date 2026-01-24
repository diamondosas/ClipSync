package learn

import(
	"fmt"
)

func main(){
		ch := make (chan int)

	for i:=0; i<100; i++{
		ch <-i
		
	}

	fmt.Println(<-ch)
}