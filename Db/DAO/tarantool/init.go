package tarantool

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func (t *TarantoolDAO) init(address, user, password string) (err error) {
	dialer := tarantool.NetDialer{
		Address:  address,
		User:     user,
		Password: password,
	}
	opts := tarantool.Opts{
		Timeout: 2 * time.Second,
	}

	t.conn, err = tarantool.Connect(context.Background(), dialer, opts)
	if err != nil {
		return fmt.Errorf("ошибка подключания к Tarantool: %v", err)
	}

	return
}
