package lenconv

import "fmt"

type Meter float64
type Feet float64

const (
	MilliMeter Meter = 0.001
	CentiMeter Meter = 0.01
	KiloMeter  Meter = 1000
)

func (m Meter) String() string { return fmt.Sprintf("%gm", m) }
func (ft Feet) String() string { return fmt.Sprintf("%gft", ft) }
