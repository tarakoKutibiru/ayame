package main

const (
	// 登録成功
	// 部屋が作成された相手を待つ
	create int = iota
	// すでに部屋はあり相手が待ってる
	// offer を出させる
	// 登録成功
	join
	// 満員だったので Reject か Error を返す
	// 登録失敗
	full

	// 部屋がなかった
	none

	// 部屋が存在した
	exists
)

type register struct {
	connection    *connection
	resultChannel chan int
}

// rawMessage には JOSN パース前の offer / answer / candidate が入る
type forward struct {
	connection *connection
	rawMessage []byte
}

type unregister struct {
	connection *connection
}
