package errgroup

import (
	"context"
	"sync"
	"time"
)

// Поитогу вышло очень похоже на errgroup
type Group struct {
	cancel  func(error)
	wg      sync.WaitGroup
	ticker  *time.Ticker
	outchan chan any
	err     error
	errOnce sync.Once
}

// Функция входа
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel, outchan: make(chan any, 1)}, ctx
}

// По сути финальная функция ожидает завершение всех горутин
func (g *Group) Wait() error {
	g.wg.Wait()

	return g.err
}

// запуск горутин, rps осуществляется за счёт timeticker
func (g *Group) Go(f func() (any, error)) {
	if g.ticker != nil {
		<-g.ticker.C
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		result, err := f()
		// если горутина завершилась с ошибкой
		// то запускаем отмену контекста чтобы не ждать завершения других горутин
		if err != nil && g.cancel != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel(g.err)
			})
			return

		}
		// если всё ок то кладём результат в канал
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
	g.outchan = make(chan any, n) // пересоздаём канал согласно кол-ву горутин
}

// Возвращаем генерик канал, заодно запускам горутинку ожидания конца, чтобы закрыть канал
// чтобы не блокировать range
func (g *Group) GetOutChan() <-chan any {
	go func() {
		g.wg.Wait()
		g.ticker.Stop()
		close(g.outchan)
	}()
	return g.outchan
}
