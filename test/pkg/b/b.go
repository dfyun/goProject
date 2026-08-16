package b

import (
	"fmt"

	"example.com/testpkg/a"
)

func init() {
	fmt.Println("init b", a.A())
}
