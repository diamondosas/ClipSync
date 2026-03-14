package main
import(
	"fmt"
)

func main(){
		ch := make (chan int)

	for i:=0; i<100; i++{
	}

	fmt.Println(<-ch)
}