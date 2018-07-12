package tempconv

import (
  "testing"
)

var EPSILON float64 = 0.00001

func TestCToF(t *testing.T) {
  actual := CToF(FreezingCs)
  expected := Fahrenheit(32)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestCToK(t *testing.T) {
  actual := CToK(FreezingCs)
  expected := Kelvin(273.15)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestFToC(t *testing.T) {
  actual := FToC(Fahrenheit(0))
  expected := Celsius(-17.777778)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestFToK(t *testing.T) {
  actual := FToK(Fahrenheit(0))
  expected := Kelvin(255.3722222)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestKToC(t *testing.T) {
  actual := KToC(Kelvin(0))
  expected := Celsius(-273.15)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestKToF(t *testing.T) {
  actual := KToF(Kelvin(0))
  expected := Fahrenheit(-459.67)
  if float64(actual - expected) > EPSILON || float64(expected - actual) > EPSILON {
    t.Errorf("actual %v want %v", actual, expected)
  }
}
