package uuid_generator

import (
	"crypto/md5"
	"github.com/gofrs/uuid"
)

func GenerateFromStr(input string) string {
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(input))
	id, err := uuid.FromBytes(md5Hasher.Sum(nil))
	if err != nil {
		panic("failed to generate uuid")
	}
	return id.String()
}

func NewUUIDV4() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id.String()
}
