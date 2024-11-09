
# warp

**warp** is an experimental Go library that introduces monadic abstractions, enabling functional programming patterns such as `Result` and `Event` monads. By leveraging Go's generics, warp facilitates cleaner error handling and the composition of asynchronous operations.

## Features

- **Result Monad**: Represents computations yielding a result or error, useful for chaining operations with error handling.
- **Event Monad**: Allows working with time-based events in a functional style.
- **Applicative and Functor Patterns**: Supports functional programming paradigms within Go’s type-safe generics.

## Installation

To install warp, run:

```sh
go get github.com/onur1/warp
```

## Usage

### Result Monad

The `Result` type encapsulates a computation that may return a value or an error, allowing you to chain dependent operations with error handling:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math"

    "github.com/onur1/warp"
    "github.com/onur1/warp/result"
)

var errNegativeSqrt = errors.New("negative square root")

func sqrt(x float64) warp.Result[float64] {
    if x < 0 {
        return result.Error[float64](errNegativeSqrt)
    }
    return result.Ok(math.Sqrt(x))
}

func main() {
    res := result.Chain(sqrt(16), sqrt)
    value, err := res(context.Background())
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", value)
    }
}
```

**Output:**

```
Result: 2
```

If you attempt to compute the square root of a negative number:

```go
res := result.Chain(sqrt(-1), sqrt)
value, err := res(context.Background())
```

**Output:**

```
Error: negative square root
```

### Event Monad

The `Event` monad models time-based events, allowing for declarative composition of asynchronous or periodic operations. Here's an example demonstrating how to use the `Event` monad to process a stream of integers, doubling each value and summing them up:

```go
package main

import (
    "context"
    "fmt"

    "github.com/onur1/warp/event"
)

func main() {
    // Create an event from a slice of integers
    numbers := event.From([]int{1, 2, 3, 4, 5})

    // Define a function to double each number
    double := func(n int) int {
        return n * 2
    }

    // Define a function to sum two numbers
    add := func(a, b int) int {
        return a + b
    }

    // Map the 'double' function over the event
    doubledNumbers := event.Map(numbers, double)

    // Fold the doubled numbers into a single sum
    sum := event.Fold(doubledNumbers, 0, add)

    // Create a channel to receive the result
    resultChan := make(chan int)

    // Start the event processing
    go sum(context.Background(), resultChan)

    // Read and print the result
    result := <-resultChan
    fmt.Printf("The sum of doubled numbers is: %d
", result)
}
```

**Output:**

```
The sum of doubled numbers is: 30
```

In this example, we create an event from a slice of integers, double each number using the `Map` function, and then sum the doubled numbers using the `Fold` function. The result is processed asynchronously and printed to the console.

## License

MIT License. See [LICENSE](LICENSE) for details.
