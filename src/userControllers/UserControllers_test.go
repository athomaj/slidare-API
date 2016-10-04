package userControllers

import "testing"

func TestFirstTest(t *testing.T) {
  a := 10
  b := 9
  if (a != b) {
    t.Error("values are different", a, b)
  }
}
