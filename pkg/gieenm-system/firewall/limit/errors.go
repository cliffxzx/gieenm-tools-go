package limit

import "fmt"

type LimitExceedError struct {
	MaxCount int
	NowCount int
}

func (e LimitExceedError) Error() string {
	return fmt.Sprintf("MaxCount is %d, but now count is %d", e.MaxCount, e.NowCount)
}
