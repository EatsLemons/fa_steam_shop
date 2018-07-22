package main

import (
	"log"
	"os"
	"runtime"

	"github.com/EatsLemons/fa_steam_shop/rest"
	"github.com/EatsLemons/fa_steam_shop/shop"
	"github.com/EatsLemons/fa_steam_shop/storage"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `long:"port" env:"FA_STEAM_SHOP_PORT" default:"8081" description:"port"`
}

func main() {
	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		log.Println(e.Error())
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Started with:")
	log.Printf("%+v", opts)

	shopAPI := shop.NewSteamAPI("https://steamcommunity.com")
	cache := storage.NewRedisCache(":6379", 10)

	srv := rest.Rest{
		ShopService: shopAPI,
		Cache:       cache,
	}

	srv.Run(opts.Port)
}
