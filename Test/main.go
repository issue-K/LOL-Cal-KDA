package main

import "fmt"

func main(){
	var s []int = make( []int,0 )
	for i :=0;i<5;i++{
		s = append( s,i )
	}
	for _,i := range s {
		fmt.Println( i )
	}
}