package steps

import "sync"

// Screen runs multiple screening steps concurrently and returns true only if
// all steps return true. This allows combining multiple technical analysis
// conditions that must all be satisfied for a trading signal.
func Screen(screeners []func() bool) bool {
	var (
		wg  = sync.WaitGroup{}
		out = make(chan bool)
	)

	for i := range screeners {
		wg.Add(1)

		go func() {
			defer wg.Done()
			v := screeners[i]()
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
