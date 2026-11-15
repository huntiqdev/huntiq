import type { Metadata, Viewport } from 'next';
import './globals.css';
import { ClientShell } from '@/components/layout/ClientShell';
import { LocaleProvider } from '@/context/LocaleContext';
import { AuthProvider } from '@/context/AuthContext';

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
  userScalable: true,
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#ffffff' },
    { media: '(prefers-color-scheme: dark)', color: '#09090b' },
  ],
  viewportFit: 'cover',
};

export const metadata: Metadata = {
  title: 'HUNTIQ — Ethiopia\'s Universal Marketplace',
  description: 'Buy, sell, and rent anything across Ethiopia. Find electronics, vehicles, property, agriculture, services, and more on HUNTIQ.',
  keywords: 'buy sell Ethiopia, Ethiopian marketplace, HUNTIQ, rent property Ethiopia, online market Ethiopia, Addis Ababa market',
  manifest: '/manifest.json',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'default',
    title: 'HUNTIQ',
  },
  icons: {
    icon: '/huntiq.png',
    shortcut: '/huntiq.png',
    apple: '/huntiq.png',
  },
};

// Runs before paint — prevents dark-mode flash without hydration errors.
const themeScript = `try{var t=localStorage.getItem('habesha_theme');if(t==='dark'||(!t&&window.matchMedia('(prefers-color-scheme: dark)').matches)){document.documentElement.classList.add('dark');}}catch(e){}`;

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{ __html: themeScript }}
        />
      </head>
      <body className="dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 transition-colors duration-300">
        <LocaleProvider>
          <AuthProvider>
            <ClientShell>{children}</ClientShell>
          </AuthProvider>
        </LocaleProvider>
      </body>
    </html>
  );
}
