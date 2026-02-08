package pkgfuncs

import "testing"

func TestApplyConfig(t *testing.T) {}

func TestParseConfig(t *testing.T) {} // want `TestParseConfig corresponds to ParseConfig \(funcs\.go:\d+\) but appears before TestApplyConfig which corresponds to ApplyConfig \(funcs\.go:\d+\)`

func TestValidateConfig(t *testing.T) {} // want `TestValidateConfig corresponds to ValidateConfig \(funcs\.go:\d+\) but appears before TestApplyConfig which corresponds to ApplyConfig \(funcs\.go:\d+\)`
