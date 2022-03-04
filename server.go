package main

var (
	makeChannel       = make(chan *register)
	joinChannel       = make(chan *register)
	unregisterChannel = make(chan *unregister)
	// ブロックされたくないので 100 に設定
	forwardChannel = make(chan forward, 100)
)

// roomId がキーになる
type room struct {
	connections map[string]*connection
}

func server() {
	// room を管理するマップはここに用意する
	var m = make(map[string]room)
	// ここはシングルなのでロックは不要
	for {
		select {
		case makeRoom := <-makeChannel:
			c := makeRoom.connection
			rch := makeRoom.resultChannel
			_, ok := m[c.roomID]
			if ok {
				// room があった
				rch <- exists
			} else {
				// room がなかった
				var connections = make(map[string]*connection)
				connections[c.ID] = c
				// room を追加
				m[c.roomID] = room{
					connections: connections,
				}
				c.debugLog().Msg("CREATED-ROOM")
				rch <- create
			}
		case joinRoom := <-joinChannel:
			c := joinRoom.connection
			rch := joinRoom.resultChannel
			r, ok := m[c.roomID]
			if ok {
				// room があった
				if len(r.connections) == 1 {
					r.connections[c.ID] = c
					m[c.roomID] = r
					rch <- join
				} else {
					// room あったけど満杯
					rch <- full
				}
			} else {
				// room がなかった
				rch <- none
			}
		case unregister := <-unregisterChannel:
			c := unregister.connection
			// room を探す
			r, ok := m[c.roomID]
			// room がない場合は何もしない
			if ok {
				_, ok := r.connections[c.ID]
				if ok {
					for _, connection := range r.connections {
						// 両方の forwardChannel を閉じる
						close(connection.forwardChannel)
						connection.debugLog().Msg("CLOSED-FORWARD-CHANNEL")
						connection.debugLog().Msg("REMOVED-CLIENT")
					}
					// room を削除
					delete(m, c.roomID)
					c.debugLog().Msg("DELETED-ROOM")
				}
			}
		case forward := <-forwardChannel:
			r, ok := m[forward.connection.roomID]
			// room がない場合は何もしない
			if ok {
				// room があった
				for connectionID, client := range r.connections {
					// 自分ではない方に投げつける
					if connectionID != forward.connection.ID {
						client.forwardChannel <- forward
					}
				}
			}
		}
	}
}
