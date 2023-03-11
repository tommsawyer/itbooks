![](https://github.com/tommsawyer/itbooks/blob/main/banner.png)
[![Scrape](https://github.com/tommsawyer/itbooks/actions/workflows/scrape.yaml/badge.svg)](https://github.com/tommsawyer/itbooks/actions/workflows/scrape.yaml)
[![Publish](https://github.com/tommsawyer/itbooks/actions/workflows/publish.yaml/badge.svg)](https://github.com/tommsawyer/itbooks/actions/workflows/publish.yaml)

Fully automated telegram [channel](https://t.me/new_it_books) that publishes all new and upcoming books about IT.

### How it works
Every day we collect information from book publishers that publish books on IT topics. These books are stored in the database. Another script checks daily for new books and publish them into telegram channel.

### Local development

1. Install golang, docker and docker-compose
2. Obtain telegram bot token as described [here](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)
3. Create telegram channel to test books publishing
4. Set up environment variables:
```
export TELEGRAM_TOKEN=token_from_step_2
export TELEGRAM_CHANNEL=@channel_name_from_step_3
```
5. Run `make postgres` to spin up testing database
6. Run `make migrate` to run migrations on testing database
7. Run `make scrape` to scrape publishers and save books into postgres
8. Run `make publish` to publish one of the scraped books

You should see one of the books published in your telegram channel at this moment. Explore `./build/itbooks --help` to see what other commands do we have.
