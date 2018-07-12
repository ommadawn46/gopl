package lenconv

func MToFt(m Meter) Feet { return Feet(m / 3.28084) }
func FtToM(ft Feet) Meter { return Meter(ft * 3.28084) }
