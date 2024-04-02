package errgroup

import (
	"context"
	"sync"
	"time"
)

type Group struct {
	cancel  func(error)
	wg      sync.WaitGroup
	ticker  *time.Ticker
	outchan chan any
	err     error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel, outchan: make(chan any, 1)}, ctx
}

func (g *Group) Wait() error {
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}

func (g *Group) Go(f func() (any, error)) {
	if g.ticker != nil {
		<-g.ticker.C
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		result, err := f()
		if err != nil && g.err == nil {
			g.err = err
			if g.cancel != nil {
				g.cancel(g.err)
			}
		}
		g.outchan <- result
	}()
}

func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.ticker = nil
		return
	}
	interval := time.Second / time.Duration(n)
	g.ticker = time.NewTicker(interval)
	g.outchan = make(chan any, n)
}

func (g *Group) GetOutChan() <-chan any {
	go func() {
		g.wg.Wait()
		close(g.outchan)
	}()
	return g.outchan
}
