package weightconv

func KgToP(kg KiloGram) Pound { return Pound(kg * 2.20462) }
func PToKg(p Pound) KiloGram { return KiloGram(p / 2.20462) }
