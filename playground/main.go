package main

import (
	"fmt"
	"os"

	"github.com/UpCraft-Solutions-Pvt-Ltd/upcraft-agent/core/plugins/spotify"
)

// Playground harness for testing one plugin in isolation.
func main() {
	if os.Getenv("SPOTIFY_ACCESS_TOKEN") == "" {
		fmt.Println("Please set SPOTIFY_ACCESS_TOKEN env var")
		return
	}

	plugin, err := spotify.NewMusicPluginFromEnv()
	if err != nil {
		fmt.Printf("failed to initialize plugin: %v\n", err)
		return
	}

	if err := plugin.Play("Hymn for the Weekend"); err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}

	fmt.Println("Success: Music play request sent.")
}
