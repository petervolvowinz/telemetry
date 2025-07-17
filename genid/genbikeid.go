package genid

import "github.com/google/uuid"

type Bikeid struct {
	id string
}

type BikeIdGenerator func() string

func UuidGenerator() string {
	return uuid.NewString()
}

func ShortIdGenerator() string {
	return uuid.NewString()[:8]
}

func (bid *Bikeid) GenBikeId(generator BikeIdGenerator) string {
	uid := generator()
	return uid
}
