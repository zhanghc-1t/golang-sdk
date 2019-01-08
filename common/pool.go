package common

import (
	"errors"
	"sync"
	"time"
)

type ConnectResource interface {
	Close() error
}

type ConnectFactory func() (ConnectResource, error)

type Connect struct {
	conn ConnectResource
	time time.Time
}

type Pool struct {
	mu       sync.Mutex
	connects chan *Connect
	factory  ConnectFactory
	closed   bool
	timeout  time.Duration
}

func (pool *Pool) Close() {
	if pool.closed {
		return
	}

	pool.mu.Lock()
	pool.closed = true
	close(pool.connects)

	for connect := range pool.connects {
		connect.conn.Close()
	}

	pool.mu.Unlock()
}

func (pool *Pool) Length() int {
	return len(pool.connects)
}

func (pool *Pool) Put(conn ConnectResource) error {
	if pool.closed {
		return errors.New("连接池已关闭")
	}

	select {
	case pool.connects <- &Connect{
		conn: conn,
		time: time.Now(),
	}:
		return nil
	default:
		conn.Close()
		return errors.New("连接池已满")
	}
}

func (pool *Pool) Get() (ConnectResource, error) {
	if pool.closed {
		return nil, errors.New("连接池已关闭")
	}

	for {
		select {
		case connect, ok := <-pool.connects:
			if !ok {
				return nil, errors.New("连接池关闭")
			}
			if time.Now().Sub(connect.time) > pool.timeout {
				connect.conn.Close()
				continue
			}

			return connect.conn, nil
		default:
			conn, err := pool.factory()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}

	}
}

func NewPool(factory ConnectFactory, cap int, timeout time.Duration) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("连接池大小不能小于1")
	}
	if timeout <= 0 {
		return nil, errors.New("连接超时时长不能小于1")
	}

	pool := &Pool{
		mu:       sync.Mutex{},
		connects: make(chan *Connect, cap),
		factory:  factory,
		closed:   false,
		timeout:  timeout,
	}

	for i := 0; i < cap; i++ {
		conn, err := factory()
		if err != nil {
			pool.Close()
			return nil, errors.New("连接池填充出错")
		}
		pool.connects <- &Connect{
			conn: conn,
			time: time.Now(),
		}
	}

	return pool, nil
}
