package algorithmia

import (
	"encoding/json"
	"fmt"
)

func debug(args ...interface{}) {
	fmt.Println(args...)
}

func jsonString(x interface{}) string {
	y, _ := json.Marshal(x)
	return string(y)
}
