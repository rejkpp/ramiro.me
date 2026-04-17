# Project

| Field | Value |
| --- | --- |
| Project | ramiro.me |
| Created | 2026-04-11 21:01 |
| Last Updated | 2026-04-17 09:30 |
| Stage | Development |

## Goals

Rebuild ramiro.me from a Hugo-based storytelling coaching site into a **life operating system** — a personal platform that holds all facets of who Ramiro is: AI developer, algorithmic trader, entrepreneur, clairvoyant meditation teacher, polyglot.

Three access tiers served by a single Go binary:
- `ramiro.me` — public site
- `x.ramiro.me` — private admin (Ramiro only)
- `circle.ramiro.me` — paid members (89 max + Ramiro, $10k/year)

Preserves the existing podcast (80 episodes, RedCircle-hosted, RSS feed URLs and GUIDs intact).

## Architecture Decisions

- **Vanilla Go + Chi + Templ + HTMX** over PocketBase or Hugo+Go hybrid — single-binary deploy, full control, type-safe templates, dynamic UX without a JS framework
- **PostgreSQL on Supabase** over SQLite or self-hosted — managed backups, no data-loss risk, generous free tier
- **pgx driver** over database/sql+lib/pq — faster, native Postgres features, connection pooling
- **Goldmark** for Markdown rendering — fast, CommonMark compliant, extensible
- **bcrypt + gorilla/sessions** for auth — session cookies, middleware-gated routes
- **Stripe Go SDK** for $10k/year member subscriptions — server-side Checkout Sessions + webhooks
- **Tailwind CSS (standalone CLI)** for styling — dark theme, utility-first, no Node toolchain
- **Three subdomains, one binary** — Chi host-based routing distinguishes public/admin/member

## Preferences

<!-- Behavioral preferences — not architecture decisions. Loaded at session start via hook. -->
- Go preferred over Node.js
- Accent color: amber primary (`#f59e0b`) + pink secondary (`#ec4899`)
- Dark theme by default, mobile-first
- Minimal dependencies; favor stdlib

## Open Issues

- **Podcast audio URLs are dead** — RedCircle CDN returns 404 for episode mp3s (subscription expired). Need to recover from RedCircle before Phase 3.

## Future Work (Podcast)

**Audio migration:**
1. Re-subscribe to RedCircle temporarily
2. Download all 80 mp3s locally
3. Upload to Cloudflare R2 (`audio.ramiro.me` subdomain)
4. Update episode URLs in database
5. Cancel RedCircle

**Private podcast feed architecture:**
- Public episodes: stored in R2, direct URLs, included in public RSS
- Private episodes: stored in R2 with **private bucket + signed URLs**
- Each circle member gets a unique feed token stored in `users.feed_token`
- Private RSS endpoint: `GET /circle/feed/{token}/rss.xml`
- Go app validates token → generates RSS with **signed URLs** (valid 48h)
- Podcast apps cache and play; if member churns, revoke token → feed dies

This is how Patreon/Supercast work. Full control, no third-party dependency.

## Decisions Log

Reverse chronological. Format: `YYYY-MM-DD — [Decision] — [Rationale]`

- 2026-04-17 — Migrate Home page prose to markdown — 4 card descriptions + featured project blurb live in `content/home/**/*.md`; structured UI (cards, bio, tags) stays in templ. Projects/Booking skipped (placeholders + pricing one-liners don't benefit)
- 2026-04-17 — Enable Goldmark Typographer extension — authors write `---`, `...`, straight quotes in markdown; renderer produces proper em-dashes, ellipses, and smart quotes without special keyboard input
- 2026-04-17 — Goldmark markdown content pipeline — prose lives in `content/**/*.md`, loaded + cached at startup via `//go:embed`, rendered into templ pages via `@content.MustGet(key)`. Structured UI (timelines, language bars, place badges, pricing tables) stays in templ; only prose moves to markdown
- 2026-04-13 — Paid booking flow: 30min ($150), 60min ($250), Clairvoyant reading ($300) — Calendly integration, no free discovery calls
- 2026-04-13 — Contact page → Booking page — simplified user flow, remove email from site
- 2026-04-12 — Dual-accent color system: amber (#f59e0b) primary + pink (#ec4899) secondary — amber for tech/action elements, pink for consciousness/decorative, gradient for hero moments. Replaces electric violet.
- 2026-04-11 — Self-host RSS feeds from Go app — eliminates RedCircle 2-account workaround, gives full RSS control, no yearly renewal
- 2026-04-11 — Rename default branch `master` → `main` — modernize 2021-era repo
- 2026-04-11 — Membership: $10,000/year, cap at 89 members + Ramiro = 90 — premium inner-circle positioning
- 2026-04-11 — Subdomains: `x.ramiro.me` (admin), `circle.ramiro.me` (members) — clean host-based separation
- 2026-04-11 — Stack: Go + Chi + Templ + HTMX + Postgres — rejected PocketBase (custom admin UX forces a Go app anyway) and Hugo+Go hybrid (two systems to sync)
