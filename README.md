# GoBlog

A simple blog built with **Go** and **PostgreSQL**, fully containerized using **Podman** and Compose.

## Quick Start

The easiest way to run the entire application (the Go server and PostgreSQL database) is using the included shell scripts which launch a native Podman pod:

```bash
# Start everything in a Podman pod
./deploy.sh
```

*(Alternatively, if you have Podman Compose installed, you can use `podman-compose up -d --build`)*

Then open **http://localhost:8080** in your browser.

## Pages

| URL | Description |
|-----|-------------|
| `/` | Blog homepage — lists all published posts |
| `/post/{slug}` | Single post view |
| `/admin` | Admin panel — manage all posts |
| `/admin/new` | Create a new post |
| `/admin/edit/{id}` | Edit an existing post |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `bloguser` | Database user |
| `DB_PASSWORD` | `blogpass` | Database password |
| `DB_NAME` | `blogdb` | Database name |
| `PORT` | `8080` | Server port |

## Cleanup

```bash
# Stop and remove the Podman pod
./cleanup.sh
```

## Project Structure

```
go-blog/
├── main.go              # Entry point, routing, DB connection
├── handlers/blog.go     # HTTP handlers (home, post, admin CRUD)
├── models/post.go       # Post model & database queries
├── templates/           # Go HTML templates
│   ├── layout.html      # Base layout (nav, footer)
│   ├── home.html        # Homepage with post cards
│   ├── post.html        # Single post view
│   ├── admin.html       # Admin post list
│   └── edit.html        # Create/edit post form
├── static/style.css     # Stylesheet (dark theme)
├── migrations/          # SQL migrations
├── Dockerfile           # Minimal multi-stage build for Go app
├── docker-compose.yml   # Podman compose configuration
├── deploy.sh            # Native Podman pod deployment script
└── cleanup.sh           # Cleanup script
```
