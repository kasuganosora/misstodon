# misstodon

[![爱发电](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fafdian.net%2Fapi%2Fuser%2Fget-profile%3Fuser_id%3D75e549844b5111ed8df552540025c377&query=%24.data.user.name&label=%E7%88%B1%E5%8F%91%E7%94%B5&color=%23946ce6)](https://afdian.net/a/gizmo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gizmo-ds/misstodon?style=flat-square)
[![Build images](https://img.shields.io/github/actions/workflow/status/gizmo-ds/misstodon/images.yaml?branch=main&label=docker%20image&style=flat-square)](https://github.com/gizmo-ds/misstodon/actions/workflows/images.yaml)
[![License](https://img.shields.io/github/license/gizmo-ds/misstodon?style=flat-square)](./LICENSE)

Misskey Mastodon-compatible APIs — use any Mastodon client (Elk, Phanpy, Ivory, Tusky, Ice Cubes...) with your [Misskey](https://github.com/misskey-dev/misskey) instance.

> **Note**
> This project is under active development. Most common Mastodon client workflows (login, timelines, posting, notifications) are functional.

## Features

- Pure Go implementation, no Node.js/JavaScript runtime dependency
- OAuth authorization flow compatible with standard Mastodon clients
- MFM (Misskey Flavored Markdown) to HTML conversion
- Proxy mode: one deployment can serve multiple Misskey instances
- Self-hosting mode: bind to a single Misskey instance via `fallback_server`

## Quick Start

### Self-Hosting (Recommended)

The typical setup: deploy Misstodon on a subdomain (e.g. `mapi.a.com`) pointing to your Misskey instance (`a.com`).

1. Edit `config.toml`:

```toml
[proxy]
fallback_server = "a.com"

[server]
bind_address = ":3000"
```

2. Run:

```bash
# Binary
./misstodon -c config.toml

# Or Docker Compose
docker-compose up -d
```

3. Point your Mastodon client to `https://mapi.a.com` and log in.

### Docker Compose

Download [docker-compose.yml](https://github.com/gizmo-ds/misstodon/raw/main/docker-compose.yml), set `MISSTODON_FALLBACK_SERVER` to your Misskey instance domain, then:

```bash
docker-compose up -d
```

> **Important**
> For security and privacy, always use HTTPS. Configure a TLS certificate or use Misstodon's AutoTLS feature.

## Advanced Usage

### Domain Name Prefixing Scheme

For proxy deployments serving multiple Misskey instances, specify the target via domain prefix:

1. Replace `_` with `__` in the Misskey domain
2. Replace `.` with `_`
3. Prepend `mt_`
4. Append your Misstodon base domain

Example: `misskey.io` → `mt_misskey_io.liuli.lol`

```bash
curl https://mt_misskey_io.liuli.lol/api/v1/instance | jq .
```

### Instance Specification via Query Parameter

```bash
curl 'https://misstodon.example.com/api/v1/instance?server=misskey.io' | jq .
```

### Instance Specification via Header

```bash
curl https://misstodon.example.com/api/v1/instance -H 'x-proxy-server: misskey.io' | jq .
```

## API Coverage

<details>
<summary>Supported Endpoints</summary>

### Discovery & Auth

- [x] `GET` /.well-known/webfinger
- [x] `GET` /.well-known/nodeinfo
- [x] `GET` /.well-known/host-meta
- [x] `GET` /nodeinfo/2.0
- [x] `GET` /oauth/authorize
- [x] `POST` /oauth/token
- [x] `POST` /api/v1/apps
- [x] `GET` /api/v1/apps/verify_credentials

### Instance

- [x] `GET` /api/v1/instance
- [x] `GET` /api/v2/instance
- [x] `GET` /api/v1/instance/peers
- [x] `GET` /api/v1/instance/rules
- [x] `GET` /api/v1/custom_emojis

### Accounts

- [x] `GET` /api/v1/accounts/lookup
- [x] `GET` /api/v1/accounts/:id
- [x] `GET` /api/v1/accounts/verify_credentials
- [x] `GET` /api/v1/accounts/relationships
- [x] `GET` /api/v1/accounts/:id/following
- [x] `GET` /api/v1/accounts/:id/followers
- [x] `GET` /api/v1/accounts/:id/lists
- [x] `GET` /api/v1/accounts/:id/featured_tags
- [x] `POST` /api/v1/accounts/:id/follow
- [x] `POST` /api/v1/accounts/:id/unfollow
- [x] `POST` /api/v1/accounts/:id/mute
- [x] `POST` /api/v1/accounts/:id/unmute
- [x] `GET` /api/v1/follow_requests
- [x] `POST` /api/v1/follow_requests/:id/authorize
- [x] `POST` /api/v1/follow_requests/:id/reject
- [x] `GET` /api/v1/bookmarks
- [x] `GET` /api/v1/favourites

### Statuses

- [x] `POST` /api/v1/statuses (text, media, polls, reply, visibility, CW)
- [x] `GET` /api/v1/statuses/:id
- [x] `GET` /api/v1/statuses/:id/context
- [x] `POST` /api/v1/statuses/:id/favourite
- [x] `POST` /api/v1/statuses/:id/unfavourite
- [x] `POST` /api/v1/statuses/:id/bookmark
- [x] `POST` /api/v1/statuses/:id/unbookmark
- [x] `POST` /api/v1/statuses/:id/reblog
- [x] `POST` /api/v1/statuses/:id/unreblog

### Timelines

- [x] `GET` /api/v1/timelines/home
- [x] `GET` /api/v1/timelines/public
- [x] `GET` /api/v1/timelines/tag/:hashtag

### Notifications

- [x] `GET` /api/v1/notifications
- [x] `GET` /api/v1/notifications/unread_count

### Polls

- [x] `GET` /api/v1/polls/:id
- [x] `POST` /api/v1/polls/:id/votes

### Search

- [x] `GET` /api/v2/search

### Media

- [x] `POST` /api/v1/media
- [x] `POST` /api/v2/media

### Trends

- [x] `GET` /api/v1/trends/statuses
- [x] `GET` /api/v1/trends/tags

### Other

- [x] `GET` /api/v1/announcements
- [x] `GET` /api/v1/conversations
- [x] `GET` /api/v1/preferences
- [x] `GET` /api/v1/markers
- [x] `POST` /api/v1/markers
- [x] `GET` /api/v1/suggestions
- [x] `GET` /api/v2/suggestions
- [x] `POST` /api/v1/reports
- [x] `GET` /api/v1/blocks
- [x] `GET` /api/v1/mutes
- [x] `GET` /api/v1/lists
- [x] `GET` /api/v1/domain_blocks
- [x] `GET` /api/v1/filters
- [x] `GET` /api/v2/filters
- [x] `GET` /api/v1/featured_tags
- [x] `GET` /api/v1/followed_tags
- [x] `GET` /api/v1/endorsements
- [x] `GET` /api/v1/scheduled_statuses
- [ ] `WS` /api/v1/streaming

</details>

## Information for Developers

[Contributing](./CONTRIBUTING.md) Information about contributing to this project.

## Sponsors

[![Sponsors](https://afdian-connect.deno.dev/sponsor.svg)](https://afdian.net/a/gizmo)

## Contributors

![Contributors](https://contributors.liuli.lol/gizmo-ds/misstodon/contributors.svg?align=left)
