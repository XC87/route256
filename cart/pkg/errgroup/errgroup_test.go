package errgroup

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestErrGroup(t *testing.T) {
	t.Run("Test success case", func(t *testing.T) {
		eg, _ := WithContext(context.Background())
		eg.SetLimit(10)
		eg.Go(func() (result any, err error) {
			return 1, nil
		})
		eg.Go(func() (result any, err error) {
			return 1, nil
		})

		outPutChannel := eg.GetOutChan()
		assert.NotNil(t, outPutChannel)

		if err := eg.Wait(); err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		// проверяем что обе задачи завершились
		assert.Equal(t, 2, len(outPutChannel))
	})

	t.Run("Test failure case", func(t *testing.T) {
		eg, ctx := WithContext(context.Background())
		eg.SetLimit(10)
		var ErrorTest = errors.New("some error")
		// проверяем что есть отмена по контексту в случае ошибки
		eg.Go(func() (result any, err error) {
			select {
			case <-ctx.Done():
				return nil, nil
			case <-time.After(10 * time.Second):
				return nil, nil
			}
		})
		eg.Go(func() (result any, err error) {
			return nil, ErrorTest
		})

		outPutChannel := eg.GetOutChan()
		assert.NotNil(t, outPutChannel)

		if err := eg.Wait(); !assert.ErrorIs(t, err, ErrorTest) {
			t.Error("Expected same error")
		}
		// проверяем что обе задачи завершились
		assert.Equal(t, 2, len(outPutChannel))
	})
}
