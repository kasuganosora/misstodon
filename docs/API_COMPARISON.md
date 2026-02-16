# API 接口对比文档

本文档对比 Mastodon、Misskey 和 Misstodon 三个项目的 API 接口，以 Mastodon API 为基准进行分析。

## 概述

| 项目 | 类型 | API 数量 | 基础路径 |
|------|------|----------|----------|
| Mastodon | Ruby on Rails | 200+ | `/api/v1`, `/api/v2` |
| Misskey | Node.js/Fastify | 435+ | `/api` (POST 方法为主) |
| Misstodon | Go/Gin | 37 | `/api/v1`, `/api/v2` |

> **注意**: Misstodon 是一个 Mastodon API 兼容层，用于将 Misskey 后端伪装成 Mastodon API。

---

## 1. 账户相关 API (Accounts)

### 1.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/accounts/verify_credentials` | GET | `POST /api/i` | ✅ 已实现 |
| `/api/v1/accounts/update_credentials` | PATCH | `POST /api/i/update` | ✅ 已实现 |
| `/api/v1/accounts/:id` | GET | `POST /api/users/show` | ✅ 已实现 |
| `/api/v1/accounts/:id/statuses` | GET | `POST /api/users/notes` | ✅ 已实现 |
| `/api/v1/accounts/:id/followers` | GET | `POST /api/users/followers` | ✅ 已实现 |
| `/api/v1/accounts/:id/following` | GET | `POST /api/users/following` | ✅ 已实现 |
| `/api/v1/accounts/relationships` | GET | `POST /api/users/relation` | ✅ 已实现 |
| `/api/v1/accounts/lookup` | GET | `POST /api/users/show` (by username) | ✅ 已实现 |
| `/api/v1/accounts/:id/follow` | POST | `POST /api/following/create` | ✅ 已实现 |
| `/api/v1/accounts/:id/unfollow` | POST | `POST /api/following/delete` | ✅ 已实现 |
| `/api/v1/accounts/:id/mute` | POST | `POST /api/mute/create` | ✅ 已实现 |
| `/api/v1/accounts/:id/unmute` | POST | `POST /api/mute/delete` | ✅ 已实现 |
| `/api/v1/favourites` | GET | `POST /api/users/reactions` | ✅ 已实现 |
| `/api/v1/follow_requests` | GET | `POST /api/following/requests/list` | ✅ 已实现 |

### 1.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/accounts` | GET | `POST /api/users` | 批量获取账户信息 |
| `/api/v1/accounts/search` | GET | `POST /api/users/search` | 搜索账户 |
| `/api/v1/accounts/:id/block` | POST | `POST /api/blocking/create` | 屏蔽用户 |
| `/api/v1/accounts/:id/unblock` | POST | `POST /api/blocking/delete` | 取消屏蔽 |
| `/api/v1/accounts/:id/pin` | POST | - | 在个人资料中推荐 |
| `/api/v1/accounts/:id/unpin` | POST | - | 取消推荐 |
| `/api/v1/accounts/:id/note` | POST | - | 设置私人备注 |
| `/api/v1/accounts/:id/remove_from_followers` | POST | `POST /api/following/invalidate` | 移除粉丝 |
| `/api/v1/accounts/familiar_followers` | GET | - | 获取熟悉的关注者 |
| `/api/v1/accounts/:id/lists` | GET | - | 获取用户所在列表 |
| `/api/v1/accounts/:id/identity_proofs` | GET | - | 身份证明 |
| `/api/v1/accounts/:id/featured_tags` | GET | - | 精选标签 |

---

## 2. 状态相关 API (Statuses)

### 2.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/statuses` | POST | `POST /api/notes/create` | ✅ 已实现 |
| `/api/v1/statuses/:id` | GET | `POST /api/notes/show` | ✅ 已实现 |
| `/api/v1/statuses/:id/context` | GET | 无 (返回空结果) | ✅ 已实现 (Stub，返回空 ancestors/descendants) |
| `/api/v1/statuses/:id/favourite` | POST | `POST /api/notes/reactions/create` | ✅ 已实现 |
| `/api/v1/statuses/:id/unfavourite` | POST | `POST /api/notes/reactions/delete` | ✅ 已实现 |
| `/api/v1/statuses/:id/bookmark` | POST | `POST /api/notes/favorites/create` | ✅ 已实现 |
| `/api/v1/statuses/:id/unbookmark` | POST | `POST /api/notes/favorites/delete` | ✅ 已实现 |
| `/api/v1/bookmarks` | GET | `POST /api/i/favorites` | ✅ 已实现 |

