package main

func main()  {
	bc := NewBlockChain("1QFbWnp7GtxrTkNqgtmhj65oyk1i2T3sWm")
	cli := CLI{bc}
	cli.Run()

	/*
	bc.AddBlock("1111111")
	bc.AddBlock("2222222")

	*/
}
