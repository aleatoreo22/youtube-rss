# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Project Overview

This is a YouTube RSS aggregator application with two main components:

1. **youtube-rss-api** (Go): A REST API server that fetches YouTube RSS feeds, stores content in SQLite, and serves it via HTTP endpoints
2. **youtube_rss_ui** (Flutter): A mobile/desktop Flutter app that consumes the API to display video content

# Architecture

## youtube-rss-api (Go)

The API follows a clean architecture pattern with these layers:

### Directory Structure
- `cmd/api/main.go` - Application entry point
- `internal/api/routes/` - HTTP route registration
- `internal/api/handler/` - HTTP request handlers
- `internal/service/` - Business logic layer
- `internal/repository/` - Data access layer
- `internal/domain/` - Domain models
- `package/` - External packages (yt-rss library, database)

### Key Components

**Service Layer** (`internal/service/`):
- `Service` - Aggregates all services
- `UserService` - User CRUD operations
- `ChannelService` - YouTube channel management
- `UserChannelService` - Links users to channels
- `ContentService` - Content fetching and syncing from YouTube RSS feeds

**Repository Layer** (`internal/repository/`):
- Uses SQLite database via `package/database/`
- `Repository` aggregates all repositories
- `ContentRepository`, `UserRepository`, `ChannelRepository`, etc.

**Handler Layer** (`internal/api/handler/`):
- `Handler` - Main HTTP handler with error handling
- `UserHandler`, `ContentHandler` - Specific endpoint handlers

### API Endpoints

The API runs on port 1234 with these endpoints:

| Method | Path | Description |
|--------|------|-------------|
| GET | /ping | Health check |
| GET | /user/:id/content/:date | Get content for a user by date |
| GET | /user/:id/content?page=:p&limit=:l | Paginated content for user |
| POST | /user/:id/channel | Add a YouTube channel to user |
| GET | /content/:date | Get all content by date |
| GET | /user/:id/channel | Get channels for a user |
| DELETE | /user/:id/channel | Remove a channel from user |

### Background Sync

The API has a background goroutine (`getYoutubeRssContent`) that:
- Runs every 5 minutes
- Fetches latest content from all subscribed YouTube RSS feeds
- Upserts new/updated content into the database

## youtube_rss_ui (Flutter)

A Flutter mobile/desktop app with:

### Directory Structure
- `lib/main.dart` - App entry point with navigation bar
- `lib/page/video.dart` - Main video feed page
- `lib/service/youtuberss/` - API client and models

### Features
- Dark theme with yellow accent color
- Bottom navigation bar with: Videos (active), History, User
- Infinite scroll with pagination
- Auto-refresh button
- Opens YouTube videos in browser

### API Client

The `YoutubeRSS` service in `lib/service/youtuberss/youtuberss.dart`:
- Base URL: `http://onehome:1234` (Docker network default)
- Methods: `getContentByDate`, `getUserContentByDate`, `addUserChannel`, `ping`, `getUserContentPaginated`

### Models

- `Content` - Video content with id, url, channel, title, image, date
- `Pagination` - Paginated response with page, limit, total, totalPages, items

# Commands

## Building

### API (Go)
```bash
cd youtube-rss-api
go build -o api ./cmd/api
```

### Flutter App
```bash
cd youtube_rss_ui
flutter pub get
flutter run
```

## Running

### API
```bash
cd youtube-rss-api
./api
# or
go run ./cmd/api
```

### Flutter App
```bash
cd youtube_rss_ui
flutter run
```

## Testing

### Flutter Tests
```bash
cd youtube_rss_ui
flutter test
```

## Linting

### Flutter
```bash
cd youtube_rss_ui
flutter analyze
```

# Development Workflow

1. **Start the API first** - The Flutter app depends on the API being running
2. **API runs on port 1234** - Make sure it's accessible from the Flutter app's network
3. **Database** - SQLite database at `youtube-rss-api/database.db`
4. **Docker** - If using Docker, API runs as `onehome:1234`

# Common Issues

1. **Connection refused** - API not running or wrong host/port
2. **No content found** - API hasn't synced yet (runs every 5 minutes) or no channels added
3. **Image not loading** - YouTube may have changed thumbnail URLs; check `maxresdefault` vs `hqdefault`

# Debugging

### API Debugging
- Check logs for sync errors
- Use `curl http://localhost:1234/ping` to verify API is running
- Check database: `youtube-rss-api/database.db`

### Flutter Debugging
- Check network connectivity to API
- Verify API base URL matches actual server address
- Check for CORS issues if running API on different host

# Data Models

## User
- id, username, email

## Channel
- id, rssUrl, title

## UserChannel
- id, userId, channelId

## Content
- id, url, channel, title, image, date

## Pagination
- page, limit, total, totalPages, items
