# ramiro.me — Implementation Plan

## Context

Ramiro's site is being rebuilt from a Hugo-based storytelling coaching site into a **life operating system** — a personal platform that holds all facets of who Ramiro is: AI developer, algorithmic trader, entrepreneur, clairvoyant meditation teacher, polyglot.

The current site (Hugo + Story theme) has no authentication, no admin, no member area, and is tightly coupled to a single storytelling identity. The rebuild needs three access tiers (public, admin, members), a rich content model, and eventually an AI-powered Instagram content system.

Ramiro prefers **Go** over Node.js. Accent color: **electric violet**.

---

## Tech Stack: Go + Templ + HTMX + Postgres

**Why migrate from Hugo:** Hugo can't handle authentication, admin CRUD, member sections, or dynamic content. The project needs a unified full-stack application.

**Why vanilla Go (not PocketBase, not Hugo+separate app):**
- **Single binary deployment** — `go build` → one binary. Ship anywhere. Database managed externally (Supabase Postgres).
- **Full control** — every route, query, and template is yours. No framework coupling.
- **Templ + HTMX** — type-safe compiled HTML templates + dynamic interactions without a JS framework. Perfect for content pages + admin forms/lists.
- **Go's stdlib is strong** — `net/http`, `encoding/xml` (RSS), `crypto/bcrypt` (auth), `embed` (static assets). Minimal dependencies.
- **PocketBase rejected** because: the moment you want custom admin UX (journal, podcast manager, content editor), you're building a custom Go app anyway, just coupled to PocketBase's abstractions.
- **Hugo+Go hybrid rejected** because: two systems to sync, auth boundary is messy, RSS feed split across services, operational complexity.

### Full Stack

| Concern | Choice | Rationale |
|---|---|---|
| **Router** | Chi (`go-chi/chi`) | Lightweight, idiomatic, great middleware support |
| **Templates** | Templ | Type-safe, compiled Go templates, component model |
| **Interactivity** | HTMX | Dynamic UX without JS framework — forms, search, filters, live preview |
| **Styling** | Tailwind CSS | Dark theme, utility-first, fast iteration |
| **Database** | PostgreSQL on Supabase | Managed, automatic backups, no data loss risk, generous free tier |
| **DB Driver** | `pgx` (jackc/pgx) | Fast, pure Go Postgres driver, connection pooling |
| **Migrations** | `golang-migrate/migrate` | SQL migration files, up/down |
| **Markdown** | Goldmark | Fast, extensible, CommonMark compliant |
| **Auth** | `bcrypt` + `gorilla/sessions` | Session cookies, middleware-gated routes |
| **Payments** | Stripe Go SDK | Server-side Checkout Sessions + webhooks |
| **Audio** | Keep RedCircle (needs verification) | Hosts 80+ episodes; private RSS feed status TBD |
| **Email** | ConvertKit (existing) | Keep for now |
| **Deployment** | Render, Fly.io, or DigitalOcean | Single binary, DB on Supabase (separate) |

---

## Content Model

### Subdomains

| Domain | Purpose |
|---|---|
| `ramiro.me` | Public site — what the world sees |
| `x.ramiro.me` | Private admin — only Ramiro |
| `circle.ramiro.me` | Inner circle — paid members (89 max + Ramiro = 90) |

All three served by the same Go binary. Chi router distinguishes by `Host` header or a simple subdomain middleware.

### Database Schema (PostgreSQL)

