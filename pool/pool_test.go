package pool

import (
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func init() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestWorker(t *testing.T) {
	pool := make(chan *worker)
	wk := newWorker(pool)
	wk.start()
	if wk == nil {
		t.Fatal("worker nil error")
	}

	wk = <-pool
	if wk == nil {
		t.Fatal("worker should register itself.")
	}

	called := false
	done := make(chan bool)

	job := func() {
		called = true
		done <- true
	}

	wk.jobChannel <- job
	<-done
	if !called {
		t.Fatal("job should be called, but no")
	}
}

func TestNewPool(t *testing.T) {
	pool := NewPool(1000, 100000)
	defer pool.Release()

	num := 10000
	var wg sync.WaitGroup
	wg.Add(num)

	var completed uint64
	for i := 0; i < num; i++ {
		arg := uint64(1)
		pool.JobQueue <- func() {
			defer wg.Done()
			atomic.AddUint64(&completed, arg)
		}
	}

	wg.Wait()
	if completed != uint64(num) {
		t.Fatal("completed count error")
	}

}

func BenchmarkPool(b *testing.B) {
	pool := NewPool(1, 10)
	defer pool.Release()

	log.SetOutput(ioutil.Discard)

	for n := 0; n < b.N; n++ {
		pool.JobQueue <- func() {
			log.Printf("I am worker! Number %d\n", n)
		}
	}
}
