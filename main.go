package main

import (
	"fmt"
	"mcsn/src"
	"os"

	"github.com/6uf/apiGO"
	"github.com/iskaa02/qalam"
	"github.com/urfave/cli/v2"
)

func init() {
	apiGO.Clear()
	apiGO.CheckFiles()
	src.Acc.LoadState()
	src.Proxys.GetProxys()
	src.Proxys.Setup()
	qalam.Printf(src.Logo(`
•     ·   ▄▄·  ▄▄ ·    ▄ 
·██ ▐███▪▐█ ▌▪▐█ ▀  •█▌▐█
▐█ ▌▐▌▐█·██ ▄▄▄▀▀▀█▄▐█▐▐▌
██ ██▌▐█▌▐███▌▐█▄▪▐███▐█▌
▀▀  █▪▀▀▀·▀▀▀  ▀▀▀▀ ▀▀ █▪ 
LZ `) + "[:crescent_moon:]\n")
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "snipe",
				Usage: `Attempts to snipe the ign of your choice onto one of your account(s)!`,
				Action: func(c *cli.Context) error {
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Snipe(c.String("u"), c.Float64("d"), "single", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "-u NAME | this is a param value, enter ur name to target it.",
					},
					&cli.Float64Flag{
						Name:  "d",
						Usage: "Snipes a few ms earlier so you can counter ping lag.",
					},
				},
			},
			{
				Name:    "auto",
				Aliases: []string{"as", "a"},
				Usage:   "Auto sniper attempts to snipe upcoming 3 character usernames.",
				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3c")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "3l",
						Usage: "Snipe only Alphabetic names.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3l")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "3n",
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3n")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "list",
						Usage: "Snipe names from your names.txt file.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "list")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
				},
			},
			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "ping helps give you a rough estimate of your delay.",
				Action: func(c *cli.Context) error {
					time := apiGO.PingMC()
					src.PrintGrad(fmt.Sprintf("Estimated Delay (Milliseconds): %v\n", time))
					return nil
				},
			},
			{
				Name:    "proxy",
				Aliases: []string{"px", "py"},
				Usage:   "Proxy, this uses proxies to byass rate limits and snipe usernames that way, useful for MFAS",
				Action: func(c *cli.Context) error {
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Snipe(c.String("u"), c.Float64("d"), "proxy", "")
					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "-u NAME | this is a param value, enter ur name to target it.",
					},
					&cli.Float64Flag{
						Name:  "d",
						Usage: "Snipes a few ms earlier so you can counter ping lag.",
					},
				},

				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3c")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "3l",
						Usage: "Snipe only Alphabetic names.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3l")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "3n",
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3n")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
					{
						Name:  "list",
						Usage: "Snipe names from your names.txt file.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "list")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
				},
			},
			{
				Name:    "turbo",
				Aliases: []string{"t"},
				Usage:   "Turbo a name just in case it drops! - Attempts to snipe a ign every minute.",
				Action: func(c *cli.Context) error {
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Snipe(c.String("u"), float64(c.Int64("t")), "turbo", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "-u NAME | this is a param value, enter ur name to target it.",
					},
					&cli.Int64Flag{
						Name:  "t",
						Usage: "the seconds between each request",
						Value: 60,
					},
				},
			},
			{
				Name:    "mckey",
				Aliases: []string{"key", "mccode", "code"},
				Usage:   "Gets your namemc claim code (for your selected account)",
				Action: func(c *cli.Context) error {
					return src.GetNAMEMCKEY()
				},
			},
			{
				Name:    "namemc",
				Aliases: []string{"n", "nmc", "skinart"},
				Usage:   `NameMC Skin Art`,
				Action: func(c *cli.Context) error {
					src.Skinart(c.String("n"), c.String("i"))
					return nil
				},
			},
			{
				Name:  "docs",
				Usage: "Prints a list of helper apis, these will either give droptimes and other information.",
				Action: func(c *cli.Context) error {
					src.PrintGrad(`
IGN Information Apis:
 - https://buxflip.com/data/3c
 - https://buxflip.com/data/droptime/:name
 - https://buxflip.com/data/search/:name
 - https://buxflip.com/data/profile/:name
 - http://api.star.shopping/droptime/:name ("user-agent", "Sniper")

Minecraft API Documentation:
 - https://mojang-api-docs.netlify.app/
 - https://wiki.vg/Mojang_API

Other Free Snipers:
 - https://github.com/snipesharp/snipesharp     (8.5/10)
   ⚬ Snipe sharp is a ease of use utility for sniping, its coded in c++ and has the same ish features to mcsn and then some.
   ⚬ its a pretty nice application and deserves more attention.

 - https://github.com/Kqzz/MCsniperGO           (8/10)
   ⚬ MCsniperGO is made and developed by KQZZ, hes the one who inspired me to make snipers and learn programming.
   ⚬ His sniper is known and used by alot of people, its pretty good for a free sniper and has features similar to mcsn.

 - https://github.com/tropicbliss/buckshot      (8/10)
   ⚬ buckshot is a sniper coded in rust, its very fast and has vast features, the owner actively maintains it
   ⚬ and the sniper itself is updated frequently.

 - https://github.com/Everest187/Artemis-Sniper (7.5/10)
   ⚬ Artemis is a sniper made by Sylestical and Everest with the help of me in its early stages.
   ⚬ its features are pretty normal, bearer caching, sniping etc. its a very fast and precise sniper!

 - https://github.com/MCsniperPY/MCsniperPY     (4/10)
   ⚬ This sniper is also made by KQZZ, its not maintained anymore and is behind the times.
   ⚬ it isnt something i would use anymore.

Run the command "guide" to learn more about snipers and how they function.
`)
					return nil
				},
			},
			{
				Name:  "guide",
				Usage: "A small sniping guide!",
				Action: func(c *cli.Context) error {
					src.PrintGrad(`
Welcome to the sniping guide!
We will be discussing:
 - How sniping works
 - How to find a good delay
 - Which sniper to use
 - Casual or Competitive sniping
 
# How Sniping Works

Sniping is a figurative word, it means in our definition "To snipe a username when it drops".
when you go for a name that is what your doing, sniping it. Either you are sending early or later in the end the goal is 
to get the username your going for, sniping although simple can be extremely complex not only for a user of a program but
also for the devs who maintain it.

to snipe a name you will need a [200] status code, indicating success! here are some codes to expect:
 - 200 [Success]
 - 403 [Miss for MFA accounts OR Cloudflare blocked because of a Datacenter IP (Only happens when using Giftcards on a vps)]
 - 429 [Youre ratelimited!]
 - 401 [Your bearer isnt valid, and you need a new one]
 - 500 [Something is wrong with the mojang api]
 - 503 [Internal Server Error] - This error may be caused due to the API overloading with requests.
 - 404 [Not Found]             - If this error occurs, you are trying to change the username of an account that does not own a copy of Minecraft.
 - 405 [Method Not Allowed]    - If this error occurs, you have not set the request method to PUT.

Applications like this one will send http requests to the mojang servers when a ign drops, it is timed
and the program tries its best to get the ign! <10.970 > 11.100 [200] Got IGN wow!>

# How To Find A Good Delay

To be honest there isnt much to it, all you can do is get a estimate number and work off it.
delay testers like the one mcsn has just times the send / recv intervals and does math / rounding
to give an approximate value, do NOT expect the values it gives to be final as it can be inaccurate!

# Which Sniper Should I Use?

Its all based on your preference, based on the free snipers that are still maintained either work fine.
its all based on who YOU trust more with your data, whos sniper has features you enjoy and which sniper has
in your opinion a more stable request system. as a beginner this isnt easy to spot (obviously).

Look for snipers that are recently updated on github
 - Go onto github search "Minecraft Name Sniper" and go to recently updated
 - Find a sniper!

Some snipers that are good for starting out are:
 - MCsniperGO
 - Snipesharp
 - MCSN

I suggest joining the mcsnipergo discord <3 https://discord.com/invite/yp69ZqtxNk

# Competitive Sniping? The Hard Truth.

When it comes to sniping you also need to expect private snipers funded by individuals who make a profit from selling there igns.
Because of this you cannot expect to get a OG name or anything too special that easily. Some snipers that are private and paid can be
good investments.

Sniping is expensive, dont come into the buisness expecting to get every username, you may get lucky (which sniping is all luck) and get a 3 char

Please have fun and GL with sniping!

`)
					return nil
				},
			},
		},
		Name:    "MCSN",
		Usage:   "A minecraft name sniper dedicated to premium free services.",
		Version: "6.10",
	}

	app.Run(os.Args)
}
