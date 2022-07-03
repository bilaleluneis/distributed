package distributed

type Node struct{}

func (n Node) Push(val int) error {
	return client.Call("RpcNode.Insert", val, nil)
}

func (n Node) Pop(remove bool) (int, error) {
	var result int
	err := client.Call("RpcNode.Retrieve", remove, &result)
	return result, err
}
