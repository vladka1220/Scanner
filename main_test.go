package main

import (
	"flag"
	"os"
	"testing"

	"basis_go/types"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestRunModeSpot(t *testing.T) {
	runOnce = true
	resetFlags()
	os.Args = []string{"cmd", "-mode=spot"}
	called := []string{}

	collectSpotPrices = func(fetchers ...map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
		called = append(called, "collectSpot")
		return types.TokenPrices{}
	}
	compareSpotPrices = func(types.TokenPrices) { called = append(called, "compareSpot") }

	main()

	if len(called) != 2 || called[0] != "collectSpot" || called[1] != "compareSpot" {
		t.Fatalf("spot mode not executed correctly: %v", called)
	}
}

func TestRunModeFutures(t *testing.T) {
	runOnce = true
	resetFlags()
	os.Args = []string{"cmd", "-mode=futures"}
	called := []string{}

	collectFuturesPrices = func(fetchers ...map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
		called = append(called, "collectFutures")
		return types.TokenPrices{}
	}
	compareFuturesPrices = func(types.TokenPrices) { called = append(called, "compareFutures") }

	main()

	if len(called) != 2 || called[0] != "collectFutures" || called[1] != "compareFutures" {
		t.Fatalf("futures mode not executed correctly: %v", called)
	}
}

func TestRunModeSpotFutures(t *testing.T) {
	runOnce = true
	resetFlags()
	os.Args = []string{"cmd", "-mode=spotfutures"}
	called := []string{}

	collectSpotPrices = func(fetchers ...map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
		called = append(called, "collectSpot")
		return types.TokenPrices{}
	}
	collectFuturesPrices = func(fetchers ...map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
		called = append(called, "collectFutures")
		return types.TokenPrices{}
	}
	compareSpotFutures = func(types.TokenPrices, types.TokenPrices) { called = append(called, "compareSpotFutures") }

	main()

	expected := []string{"collectSpot", "collectFutures", "compareSpotFutures"}
	for i, v := range expected {
		if i >= len(called) || called[i] != v {
			t.Fatalf("spotfutures mode not executed correctly: %v", called)
		}
	}
}

func TestRunModePump(t *testing.T) {
	runOnce = true
	resetFlags()
	os.Args = []string{"cmd", "-mode=pump"}
	called := false
	monitorPumps = func() { called = true }

	main()

	if !called {
		t.Fatalf("pump mode not executed")
	}
}

func TestRunModeDefault(t *testing.T) {
	runOnce = true
	resetFlags()
	os.Args = []string{"cmd", "-mode=unknown"}
	called := false
	collectSpotPrices = func(fetchers ...map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
		called = true
		return nil
	}

	main()

	if called {
		t.Fatalf("default branch should not call any collect functions")
	}
}
