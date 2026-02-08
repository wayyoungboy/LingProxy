# LingProxy Website

This website is built using [Docusaurus](https://docusaurus.io/), a modern static website generator.

## Installation

```bash
npm install
```

## Local Development

```bash
npm start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

## Build

```bash
npm run build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Deployment

### Using GitHub Pages

1. Update `docusaurus.config.js`:
   ```js
   url: 'https://lingproxy.github.io',
   baseUrl: '/lingproxy/',
   ```

2. Run:
   ```bash
   GIT_USER=<Your GitHub username> USE_SSH=true npm run deploy
   ```

### Using Docker

```bash
docker build -t lingproxy-website .
docker run -p 5000:80 lingproxy-website
```

## Documentation Structure

- `docs/` - Main documentation files
- `docs/analysis/` - Technical analysis documents
- `i18n/` - Internationalization files
- `static/` - Static assets

## Internationalization

The website supports both English and Chinese:
- English: Default locale
- Chinese: Available via language switcher

## Learn More

- [Docusaurus Documentation](https://docusaurus.io/docs)
- [Docusaurus GitHub](https://github.com/facebook/docusaurus)
