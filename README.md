# Incident War Room - Frontend

Read-only web dashboard for incidents managed through the Telegram bot.
Built with React + Vite, talking directly to `incident-service`'s HTTP API
(see `docs/api/incident-service-openapi.yaml`).

This implements the four pages described in the feature list document:
incident list, incident overview, timeline, and photo gallery.

## Running locally

You need Node.js 18+ and the incident-service running (default
`http://localhost:8080`).

```bash
cd frontend
npm install
cp .env.example .env   # adjust the API base URLs if needed
npm run dev
```

This starts the dev server at `http://localhost:5173`.

## Building for production

```bash
npm run build
```

Output goes to `dist/`. `npm run preview` serves that build locally if you
want to sanity-check it before deploying.

## Running with Docker

```bash
docker build -t incident-war-room-frontend .
docker run -p 3000:80 incident-war-room-frontend
```

The API base URLs are baked in at build time (Vite inlines them into the
bundle), so pass them as build args if they differ from the defaults:

```bash
docker build \
  --build-arg VITE_INCIDENT_API_BASE=http://incident-service:8080 \
  -t incident-war-room-frontend .
```

### Adding this to the project's docker-compose.yml

A snippet like this should slot into the existing `docker-compose.yml` at
the repo root:

```yaml
  frontend:
    build:
      context: ./frontend
      args:
        VITE_INCIDENT_API_BASE: http://localhost:8080
    ports:
      - "3000:80"
    depends_on:
      - incident-service
```

Note: because Vite inlines the API base URL into the JS bundle at build
time, `VITE_INCIDENT_API_BASE` needs to be a URL the *browser* can reach
(e.g. `localhost:8080`), not the internal Docker network name
(`http://incident-service:8080`) - unless you're proxying through nginx.
Worth double-checking with whoever sets up the deployment.

## What's intentionally not here yet

Carried over from the open questions in the feature list document - these
were simplifications made to get a working version built, not final
decisions:

- **No report generation trigger.** The UI only displays `reportUrl` when
  it's already present on an incident; it never calls
  `report-service`'s `/api/v1/reports/generate`.
- **No authentication.** Matches the API specs, which don't define any.
- **No demo-data fallback.** Your Claude Design prototype showed a "demo
  data" banner when the API was unreachable. This build shows a plain
  error state with a retry button instead, to keep the first version
  simpler. Easy to add back if the team wants it.
- **No Telegram deep-linking.** `chatId` and `topicId` are available on
  the incident object but aren't used for a "view in Telegram" link yet.
- **Navigation chrome is simplified.** The prototype used a sidebar on
  the detail page; this build uses a single top bar across all pages
  instead, to reduce the first cut's complexity. Layout-only change, no
  functionality lost.

## Project structure

```
src/
  api.js          API client (GET requests to incident-service only)
  utils.js        Formatting/color helpers (severity, status, dates, avatars)
  styles.css      Global styles
  components/     Small reusable pieces (badges, avatar, skeleton, states)
  pages/
    IncidentList.jsx
    IncidentDetail.jsx
    TimelineTab.jsx
    PhotosTab.jsx
```
