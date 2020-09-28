# feedy

> Scrape RSS feeds and add articles to pocket

## Getting started

Clone repo and download dependencies:

```bash
git clone git@github.com:eirsyl/feedy.git
cd feedy
make vendor
```

Build application:

```bash
make build
```

## Usage

Go the releases page to download a precompiled version of feedy!

```bash
feedy --help
RSS feed scraper

Usage:
  feedy [command]

Available Commands:
  feed        Manage feeds
  help        Help about any command
  login       Authenticate with pocket
  scrape      Scrape watched feeds

Flags:
      --concurrency int     feeds to scrape concurrent (default 10)
  -c, --configFile string   config file path
  -h, --help                help for feedy
      --version             version for feedy

Use "feedy [command] --help" for more information about a command.
```

### Add feed

```bash
feedy -c /var/lib/feedy/config.db feed add <feed> <tag1> <tag2>
```

### List subscribed feeds

```bash
feedy -c /var/lib/feedy/config.db feed list
```

### Run the scraper as a system service

```bash
feedy -c /var/lib/feedy/config.db scrape --autostop=false
```

Systemd file for running the service in the background:

```bash
[Unit]
Description=Feedy RSS scraper
After=network.target

[Service]
Type=simple
User=feedy
ExecStart=/usr/local/sbin/feedy -c /root/config.db scrape --autostop=false
CPUAccounting = yes
MemoryAccounting = yes

[Install]
WantedBy=multi-user.target
```