```sql
-- Core content table (stories, ideas, celebrations, challenges, blog, projects)
CREATE TABLE content (
    id          TEXT PRIMARY KEY,
    type        TEXT NOT NULL, -- 'story','project','idea','celebration','challenge','relationship','blog','private'
    title       TEXT NOT NULL,
    slug        TEXT UNIQUE NOT NULL,
    summary     TEXT,
    body        TEXT NOT NULL, -- Markdown
    visibility  TEXT NOT NULL DEFAULT 'private', -- 'public','members','private'
    status      TEXT NOT NULL DEFAULT 'draft',   -- 'draft','published'
    tags        TEXT, -- JSON array
    category    TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Podcast episodes (migrated from Hugo markdown files)
CREATE TABLE episodes (
    id              TEXT PRIMARY KEY,
    season          TEXT,
    episode         TEXT,
    title           TEXT NOT NULL,
    slug            TEXT UNIQUE NOT NULL,
    date            TIMESTAMP NOT NULL,
    mp3_url         TEXT NOT NULL, -- RedCircle UUID
    guid            TEXT UNIQUE NOT NULL,
    length          INTEGER, -- bytes
    duration        INTEGER, -- seconds
    episode_type    TEXT, -- 'guest','ramiro','lyrics','bonus','intro'
    itunes_summary  TEXT,
    body            TEXT, -- Show notes (markdown)
    visibility      TEXT NOT NULL DEFAULT 'public',
    status          TEXT NOT NULL DEFAULT 'published',
    tags            TEXT, -- JSON array
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users (admin + members)
CREATE TABLE users (
    id                  TEXT PRIMARY KEY,
    email               TEXT UNIQUE NOT NULL,
    password_hash       TEXT, -- bcrypt, NULL for magic-link-only members
    name                TEXT,
    role                TEXT NOT NULL DEFAULT 'member', -- 'admin','member'
    languages           TEXT, -- JSON array, must have 2+ for members
    stripe_customer_id  TEXT,
    subscription_status TEXT DEFAULT 'none', -- 'none','active','cancelled','past_due'
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions
CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id),
    expires_at TIMESTAMP NOT NULL
);

-- Journal entries (admin-only private notes)
CREATE TABLE journal (
    id         TEXT PRIMARY KEY,
    title      TEXT,
    body       TEXT NOT NULL,
    mood       TEXT,
    tags       TEXT, -- JSON array
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Migration from Hugo
- 80+ episode `.md` files → parse YAML front matter → insert into `episodes` table
- Training episodes (`content/access/`) → insert with `visibility: 'members'`
- Images (`static/img/`) → copy to `static/img/` in new project
- Stripe price IDs, RedCircle channel tokens, ConvertKit form IDs → environment variables

---

## Project Structure

```
ramiro.me/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, wire up routes + middleware
├── internal/
│   ├── handler/                 # HTTP handlers grouped by section
│   │   ├── public.go            # Home, About, Projects, Stories, Contact
│   │   ├── podcast.go           # Podcast listing, episode detail, RSS feed
│   │   ├── admin.go             # Dashboard, content browser, journal, editor
│   │   ├── member.go            # Member area handlers
│   │   ├── auth.go              # Login, logout, session management
│   │   └── api.go               # Stripe webhooks, HTMX endpoints
│   ├── middleware/
│   │   ├── auth.go              # RequireAdmin, RequireMember middleware
│   │   └── logging.go
│   ├── model/                   # Data types
│   │   ├── content.go
│   │   ├── episode.go
│   │   ├── user.go
│   │   └── journal.go
│   ├── store/                   # Database access layer
│   │   ├── store.go             # Postgres connection (pgx) + helpers
│   │   ├── content.go           # Content CRUD
│   │   ├── episode.go           # Episode CRUD
│   │   ├── user.go              # User CRUD
│   │   └── journal.go           # Journal CRUD
│   ├── rss/
│   │   └── podcast.go           # iTunes-compatible RSS XML generation
│   └── markdown/
│       └── render.go            # Goldmark markdown → HTML
├── templates/                   # Templ component files (.templ)
│   ├── layout/
│   │   ├── base.templ           # HTML shell (head, body, scripts)
│   │   ├── nav.templ            # Navigation (public)
│   │   ├── footer.templ
│   │   └── admin.templ          # Admin layout (sidebar + content)
│   ├── pages/
│   │   ├── home.templ
│   │   ├── about.templ
│   │   ├── projects.templ
│   │   ├── stories.templ
│   │   ├── contact.templ
│   │   └── login.templ
│   ├── podcast/
│   │   ├── list.templ           # Episode listing
│   │   ├── detail.templ         # Single episode + player
│   │   └── player.templ         # Audio player component
│   ├── admin/
│   │   ├── dashboard.templ
│   │   ├── content_browser.templ
│   │   ├── content_editor.templ # Markdown editor with HTMX preview
│   │   ├── journal.templ
│   │   └── podcast_manager.templ
│   ├── member/
│   │   ├── home.templ
│   │   └── join.templ
│   └── components/              # Reusable Templ components
│       ├── card.templ
│       ├── tag.templ
│       ├── pagination.templ
│       └── form.templ
├── migrations/                  # SQL migration files
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── static/                      # Static assets (embedded via go:embed)
│   ├── css/
│   │   └── app.css              # Tailwind output
│   ├── js/
│   │   └── htmx.min.js          # HTMX library
│   └── img/                     # Images (migrated from Hugo)
├── tools/
│   └── migrate/
│       └── main.go              # Hugo content → Postgres migration script
├── tailwind.config.js
├── go.mod
├── go.sum
├── Makefile                     # build, dev, migrate, seed commands
└── .env.example
```

---

## Design System

**Color palette — Near-black + Electric Violet:**
- Background: `#0a0a0f`
- Surface: `#141420` (cards, panels)
- Accent primary: `#7c3aed` (electric violet)
- Accent light: `#a78bfa` (hover states)
- Text primary: `#f0f0f5`
- Text secondary: `#9090a0`
- Border: `#1e1e2e`

