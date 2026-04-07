# LingProxy Claude Design System Refactor

## Overview

Refactor the LingProxy admin panel to adopt the Claude (Anthropic) design system, creating a warm, human-centered interface for developer users.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Target User | Developer/Engineer | Needs clear resource monitoring and cost tracking |
| Color Mode | Light mode | Claude's signature parchment tone creates human warmth |
| Refactor Depth | Deep refactor | Full Claude aesthetic requires complete consistency |
| Sidebar Style | Dark sidebar (Claude colors) | Maintains navigation contrast while using Claude's warm dark tones |
| Data Display | Soft tables + Card-style + Emphasized typography | Developer-focused but warm aesthetic |

## Color System

### Primary Colors

| Role | Claude Color | Hex | Usage |
|------|-------------|-----|-------|
| Brand/CTA | Terracotta Brand | `#c96442` | Primary buttons, brand moments |
| Coral Accent | Coral Accent | `#d97757` | Links on dark, secondary emphasis |
| Error | Error Crimson | `#b53333` | Error states |
| Focus | Focus Blue | `#3898ec` | Input focus rings (accessibility) |

### Surface Colors

| Role | Claude Color | Hex | Usage |
|------|-------------|-----|-------|
| Page Background | Parchment | `#f5f4ed` | Main page background - warm cream |
| Card Surface | Ivory | `#faf9f5` | Cards, elevated containers |
| Pure White | Pure White | `#ffffff` | Button surfaces, max contrast |
| Button BG | Warm Sand | `#e8e6dc` | Secondary button backgrounds |
| Dark Surface | Dark Surface | `#30302e` | Dark containers, nav borders |
| Deep Dark | Anthropic Near Black | `#141413` | Dark theme background, sidebar |

### Text Colors

| Role | Claude Color | Hex | Usage |
|------|-------------|-----|-------|
| Primary Text | Anthropic Near Black | `#141413` | Main text color |
| Secondary | Olive Gray | `#5e5d59` | Secondary body text |
| Tertiary | Stone Gray | `#87867f` | Tertiary text, metadata |
| Dark-on-light | Charcoal Warm | `#4d4c48` | Button text on warm surfaces |
| Light-on-dark | Warm Silver | `#b0aea5` | Text on dark surfaces |

### Border Colors

| Role | Claude Color | Hex | Usage |
|------|-------------|-----|-------|
| Light Border | Border Cream | `#f0eee6` | Standard light border |
| Prominent Border | Border Warm | `#e8e6dc` | Section dividers |
| Dark Border | Dark Surface | `#30302e` | Dark surface borders |

## Typography

### Font Families

| Role | Claude Font | Fallback |
|------|-------------|----------|
| Headlines | Anthropic Serif | `Georgia, serif` |
| Body/UI | Anthropic Sans | `Inter, system-ui, sans-serif` |
| Code | Anthropic Mono | `Consolas, monospace` |

### Typography Scale

| Role | Size | Weight | Line Height | Usage |
|------|------|--------|-------------|-------|
| Hero/Display | 64px | 500 | 1.10 | Dashboard hero |
| Section Heading | 52px | 500 | 1.20 | Page titles |
| Sub-heading Large | 36px | 500 | 1.30 | Card titles |
| Sub-heading | 32px | 500 | 1.10 | Section markers |
| Feature Title | 20px | 500 | 1.20 | Small headings |
| Body Large | 20px | 400 | 1.60 | Intro paragraphs |
| Body Standard | 16px | 400 | 1.50 | Standard body |
| Body Small | 15px | 400 | 1.60 | Compact text |
| Caption | 14px | 400 | 1.43 | Metadata |
| Label | 12px | 400 | 1.25 | Badges, labels |
| Code | 15px | 400 | 1.60 | Inline code |

## Component Styles

### Buttons

**Primary CTA (Terracotta)**
```css
background: #c96442;
color: #faf9f5;
border-radius: 8px;
padding: 8px 16px;
font-weight: 500;
ring-shadow: 0px 0px 0px 1px #c96442;
```

