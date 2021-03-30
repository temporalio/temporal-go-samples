package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/temporalio/samples-go/cryptconverter"
	"go.temporal.io/sdk/converter"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TEMPORAL_PLUGIN_DATA_CONVERTER",
	MagicCookieValue: "cookie",
}

func main() {
	var pluginMap = map[string]plugin.Plugin{
		"DataConverter": &converter.DataConverterPlugin{
			Impl: cryptconverter.NewCryptDataConverter(
				converter.GetDefaultDataConverter(),
			),
		},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}