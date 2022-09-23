# Backer

Rsync wrapper to do multiple transfers from a single yaml file.

## Why

I'm sure there is **lots** of implementations already exist doing the same, but i wanted to some golang rsync wrapper that will do some routine backups using my own yaml file specs.
Since i already have multiple directories that i update daily, i wanted to always make sure that i have them synced and moved to a persistent location since some of these files are not and won't be hosted on github.

## How

Backer will generate a new config file under `~/.config/backer/config.yaml` there you're expected to modify and enter your sources & destinations.

### Config file

```yaml
backer:
  # you can add multiple transfer levels
  transfer:
    # make sure that both source & destination paths are absolute
    - source: /home/john/personal-stuff/
      destination: /run/media/john/some-mounted-hdd/personal/
    - source: /home/john/pictures/
      destination: /run/media/john/some-mounted-hdd/pictures/
  # these exclude args will be later fed to rsync --exclude, by default this will be empty
  exclude:
    - log
    - logs
    - "*.log"
    - "node_modules"
    - vendor
    - .git
    - .next
  # here you can just add rsync options you want to supply by default this will be empty
  rsync_options:
    - -avAXEWSlHh
    - --no-compress
    - --info=progress2
```

### The CLI

Currently you could use backer directly to copy from a single source to single destination:

```bash
backer /home/john/a/ /home/john/b/
```

You can also add rsync options:

```bash
backer /home/john/a/ /home/john/b/ -o '-avAXEWSlHh' -o '--no-compress'
```

## Cron

Currently since I'm using Arch - by the way ;) - I use systemd to periodically run backer every hour, you can find the systemd service and timer under `./systemd`

You will need to make some proper changes to the service to point to the proper user & group and backer location

```
[Service]
Type=simple
User=mrgeek # change this
Group=mrgeek # change this
ExecStart=/home/mrgeek/go/bin/backer # change this
```

then you can run `./sync` to sync and start backer timer

```bash
cd ./systemd
./sync
```

## Install

If you already have go installed then:

```bash
go install github.com/omarahm3/backer@latest
```

Or from releases page
