# go-elsewherr

This project was inspired from the existing project [Elsewherr](https://github.com/Adman1020/Elsewherr). 

## Introduction

### What is it?

go-elsewherr will see if your movies from Radarr (or series from Sonarr) are available on a streaming service,
and add a tag against the movie if it is.

### How does it work?

The script will check The Movie Database (https://www.themoviedb.org/) via their API,
which in turn uses Just Watch (https://www.justwatch.com/), to get all streaming services each movie is on.
It then adds this tag in Radarr (or Sonarr).

### Why?

I use Radarr (and Sonarr) as a manager for my movies (and series). When I decide to watch a movie
(or series), I can see directly from inside the entry where it is available.

## Environment variables

Configurations such as Radarr and Sonarr url/access key can be set through environment variables.
See the [.env.example](.env.example) file (read the comments).

## Setup

Copy [.env.example](.env.example) to `.env` and use the `docker compose` command to create and start the container.
```shell
docker compose up --build
```

You might want to set up a cronjob to run this regularly to keep the list up to date.







