# each

### Iterate
```go
package main

import (
	"fmt"
	each "github.com/evgeny-klyopov/each"
)

func main() {
	var errs *[]error
	rows := make([]string, 0)

	for i := 0; i < 100; i++ {
		rows = append(rows, fmt.Sprintf("name%d", i))
	}

	errs = each.Iterate(rows, 2, func(lines []string, hasError bool) error {
		if hasError == true {
			return nil
		}

		println("write lines", lines)

		return nil
	})

	fmt.Println(errs)
}
```

### CustomIterate
```go
package main

import (
	"fmt"
	each "github.com/evgeny-klyopov/each"
)

func main() {
	rows := make([]string, 0)

	for i := 0; i < 100; i++ {
		rows = append(rows, fmt.Sprintf("name%d", i))
	}

	e := each.NewEach(2, func(lines []string, hasError bool) error {
		if hasError == true {
			return nil
		}

		println("write lines", lines)

		return nil
	})

	for _, val := range rows {
		hasError := e.Add(val)
		if true == hasError {
			break
		}
	}
	e.Close()

	fmt.Println(e.GetErrors())
}
```