### 2.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/statuses/:id` | DELETE | `POST /api/notes/delete` | 删除状态 |
| `/api/v1/statuses/:id` | PUT/PATCH | `POST /api/notes/update` | 编辑状态 |
| `/api/v1/statuses/:id/reblogged_by` | GET | `POST /api/notes/renotes` | 获取转发者 |
| `/api/v1/statuses/:id/favourited_by` | GET | - | 获取点赞者 |
| `/api/v1/statuses/:id/reblog` | POST | `POST /api/notes/create` (renote) | 转发 |
| `/api/v1/statuses/:id/unreblog` | POST | `POST /api/notes/delete` | 取消转发 |
| `/api/v1/statuses/:id/mute` | POST | - | 静音对话 |
| `/api/v1/statuses/:id/unmute` | POST | - | 取消静音 |
| `/api/v1/statuses/:id/pin` | POST | - | 置顶 |
| `/api/v1/statuses/:id/unpin` | POST | - | 取消置顶 |
| `/api/v1/statuses/:id/history` | GET | - | 编辑历史 |
| `/api/v1/statuses/:id/source` | GET | - | 获取原始内容 |
| `/api/v1/statuses/:id/translate` | POST | `POST /api/notes/translate` | 翻译 |
| `/api/v1/statuses/:id/quotes` | GET | - | 引用列表 |

---

## 3. 时间线相关 API (Timelines)

### 3.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/timelines/home` | GET | `POST /api/notes/timeline` | ✅ 已实现 |
| `/api/v1/timelines/public` | GET | `POST /api/notes/global-timeline` / `local-timeline` | ✅ 已实现 |
| `/api/v1/timelines/tag/:hashtag` | GET | `POST /api/notes/search-by-tag` | ✅ 已实现 |

### 3.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/timelines/list/:id` | GET | - | 列表时间线 |
| `/api/v1/timelines/link` | GET | - | 链接时间线 |

---

## 4. 趋势相关 API (Trends)

### 4.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/trends/tags` | GET | `POST /api/hashtags/trend` | ✅ 已实现 |
| `/api/v1/trends/statuses` | GET | `POST /api/notes/featured` | ✅ 已实现 |

### 4.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/trends` | GET | - | 趋势 (旧版，同 tags) |
| `/api/v1/trends/links` | GET | - | 热门链接 |

---

## 5. 媒体相关 API (Media)

### 5.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/media` | POST | `POST /api/drive/files/create` | ✅ 已实现 |
| `/api/v2/media` | POST | `POST /api/drive/files/create` | ✅ 已实现 |

### 5.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/media/:id` | GET | `POST /api/drive/files/show` | 获取媒体信息 |
| `/api/v1/media/:id` | PUT/PATCH | `POST /api/drive/files/update` | 更新媒体 |
| `/api/v1/media/:id` | DELETE | `POST /api/drive/files/delete` | 删除媒体 |

---

## 6. 通知相关 API (Notifications)

### 6.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/notifications` | GET | `POST /api/i/notifications` | ✅ 已实现 |

### 6.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/notifications/:id` | GET | - | 获取单个通知 |
| `/api/v1/notifications/clear` | POST | `POST /api/notifications/mark-all-as-read` | 清除所有通知 |
| `/api/v1/notifications/:id/dismiss` | POST | - | 忽略单个通知 |
| `/api/v1/notifications/unread_count` | GET | - | 未读计数 |
| `/api/v1/notifications/policy` | GET/PUT | - | 通知策略 |
| `/api/v1/notifications/requests` | GET | - | 通知请求 |

---

## 7. 实例相关 API (Instance)

### 7.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/instance` | GET | `POST /api/meta` | ✅ 已实现 |
| `/api/v1/instance/peers` | GET | 无 (返回空数组) | ✅ 已实现 (Stub) |
| `/api/v1/custom_emojis` | GET | `POST /api/emojis` | ✅ 已实现 |

### 7.2 未实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/instance/rules` | GET | - | 实例规则 |
| `/api/v1/instance/domain_blocks` | GET | - | 被屏蔽域名 |
| `/api/v1/instance/extended_description` | GET | - | 扩展描述 |
| `/api/v1/instance/privacy_policy` | GET | - | 隐私政策 |
| `/api/v1/instance/terms_of_service` | GET | - | 服务条款 |
| `/api/v1/instance/activity` | GET | - | 活动统计 |
| `/api/v1/instance/translation_languages` | GET | - | 支持的翻译语言 |
| `/api/v2/instance` | GET | `POST /api/meta` | v2 实例信息 |

---

## 8. OAuth 相关 API

### 8.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/oauth/authorize` | GET | `POST /api/auth/session/generate` | ✅ 已实现 (Legacy Auth) |
| `/oauth/token` | POST | `POST /api/auth/session/userkey` | ✅ 已实现 (Legacy Auth) |
| `/oauth/redirect` | GET | - | ✅ 已实现 (非标准端点，用于 OAuth 回调重定向) |

