package data

import (
	"fmt"
	"strconv"
	"time"
)

type CreatedAt struct {
	OurTime time.Time
}

type Latency float64

func(latency Latency) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%.2f ms", latency)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}

func(createdAt CreatedAt) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%q", createdAt.OurTime.Format(time.ANSIC))

	return []byte(jsonValue), nil
}
