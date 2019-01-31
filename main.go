package main

func main()  {
	bc := NewBlockChain("张三")
	cli := CLI{bc}
	cli.Run()

	/*
	bc.AddBlock("1111111")
	bc.AddBlock("2222222")

	*/
}
