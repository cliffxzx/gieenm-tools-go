package test

import (
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/seeder"
)

func TestSeeder(t *testing.T) {
	seeder.Init()
	seeder.Seeder()
}