**Typography:**
- Headings: Space Grotesk (bold, modern, slightly technical)
- Body: Inter (clean, readable)
- Code/mono: JetBrains Mono

**Tailwind config:**
```js
theme: {
  extend: {
    colors: {
      bg: { DEFAULT: '#0a0a0f', surface: '#141420', border: '#1e1e2e' },
      accent: { DEFAULT: '#7c3aed', light: '#a78bfa' },
      text: { DEFAULT: '#f0f0f5', muted: '#9090a0' },
    },
    fontFamily: {
      heading: ['"Space Grotesk"', 'sans-serif'],
      body: ['Inter', 'system-ui', 'sans-serif'],
      mono: ['"JetBrains Mono"', 'monospace'],
    }
  }
}
```

**Principles:** Mobile-first. Dark by default. Confident whitespace. Subtle CSS transitions (no heavy JS animation library needed). No clutter. Feels like one person, not a committee.

---

## Phased Implementation

### Phase 1 — Foundation
1. Initialize Go module, set up Chi router, Templ, HTMX
2. Set up Supabase Postgres + migrations
3. Create base layout template (dark theme, fonts, nav, footer)
4. Build Home page (hero: "Intuition, Intelligence, Impact", intro sections)
5. Build About page (polymath journey — timeline format)
6. Build Contact page (services, paid call CTA, social links)
7. Set up Tailwind build pipeline (via standalone CLI)
8. Set up Makefile (`make dev`, `make build`, `make migrate`)
9. Deploy to Render/Fly.io (single binary, Postgres on Supabase)

**Key files:** `cmd/server/main.go`, `internal/handler/public.go`, `templates/layout/base.templ`, `templates/pages/home.templ`, `migrations/001_initial.up.sql`, `Makefile`

### Phase 2 — Content & Projects
1. Implement Goldmark markdown rendering pipeline
2. Build content CRUD in store layer
3. Build Stories listing page (with tag filtering via HTMX)
4. Build story detail page (markdown rendered to HTML)
5. Build Projects grid + detail pages
6. Seed initial public content

### Phase 3 — Podcast
1. Write Hugo-to-Postgres migration script (`tools/migrate/main.go`) — parse 80+ episode `.md` files, insert into `episodes` table
2. Build podcast listing page (season grouping, filtering)
3. Build episode detail page with HTML5 audio player (RedCircle URLs)
4. Build iTunes-compatible RSS feed (`internal/rss/podcast.go`) using `encoding/xml`
5. Add URL rewrites: `/eps/index.xml` → RSS handler (preserve podcast subscriber URLs)
6. Validate RSS feed output matches current Hugo feed structure