### 8.2 未实现的接口

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/oauth/revoke` | POST | 撤销令牌 |
| `/oauth/token/info` | POST | 获取令牌信息 |

---

## 9. 流媒体相关 API (Streaming)

### 9.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/streaming` | WebSocket | WebSocket | ✅ 已实现 |

---

## 10. 应用相关 API (Apps)

### 10.1 已实现的接口

| Mastodon API | HTTP 方法 | Misskey 等效端点 | Misstodon 实现状态 |
|--------------|-----------|------------------|-------------------|
| `/api/v1/apps` | POST | `POST /api/app/create` | ✅ 已实现 |

### 10.2 未实现的接口

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/apps/verify_credentials` | GET | 验证应用凭据 |

---

## 11. 其他未实现的重要 API

### 11.1 搜索 (Search)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v2/search` | GET | `POST /api/search` | 搜索 |

### 11.2 列表 (Lists)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/lists` | GET | - | 获取列表 |
| `/api/v1/lists/:id` | GET | - | 获取单个列表 |
| `/api/v1/lists` | POST | `POST /api/users/lists/create` | 创建列表 |
| `/api/v1/lists/:id` | PUT | `POST /api/users/lists/update` | 更新列表 |
| `/api/v1/lists/:id` | DELETE | `POST /api/users/lists/delete` | 删除列表 |
| `/api/v1/lists/:id/accounts` | GET | `POST /api/users/lists/members` | 列表成员 |
| `/api/v1/lists/:id/accounts` | POST | `POST /api/users/lists/push` | 添加成员 |
| `/api/v1/lists/:id/accounts` | DELETE | `POST /api/users/lists/pull` | 移除成员 |

### 11.3 过滤器 (Filters)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/filters` | GET/POST | 过滤器列表/创建 |
| `/api/v2/filters` | GET/POST | v2 过滤器 |
| `/api/v2/filters/:id` | GET/PUT/DELETE | 单个过滤器操作 |

### 11.4 对话 (Conversations)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/conversations` | GET | `POST /api/messaging/history` | 对话列表 |

### 11.5 标签 (Tags)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/tags/:id` | GET | `POST /api/hashtags/show` | 标签信息 |
| `/api/v1/tags/:id/follow` | POST | `POST /api/hashtags/follow` | 关注标签 |
| `/api/v1/tags/:id/unfollow` | POST | `POST /api/hashtags/unfollow` | 取消关注 |
| `/api/v1/followed_tags` | GET | `POST /api/i/followed-tags` | 已关注标签 |

### 11.6 投票 (Polls)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/polls/:id` | GET | - | 获取投票 |
| `/api/v1/polls/:id/votes` | POST | `POST /api/notes/polls/vote` | 投票 |

### 11.7 公告 (Announcements)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/announcements` | GET | `POST /api/announcements` | 公告列表 |
| `/api/v1/announcements/:id/dismiss` | POST | - | 忽略公告 |

### 11.8 建议 (Suggestions)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/suggestions` | GET | `POST /api/users/recommendation` | 用户建议 |
| `/api/v2/suggestions` | GET | `POST /api/users/recommendation` | v2 用户建议 |

### 11.9 精选标签 (Featured Tags)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/featured_tags` | GET/POST | 精选标签 |
| `/api/v1/featured_tags/suggestions` | GET | 建议的精选标签 |

### 11.10 计划状态 (Scheduled Statuses)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/scheduled_statuses` | GET | 计划发布的状态 |
| `/api/v1/scheduled_statuses/:id` | GET/PUT/DELETE | 单个计划状态 |

### 11.11 域名屏蔽 (Domain Blocks)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/domain_blocks` | GET/POST/DELETE | - | 域名屏蔽 |

### 11.12 推送订阅 (Push Subscriptions)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/push/subscription` | GET/POST/PUT/DELETE | 推送订阅 |

### 11.13 标记 (Markers)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/markers` | GET/POST | 阅读位置标记 |

### 11.14 偏好设置 (Preferences)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/preferences` | GET | 用户偏好 |

### 11.15 目录 (Directory)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/directory` | GET | 用户目录 |

### 11.16 屏蔽和静音 (Blocks & Mutes)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/blocks` | GET | `POST /api/blocking/list` | 屏蔽列表 |
| `/api/v1/mutes` | GET | `POST /api/mute/list` | 静音列表 |

### 11.17 举报 (Reports)

| Mastodon API | HTTP 方法 | Misskey 等效端点 | 说明 |
|--------------|-----------|------------------|------|
| `/api/v1/reports` | POST | `POST /api/users/report-abuse` | 举报用户 |

### 11.18 推荐 (Endorsements)

| Mastodon API | HTTP 方法 | 说明 |
|--------------|-----------|------|
| `/api/v1/endorsements` | GET | 推荐的用户 |

