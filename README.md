# Nagios

## Usage

```go
package main

import (
  "math"
  "github.com/lscheidler/go-nagios"
)

func main() {
  nagios := nagios.Init()
  defer nagios.Exit()

  // do not add perfdata to message, default is true
  nagios.ShowPerfdata = false

  var (
    value float64 = 5
    warning float64 = 4
    critical float64 = math.NaN()

    valueP float64 = 20
    maxP float64 = 100
    warningP float64 = 40
    criticalP float64 = 80
  )

  // check value against warning threshold and ignore critical, because it is NaN
  nagios.CheckThreshold("example", value, warning, critical)

  // check percentage calculated from valueP and maxP against warningP and criticalP thresholds
  nagios.CheckPercentageThreshold("percentage.example", valueP, maxP, warningP, criticalP)
}
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/lscheidler/go-nagios.

## License

The code is available as open source under the terms of the [Apache 2.0 License](http://opensource.org/licenses/Apache-2.0).
