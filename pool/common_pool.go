package pool

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Common_pool struct {
	Dial func() (interface{}, error)

	TestOnBorrow func(c interface{}, t time.Time) error

	MaxIdle   int
	MaxActive int

	IdleTimeout time.Duration

	Wait bool

	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	active int

	// Stack of idleConn with most recently used at the front.
	idle list.List
}

var nowFunc = time.Now

type idleConn struct {
	p *Common_pool
	c *clientv3.Client
	t time.Time
}

type pooledConnection struct {
	p     *Common_pool
	c     *clientv3.Client
	state int
}

func (p *Common_pool) Get() interface{} {
	return nil
}

func (p *Common_pool) get() (interface{}, error) {
	p.mu.Lock()

	if timeout := p.IdleTimeout; timeout > 0 {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Back()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			if ic.t.Add(timeout).After(nowFunc()) {
				break
			}
			p.idle.Remove(e)
			p.release()
			p.mu.Unlock()
			ic.c.Close()
			p.mu.Lock()
		}
	}

	for {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Front()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			p.idle.Remove(e)
			test := p.TestOnBorrow
			p.mu.Unlock()
			if test == nil || test(ic.c, ic.t) == nil {
				return ic.c, nil
			}
			ic.c.Close()
			p.mu.Lock()
			p.release()
		}
		if p.closed {
			p.mu.Unlock()
			return nil, errors.New("error pool")
		}
		if p.MaxActive == 0 || p.active < p.MaxActive {
			dial := p.Dial
			p.active += 1
			p.mu.Unlock()
			c, err := dial()
			if err != nil {
				p.mu.Lock()
				p.release()
				p.mu.Unlock()
				c = nil
			}
			return c, err
		}
		if !p.Wait {
			p.mu.Unlock()
			return nil, errors.New("pool errors")
		}

		if p.cond == nil {
			p.cond = sync.NewCond(&p.mu)
		}
		p.cond.Wait()
	}
}

func (p *Common_pool) Close() error {
	p.mu.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.active -= idle.Len()
	if p.cond != nil {
		p.cond.Broadcast()
	}
	p.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(idleConn).c.Close()
	}
	return nil
}

func (p *Common_pool) release() {
	p.active -= 1
	if p.cond != nil {
		p.cond.Signal()
	}
}

func (ic *idleConn) Close() error {
	// c := ic.c
	return nil
}
