package steps

import "sync"

func Executor(steps []func() bool) bool {
	var (
		wg  = sync.WaitGroup{}
		out = make(chan bool)
	)

	for i := range steps {
		wg.Add(1)

		go func() {
			defer wg.Done()
			v := steps[i]()
			out <- v
		}()
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	res := true
	for v := range out {
		res = res && v
	}

	return res
}
