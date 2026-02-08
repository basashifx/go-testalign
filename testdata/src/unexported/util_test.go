package unexported

import "testing"

func Test_normalize(t *testing.T) {}

func Test_validate(t *testing.T) {} // want `Test_validate corresponds to validate \(util\.go:\d+\) but appears before Test_normalize which corresponds to normalize \(util\.go:\d+\)`

func Test_format(t *testing.T) {}
