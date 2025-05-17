package ticker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/woshilaixuex/csd_chat_backend/app/util/ticker"
)

func TestTimeWheel(t *testing.T) {
	tw := ticker.NewTimeWheel(1*time.Second, 10)
	tw.Start()

	go func() {
		for i := 0; i < 10; i++ {
			tw.AddTaskWithKey(fmt.Sprintf("print:100.%v", i), func() {
				fmt.Println(100)
			},
				2*time.Second, true)
		}
	}()

	time.Sleep(25 * time.Second)
}
