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
	Port int `long:"port" env:"FA_SHOP_PORT" default:"8081" description:"port"`

	SteamShopAddress string `long:"steam-shop-address" env:"FA_STEAM_SHOP_ADDRESS" default:"https://steamcommunity.com" description:"Steam API domain"`

	RedisConnect string `long:"redis-connection" env:"FA_REDIS_CACHE_CONNECT" default:":6379" description:"Redis connection string"`
	RedisTTL     int    `long:"redis-ttl" env:"FA_REDIS_TTL" default:"360" description:"Redis cache time to live in seconds"`
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

	shopAPI := shop.NewSteamAPI(opts.SteamShopAddress)
	cache := storage.NewRedisCache(opts.RedisConnect, opts.RedisTTL)

	srv := rest.Rest{
		ShopService: shopAPI,
		Cache:       cache,
	}

	srv.Run(opts.Port)
}
