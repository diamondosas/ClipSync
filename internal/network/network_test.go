package network_test

import (
	"clipsync/internal/network"
	"testing"
	"context"
)

func TestRegisterandBrowse(t *testing.T){
	ctx := context.WithoutCancel(context.Background())
	go network.RegisterDevice(ctx)

}