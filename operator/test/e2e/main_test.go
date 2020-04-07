package e2e

import (
	"flag"
	"testing"

	f "github.com/operator-framework/operator-sdk/pkg/test"
)

type testArgs struct {
	skipCleanUp *bool
}

var args = &testArgs{}

func TestMain(m *testing.M) {
	args.skipCleanUp = flag.Bool("skipcleanup", false, "skip test resources clean up")

	f.MainEntry(m)
	//pflag.Parse()
}
