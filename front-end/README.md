# Go Auth Template - Frontend

A modern Next.js frontend for the Go Auth Template with WebAuthn (passkey) support and magic link email authentication.

## Features

- 🔐 **Dual Authentication**: Email magic links + WebAuthn passkeys
- 🚀 **Next.js 15**: Latest version with App Router and Server Components
- 🎨 **Modern UI**: Tailwind CSS with Radix UI components
- 🔒 **Type Safety**: Full TypeScript support
- 📱 **Responsive**: Mobile-first design
- ⚡ **Fast**: Turbopack for development

## Tech Stack

- **Framework**: Next.js 15 with App Router
- **Styling**: Tailwind CSS v4
- **UI Components**: Radix UI
- **Authentication**: Custom session management with JWT
- **WebAuthn**: Browser-native passkey support
- **TypeScript**: Full type safety

## Quick Start

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env.local
   ```
   Edit `.env.local` with your values.

3. **Start development server**:
   ```bash
   npm run dev
   ```

4. **Open in browser**: [http://localhost:3000](http://localhost:3000)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_URL` | Go backend API URL | `http://localhost:8080` |
| `NEXT_PUBLIC_API_URL` | Public API URL for client-side | `http://localhost:8080` |
| `NEXTAUTH_URL` | App URL | `http://localhost:3000` |
| `NEXTAUTH_SECRET` | Session secret | Required in production |

## Project Structure

```
src/
├── app/                    # Next.js App Router
│   ├── api/               # API routes (proxy to Go backend)
│   ├── dashboard/         # Protected dashboard page
│   ├── login/            # Authentication pages
│   └── layout.tsx        # Root layout
├── components/
│   ├── auth/             # Authentication components
│   └── ui/               # Reusable UI components
├── lib/
│   ├── auth/             # Authentication utilities
│   └── utils.ts          # Utility functions
├── auth.ts               # Server actions for auth
└── middleware.ts         # Next.js middleware
```

## Authentication Flow

### Email Magic Link
1. User enters email
2. System checks if user exists and has passkeys
3. Sends magic link code via email
4. User enters code to authenticate

### WebAuthn Passkeys
1. User enters email  
2. System detects available passkeys
3. User selects passkey authentication
4. Browser prompts for biometric/passkey
5. Instant authentication

## Available Scripts

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Fix ESLint errors
- `npm run type-check` - Check TypeScript types
- `npm run clean` - Clean build artifacts

## API Integration

The frontend communicates with the Go backend through:
- **Next.js API routes** - Proxy requests for client components
- **Server actions** - Direct backend calls for server components

All API calls include proper client type detection for multi-platform support.

## Deployment

1. **Build the application**:
   ```bash
   npm run build
   ```

2. **Set production environment variables**

3. **Deploy to your platform** (Vercel, Netlify, etc.)

## Customization

### Styling
- Modify `tailwind.config.js` for design system changes
- Update components in `src/components/ui/` for UI modifications

### Authentication
- Customize auth flows in `src/components/auth/`
- Update session management in `src/lib/auth/session.ts`

### Branding
- Replace logos and icons in `public/`
- Update metadata in `src/app/layout.tsx`

## Contributing

1. Follow the existing code style
2. Run `npm run lint` and `npm run type-check` before committing
3. Ensure all tests pass
4. Update documentation as needed

## License

[Your License Here]