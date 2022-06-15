# NRKReader

Playing around with scraping articles from nrk.no. Will list all headlines and link to the article.

## Usage

```
$ go run ./cmd/nrkreader
```

# NRKTracker

Stores the headlines and day they were seen at `$HOME/.config/nrktracker/news.json.gz`.
We can track the headlines by running the tracker daily. For instance in a cronjob.

## Usage
```
$ go run ./cmd/nrktracker
```