### Phase 4 — Auth & Admin
1. Implement session-based auth (bcrypt + gorilla/sessions)
2. Build login page
3. Add `RequireAdmin` middleware for `x.ramiro.me` subdomain routes
4. Build admin dashboard (content stats, recent entries, quick actions)
5. Build content browser (filter by type/visibility/status/tags — HTMX-powered)
6. Build content editor (textarea + HTMX live markdown preview via Goldmark)
7. Build life journal (private entries CRUD)
8. Build podcast manager (episode list, create/edit metadata)

### Phase 5 — Members Section (circle.ramiro.me)
1. Implement Stripe server-side Checkout Sessions + webhooks for $10k/year annual subscription
2. Build member join page — language check (2+ required), capacity check (89 max), Stripe checkout
3. Add `RequireMember` middleware for circle subdomain routes
4. Build member home page (premium feel — this is a $10k community)
5. Member-visible content feeds (stories, celebrations, challenges with `visibility: 'members'`)
6. Show "X of 89 spots taken" on join page
7. Training content migration (from current `content/access/` section)

### Phase 6 — Instagram AI Content System
1. Content calendar UI in admin (HTMX-powered month/week view)
2. AI-assisted content generator (call Claude API from Go, return draft captions/carousel text)
3. Brand kit storage (avatar, colors, fonts, tone-of-voice — stored in DB)
4. Asset templates for carousels, captions, reels

---

## Content Migration Strategy

### Preserved:
- 80+ episode markdown files → parsed and imported to Postgres via migration script
- RedCircle audio URLs (mp3 UUIDs) → stored in `episodes.mp3_url` (need to verify URLs still work)
- RSS feed URLs (`/eps/index.xml`) → served by new RSS handler at same path
- Episode GUIDs → preserved exactly (podcast apps use these as primary keys)
- Images (`static/img/`) → copied to new `static/img/`
- Stripe price IDs → moved to environment variables

### Dropped:
- Hugo theme (Story) — replaced by Templ + Tailwind dark design
- All 63 Hugo partials — replaced by Templ components
- MediaElementPlayer + jQuery — replaced by HTML5 `<audio>` + minimal JS
- `config.toml` — replaced by Go config + env vars

---

## Verification Plan

### Phase 1:
- `make dev` starts the server, pages render at localhost
- Dark theme displays correctly
- Mobile responsive at 375px, 768px, 1024px
- Navigation works on all screen sizes
- `make build` produces single binary
- Binary deploys and runs on target platform

### Phase 3:
- Migration script imports all 80+ episodes without errors
- Episode listing shows all episodes with correct metadata
- Audio player plays from RedCircle URLs
- RSS feed at `/eps/index.xml` validates at W3C feed validator
- RSS output matches current Hugo RSS structure (diff test)
- Episode GUIDs are identical to current feed

### Phase 4:
- Admin login works with password
- Unauthenticated users get 302 redirect from `x.ramiro.me/*`
- Content CRUD works (create, read, update, delete)
- Journal entries never appear on public site
- Markdown preview renders correctly via HTMX

### Phase 5:
- Stripe checkout creates subscription
- Webhook updates user subscription status
- Member login gates circle subdomain routes
- Capacity enforced at 89 members
- Training content only visible to members

---

## Open Questions

1. **Hosting**: Render.com, Fly.io, or DigitalOcean droplet?
2. **Podcast RSS**: Need to verify RedCircle private RSS feed status. Are the mp3 UUIDs still valid?
3. **Tagline**: "Intuition, Intelligence, Impact" as primary. Where (if anywhere) to use "helping the crazy ones be as FUCKING crazy as they can be"?
4. **Member onboarding**: At $10k/year, should there be a personal vetting step (call with Ramiro) before payment, or direct checkout?

## Resolved Decisions

- **Database**: PostgreSQL on Supabase (managed, backed up, no data loss risk)
- **Subdomains**: `x.ramiro.me` (private admin), `circle.ramiro.me` (members)
- **Membership**: $10,000 USD/year annual subscription via Stripe. Max 89 members + Ramiro = 90.
- **Accent color**: Electric violet (`#7c3aed`)
- **Framework**: Vanilla Go + Chi + Templ + HTMX (no heavy framework)
