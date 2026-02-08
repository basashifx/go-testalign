package basic

import "testing"

func TestService_Delete(t *testing.T) {}

func TestService_Create(t *testing.T) {} // want `TestService_Create corresponds to Service\.Create \(service\.go:\d+\) but appears before TestService_Delete which corresponds to Service\.Delete \(service\.go:\d+\)`

func TestService_Read(t *testing.T) {} // want `TestService_Read corresponds to Service\.Read \(service\.go:\d+\) but appears before TestService_Delete which corresponds to Service\.Delete \(service\.go:\d+\)`
