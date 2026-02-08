package externalapi_test

import (
	"externalapi"
	"testing"
)

var _ = externalapi.API{}

func TestAPI_Delete(t *testing.T) {}

func TestAPI_Get(t *testing.T) {} // want `TestAPI_Get corresponds to API\.Get \(api\.go:\d+\) but appears before TestAPI_Delete which corresponds to API\.Delete \(api\.go:\d+\)`

func TestAPI_Post(t *testing.T) {} // want `TestAPI_Post corresponds to API\.Post \(api\.go:\d+\) but appears before TestAPI_Delete which corresponds to API\.Delete \(api\.go:\d+\)`
