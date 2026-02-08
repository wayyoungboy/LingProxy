# Getting Started with LingProxy Website

This guide will help you set up and run the LingProxy documentation website built with Docusaurus.

## Prerequisites

- Node.js 18.0 or higher
- npm or yarn

## Installation

1. **Install Dependencies**

```bash
cd website
npm install
```

## Development

2. **Start Development Server**

```bash
npm start
```

This command:
- Starts a local development server
- Opens a browser window automatically
- Most changes are reflected live without restarting the server

The website will be available at `http://localhost:5000`

## Build

3. **Build for Production**

```bash
npm run build
```

This command generates static content into the `build` directory.

4. **Preview Production Build**

```bash
npm run serve
```

This serves the built site locally for testing.

## Project Structure

```
website/
├── docs/                    # Documentation files
│   ├── introduction.md     # Introduction page
│   ├── quick-start.md      # Quick start guide
│   ├── configuration.md    # Configuration guide
│   ├── api-reference.md   # API reference
│   ├── streaming.md       # Streaming guide
│   ├── docker.md          # Docker deployment
│   ├── architecture.md    # Architecture guide
│   ├── development.md     # Development guide
│   └── analysis/          # Technical analysis docs
├── src/
│   ├── components/         # React components
│   ├── css/               # Custom styles
│   └── pages/             # Custom pages
├── static/                 # Static assets
├── i18n/                  # Internationalization
├── docusaurus.config.js   # Docusaurus configuration
└── sidebars.js            # Sidebar configuration
```

## Features

- ✅ **Multi-language Support**: English and Chinese
- ✅ **Search**: Built-in search functionality
- ✅ **Dark Mode**: Automatic dark mode support
- ✅ **Responsive**: Mobile-friendly design
- ✅ **Mermaid Diagrams**: Support for Mermaid diagrams
- ✅ **Code Highlighting**: Syntax highlighting for code blocks

## Customization

### Update Site Metadata

Edit `docusaurus.config.js`:
- `title`: Site title
- `tagline`: Site tagline
- `url`: Production URL
- `baseUrl`: Base URL path

### Add New Documentation

1. Add markdown files to `docs/` directory
2. Update `sidebars.js` to include new pages
3. Add frontmatter to markdown files:
   ```markdown
   ---
   sidebar_position: 1
   ---
   ```

### Customize Styling

Edit `src/css/custom.css` to customize colors, fonts, and styles.

## Deployment

### GitHub Pages

1. Update `docusaurus.config.js`:
   ```js
   url: 'https://lingproxy.github.io',
   baseUrl: '/lingproxy/',
   ```

2. Deploy:
   ```bash
   GIT_USER=<Your GitHub username> USE_SSH=true npm run deploy
   ```

### Docker

```bash
docker build -t lingproxy-website .
docker run -p 3000:80 lingproxy-website
```

### Other Platforms

The `build` directory contains static files that can be deployed to any static hosting service:
- Vercel
- Netlify
- AWS S3 + CloudFront
- Any web server

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
npm start -- --port 3001
```

### Build Errors

```bash
# Clear cache and rebuild
npm run clear
npm run build
```

### Missing Dependencies

```bash
# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

## Learn More

- [Docusaurus Documentation](https://docusaurus.io/docs)
- [Docusaurus Tutorial](https://docusaurus.io/docs/tutorial-basics)
