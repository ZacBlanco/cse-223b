# Web-crawler

We implement a Web Crawler inspired by *UbiCrawler* to evaluate our serverless modifications. Our web crawler downloads
web pages starting from some seed of domains. The web crawler uses the Actor model design and divides domain
responsibility between other actors. We use state to avoid re-downloading pages it has already visited. The Actor design
allows us to scale the number of concurrent downloads based on our specification. To avoid DoS-ing domains, we limit
each actor to asynchronously download from at most 30 different domains. Actors checkpoint periodically and can be
continued if stopped.

## Setup
1. `make`
1. Golang

## Tutorial

1. [Start up OpenWhisk](#Start up OpenWhisk)
1. [Upload functions](#Upload functions)
1. [Start the crawler](#Start the crawler)

### Start up OpenWhisk

What to do and what to expect.

### Upload functions

What to do and what to expect.

### Start the crawler

What to do and what to expect.

## Resources

1. [*UbiCrawler*](http://static.aminer.org/pdf/PDF/001/073/501/ubicrawler_a_scalable_fully_distributed_web_crawler.pdf)
