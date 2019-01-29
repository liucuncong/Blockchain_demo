package main

func main()  {
	bc := NewBlockChain()
	cli := CLI{bc}
	cli.Run()

	/*
	bc.AddBlock("1111111")
	bc.AddBlock("2222222")

	*/
}