---

## 12. Well-Known 端点

| 端点 | Mastodon | Misskey | Misstodon |
|------|----------|---------|-----------|
| `/.well-known/nodeinfo` | ✅ | ✅ | ✅ 已实现 |
| `/.well-known/webfinger` | ✅ | ✅ | ✅ 已实现 |
| `/.well-known/host-meta` | ✅ | ✅ | ✅ 已实现 |

---

## 13. NodeInfo 端点

| 端点 | Mastodon | Misskey | Misstodon |
|------|----------|---------|-----------|
| `/nodeinfo/2.0` | ✅ | ✅ | ✅ 已实现 |

---

## 14. Admin API

> **注意**: Misstodon 目前未实现任何 Admin API

### Mastodon Admin API 列表 (部分)

| Mastodon API | 功能 |
|--------------|------|
| `/api/v1/admin/accounts` | 账户管理 |
| `/api/v1/admin/reports` | 报告管理 |
| `/api/v1/admin/domain_allows` | 允许的域名 |
| `/api/v1/admin/domain_blocks` | 屏蔽的域名 |
| `/api/v1/admin/email_domain_blocks` | 屏蔽的邮箱域名 |
| `/api/v1/admin/ip_blocks` | IP 屏蔽 |
| `/api/v1/admin/trends/*` | 趋势管理 |
| `/api/v1/admin/measures` | 统计指标 |
| `/api/v1/admin/dimensions` | 维度统计 |
| `/api/v1/admin/retention` | 留存统计 |

---

## 15. 统计汇总

### 15.1 Misstodon 实现进度

| 类别 | Mastodon 接口数 | Misstodon 已实现 | 实现率 |
|------|----------------|------------------|--------|
| 账户 (Accounts) | 26 | 14 | 53.8% |
| 状态 (Statuses) | 22 | 8 | 36.4% |
| 时间线 (Timelines) | 5 | 3 | 60.0% |
| 趋势 (Trends) | 4 | 2 | 50.0% |
| 媒体 (Media) | 5 | 2 | 40.0% |
| 通知 (Notifications) | 7 | 1 | 14.3% |
| 实例 (Instance) | 11 | 3 | 27.3% |
| OAuth | 4 | 3 | 75.0% |
| 流媒体 (Streaming) | 1 | 1 | 100% |
| 应用 (Apps) | 2 | 1 | 50.0% |
| 其他 | 80+ | 0 | 0% |

### 15.2 Misskey 特有功能 (Mastodon 无对应)

| Misskey 端点 | 功能说明 |
|--------------|----------|
| `/api/antennas/*` | 天线功能 (高级过滤器) |
| `/api/channels/*` | 频道功能 |
| `/api/chat/*` | 聊天功能 |
| `/api/clips/*` | 剪辑功能 |
| `/api/pages/*` | 页面功能 |
| `/api/flash/*` | Play 功能 |
| `/api/gallery/*` | 画廊功能 |
| `/api/roles/*` | 角色系统 |
| `/api/i/registry/*` | 注册表存储 |
| `/api/drive/folders/*` | 文件夹管理 |
| `/api/renote-mute/*` | 转发静音 |
| `/api/i/2fa/*` | 双因素认证详细配置 |

---

## 16. 版本信息

| 项目 | 当前版本 | 协议 |
|------|----------|------|
| Mastodon | v4.3+ | AGPL-3.0 |
| Misskey | 2024.x | AGPL-3.0 |
| Misstodon | - | AGPL-3.0 |

---

## 附录 A: Misskey API 调用方式

Misskey API 与 Mastodon 有显著不同：

1. **HTTP 方法**: 默认使用 `POST`，而非 RESTful 的 GET/PUT/DELETE
2. **参数传递**: 通过 JSON body 传递，而非 URL 参数
3. **认证**: 使用 `i` 参数传递 access token
4. **响应格式**: JSON

### 示例对比

**Mastodon:**
```http
GET /api/v1/accounts/123
Authorization: Bearer <token>
```

**Misskey:**
```http
POST /api/users/show
Content-Type: application/json

{
  "userId": "123",
  "i": "<token>"
}
```

---

## 附录 B: Misstodon 代理服务器机制

Misstodon 支持通过以下方式指定代理的 Misskey 服务器（按优先级从高到低排列）：

1. **Host 头**: `mt_misskey_io.example.com` (编码格式：`_` 代表 `.`，`__` 代表原始 `_`)
2. **路径参数**: `/:proxyServer/api/v1/...`
3. **URL 参数**: `?server=misskey.io`
4. **请求头**: `X-Proxy-Server: misskey.io`
5. **配置文件**: `FallbackServer` (兜底默认服务器)

---

*文档生成日期: 2025年2月*
*最后更新日期: 2026年2月*
