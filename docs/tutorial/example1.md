# Example 1 - Using the CLI to create a podcast

## Register a Production

```shell
$ po new-show my-podcast
```

## Create a Show

```shell
$ po template show my-podcast
```

creates the following `show.yaml` file:

```yaml
apiVersion: v1
kind: show
metadata:
    name: my-podcast
    labels:
        block: "no"
        complete: "no"
        explicit: "no"
        guid: efc54386a6c6
        language: en_US
        type: Episodic
description:
    title: TITLE
    summary: SUMMARY
    link:
        uri: https://podops.dev/s/my-podcast
    category:
        name: Technology
        subcategory:
            - Podcasting
    owner:
        name: podcast owner
        email: hello@new-show.me
    author: new-show author
    copyright: podcast copyright
image:
    uri: https://cdn.podops.dev/c/default/cover.png
    rel: external
```

## Upload the media file

```shell
$ po upload first-episode.mp3
```

## Create a first Episode

```shell
$ po template -p my-podcast episode first-episode
```

creates the following `episode.yaml` file:

```yaml
apiVersion: v1
kind: episode
metadata:
    name: first-episode
    labels:
        block: "no"
        date: Sat, 13 Mar 2021 19:15:17 +0000
        episode: "1"
        explicit: "no"
        guid: c3ac4d7293a2
        parent_guid: efc54386a6c6
        season: "1"
        type: Full
description:
    title: Episode Title
    summary: Episode Subtitle or short summary
    episodeText: A long-form description of the episode with notes etc.
    link:
        uri: https://podops.dev/s/my-podcast/first-episode
    duration: 1
image:
    uri: https://cdn.podops.dev/c/default/episode.png
    rel: external
enclosure:
    uri: c3ac4d7293a2/first-episode.mp3
    rel: local
    type: audio/mpeg
    size: 1
```

## Create the resources

```shell
$ po create show.yaml
$ po create episode.yaml
```

## Build the Feed

```shell
$ po build
```