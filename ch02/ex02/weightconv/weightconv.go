package weightconv

import "fmt"
type KiloGram float64
type Pound float64

const (
  Gram KiloGram = 0.001
)

func (kg KiloGram) String() string { return fmt.Sprintf("%gkg", kg) }
func (p Pound) String() string { return fmt.Sprintf("%glb", p) }
