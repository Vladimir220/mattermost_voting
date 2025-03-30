package tarantool

func (t TarantoolDAO) Close() {
	t.conn.CloseGraceful()
}
