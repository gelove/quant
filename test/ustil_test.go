package test

import (
	"fmt"
	"os"
	"quant/cmd"
	"quant/internal/app/dto"
	"quant/pkg/utils/date"
	"quant/pkg/utils/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/shopspring/decimal"
)

func init() {
	cmd.InitConfig("/Users/allen/Projects/Go/mod/quant/config/config.yml")
}

func BenchmarkMutex(b *testing.B) {
	var number int
	var lock sync.Mutex
	// for i := 0; i < b.N; i++ {
	// 	go func() {
	// 		defer lock.Unlock()
	// 		lock.Lock()
	// 		number++
	// 	}()
	// }
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			func() {
				defer lock.Unlock()
				lock.Lock()
				number++
			}()
		}
	})
	b.Logf("BenchmarkMutex number: %d, %d", number, b.N)
}

func BenchmarkAtomic(b *testing.B) {
	var number int32
	// for i := 0; i < b.N; i++ {
	// 	go func() {
	// 		atomic.AddInt32(&number, 1)
	// 	}()
	// }
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt32(&number, 1)
		}
	})
	b.Logf("BenchmarkAtomic number: %d, %d", number, b.N)
}

func TestPoint(t *testing.T) {
	var a, b *dto.Frame
	t.Logf("TestPoint a => %p, &a => %v, b => %p, &b => %v, a == b is %v", a, &a, b, &b, a == b)
	var unsafepL = (*unsafe.Pointer)(unsafe.Pointer(&a))
	px1 := (*dto.Frame)(atomic.LoadPointer(unsafepL))
	t.Logf("TestPoint px1 => %p %p", px1, &px1)
}

func TestTime(t *testing.T) {
	now := time.Now()
	nowLocal := now.Local()
	t.Logf("TestDecimal => now: %v nowLocal: %v", now, nowLocal)
	startTime := "2021-06-19T00:00:00"
	date.GetMilliUnix(startTime, date.YMD_HIS)

	var transactTime int64 = 1600000000595
	res := time.Unix(transactTime/1e3, transactTime%1e3*int64(time.Millisecond))
	t.Logf("TestDecimal res => %v", res)
	t.Logf("TestDecimal 1e6 == 1000000 => %v", 1e6 == 1000000)
}

func TestDecimal(t *testing.T) {
	step := decimal.NewFromFloat(0.001)
	res := decimal.NewFromFloat(20.2535).Div(step).Floor().Mul(step)
	t.Logf("TestDecimal res => %+v", res)
}

func TestJson(t *testing.T) {
	v := `{"stream":"btcusdt@depth10","data":{"lastUpdateId":1228645,"bids":[["41275.26000000","0.00240000"],["41266.04000000","0.00240000"]],"asks":[["41453.97000000","0.00241100"],["41455.01000000","0.01206200"]]}}`
	res := &dto.Frame{}
	json.MustDecode([]byte(v), res)
	t.Logf("TestJson res => %#v", res)
	depth := &dto.Depth{}
	json.MustTransform(res.Data, depth)
	t.Logf("TestJson depth => %#v", depth)
}

func TestPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Logf("TestGetAccount err => %+v", err)
		return
	}
	res := fmt.Sprintf("%s%c%s", home, os.PathSeparator, "quant")
	t.Logf("TestPath res => %#v", res)
}
