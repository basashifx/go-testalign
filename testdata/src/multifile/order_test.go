package multifile

import "testing"

func TestOrder_Cancel(t *testing.T) {}

func TestOrder_Place(t *testing.T) {} // want `TestOrder_Place corresponds to Order\.Place \(order\.go:\d+\) but appears before TestOrder_Cancel which corresponds to Order\.Cancel \(order\.go:\d+\)`
