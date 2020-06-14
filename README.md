# Simple crypto currencies bot

> Minimalistic telegram-bot for observing crypto-assets through messages.

## Commands

`/start` - show usage manual

`/add <symbol> <symbol>` - add crypto currency traiding pair to report

`/remove` - remove existing pair from reports

`/interval <minutes>` - setup notifications interval

`/suspend` - stop notifications

`/resume` - resume notifications

## Deploy

You should have `docker` and `docker-compose`.

Copy `.env.example` file and setup your-own variables.

Then, just run `docker-compose up -d`.

## Localization

At this moment bot supports only Russian language. If you want to have your language support -- feel free to open issue.