**Secondary (Warm Sand)**
```css
background: #e8e6dc;
color: #4d4c48;
border-radius: 8px;
padding: 8px 12px;
ring-shadow: 0px 0px 0px 1px #d1cfc5;
```

**Dark (Sidebar)**
```css
background: #30302e;
color: #faf9f5;
border-radius: 8px;
ring-shadow: 0px 0px 0px 1px #30302e;
```

### Cards & Containers

```css
background: #faf9f5;
border: 1px solid #f0eee6;
border-radius: 8px;
shadow: rgba(0,0,0,0.05) 0px 4px 24px;
padding: 24px;
```

### Tables

```css
border-radius: 8px;
border: 1px solid #f0eee6;
header-bg: #f8f7f3;
header-text: #5e5d59;
row-border: 1px 0px 0px #f0eee6;
```

### Inputs

```css
background: #faf9f5;
border: 1px solid #f0eee6;
border-radius: 12px;
padding: 8px 12px;
focus-ring: #3898ec;
```

### Sidebar

```css
background: #141413;
border-right: 1px solid #30302e;
active-item-bg: #c96442;
active-item-shadow: 0px 4px 12px rgba(201, 100, 66, 0.3);
```

## Layout Principles

### Spacing System

- Base unit: 8px
- Scale: 3px, 4px, 6px, 8px, 10px, 12px, 16px, 20px, 24px, 30px
- Section padding: 24px
- Header height: 72px

### Border Radius Scale

| Name | Size | Usage |
|------|------|-------|
| Sharp | 4px | Minimal inline |
| Subtle | 6px | Small buttons |
| Comfortable | 8px | Standard cards, buttons |
| Generous | 12px | Primary buttons, inputs |
| Very | 16px | Featured containers |
| Highly | 24px | Tag-like elements |
| Maximum | 32px | Hero containers |

## Shadow System

| Level | Shadow | Usage |
|-------|--------|-------|
| Flat | None | Background surfaces |
| Contained | 1px border | Standard cards |
| Ring | 0px 0px 0px 1px | Interactive states |
| Whisper | rgba(0,0,0,0.05) 0px 4px 24px | Elevated cards |

## Files to Modify

### Core Styles
- `frontend/src/style.css` - Global CSS variables and base styles
- `frontend/src/App.vue` - Root styles

### Layout Components
- `frontend/src/components/Layout.vue` - Header and main layout
- `frontend/src/components/Sidebar.vue` - Sidebar navigation

### Views (All 12 pages)
- `frontend/src/views/Dashboard.vue`
- `frontend/src/views/Login.vue`
- `frontend/src/views/Endpoints.vue`
- `frontend/src/views/Models.vue`
- `frontend/src/views/Tokens.vue`
- `frontend/src/views/Users.vue`
- `frontend/src/views/Policies.vue`
- `frontend/src/views/Logs.vue`
- `frontend/src/views/Requests.vue`
- `frontend/src/views/LLMResources.vue`
- `frontend/src/views/LLMResourceUsage.vue`
- `frontend/src/views/Settings.vue`

## Element Plus Overrides

Override Element Plus component styles to match Claude design:

- `el-button` - Terracotta primary, Warm Sand secondary
- `el-card` - Ivory background, cream borders
- `el-table` - Soft borders, warm header
- `el-input` - Generous radius, warm borders
- `el-menu` - Dark sidebar colors
- `el-dropdown` - Warm hover states
- `el-dialog` - Ivory surface, generous radius

## Success Criteria

1. All pages display warm parchment background
2. Typography uses serif for headings, sans for body
3. All buttons match Claude button styles
4. Sidebar uses Claude dark colors with terracotta active state
5. Tables have soft warm borders
6. Stat numbers use serif typography
7. Consistent ring-shadow depth system
8. Responsive behavior maintained