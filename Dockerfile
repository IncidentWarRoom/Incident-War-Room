# --- build stage ---
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
# Build-time API base URLs (Vite inlines these into the bundle).
# Override at build time with --build-arg if your deployment needs different values.
ARG VITE_INCIDENT_API_BASE=http://localhost:8080
ARG VITE_REPORT_API_BASE=http://localhost:8000
ENV VITE_INCIDENT_API_BASE=$VITE_INCIDENT_API_BASE
ENV VITE_REPORT_API_BASE=$VITE_REPORT_API_BASE
RUN npm run build

# --- serve stage ---
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
