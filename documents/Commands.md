### Hello, welcome to the commands guide.

A little FYI before i continue, the command prefix depends on the version your using heres the different versions.
- `go run .` | src
- `main.exe` | executable

I will be using the src prefix for this tutorial!

# Commands

To get started with the commands heres a few important details

![Meow](https://media.discordapp.net/attachments/924456247237947424/939398249377329212/unknown.png)

Each command has its own description.

With that out of the way lets continue;

- `go run . ping` This command returns the mean of your delay. you can use this for your delay in the next cmd.
- `go run . snipe -u username -d delay` This command snipes a username your going for, the -d is your delay input EG 10.
- `go run . auto 3c -d delay` attempts to snipe upcoming 3 character names (any and all)
- `go run . auto 3l -d delay` attempts to snipe upcoming 3l names, "lol".
- `go run . auto 3n -d delay` attempts to snipe upcoming 3n names, "123".
- `go run . auto list -d delay` snipes names from your names.txt file.

- `go run . proxy -u username -d delay` Snipes names using proxies.
- `go run . proxy 3c -d delay` this attempts to snipe all upcoming 3c using proxies.
- `go run . proxy 3l -d delay` this attempts to snipe all upcoming 3l using proxies.
- `go run . proxy 3n -d delay` this attempts to snipe all upcoming 3n using proxies.

# Fun Commands and Extra Info

These commands are mainly added for preference or user experience.

- `go run . namemc -n nameofskinart -i skinartimg.jpg` This applies skin art to your account, it handles any image type PNG JPG etc. For this to function, navigate to your namemc profile from your browser, everytime your skin is changed refresh your profiles page until namemc caches your skin!~
- `go run . turbo -u name` this turbos a name, essentially sending requests every minute or so (typically to a name that isnt available) this is just in case a drop drops for no reason because of mojang support.
- `go run . bot` this command launches a discord bot that you host, you will need to add your bots key to the config, "DiscordBotKey".

Other info, you do NOT need to add a -d value to the auto snipe commands.

- `go run . auto 3c` this functions normally like the other auto commands, but instead uses the build in delay calculator to grab your delay for snipes.
- `go run . proxy 3c` this functions like the normal auto 3c functions, instead it uses thedelay it gets from the autooffset function as stated above aswell ^^
