package assert

import (
	"log"
)

func NotNil(path string, value any) {
	if value == nil {
		log.Fatalf("Assertion failed(%s): value is nil!", path)
	}
}

func SafeWrite(_ int, err error) {
	if err != nil {
		log.Fatal(err)
	}
}
