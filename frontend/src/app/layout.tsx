import type { Metadata, Viewport } from 'next';
import './globals.css';
import { ClientShell } from '@/components/layout/ClientShell';
import { LocaleProvider } from '@/context/LocaleContext';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
  userScalable: true,
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#1769E0' },
    { media: '(prefers-color-scheme: dark)', color: '#091426' },
  ],
  viewportFit: 'cover',
};

export const metadata: Metadata = {
  title: 'HUNTIQ — Ethiopia\'s Universal Marketplace & Platform Services',
  description: 'Buy, sell, and rent anything across Ethiopia with verified trust, Telebirr integration, and broker-grade transparency.',
  keywords: 'buy sell Ethiopia, Ethiopian marketplace, HUNTIQ, rent property Ethiopia, online market Ethiopia, Addis Ababa market, Telebirr',
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
const themeScript = `try{var t=localStorage.getItem('huntiq_theme')||localStorage.getItem('habesha_theme');if(t==='dark'||((!t||t==='system')&&window.matchMedia('(prefers-color-scheme: dark)').matches)){document.documentElement.classList.add('dark');}else{document.documentElement.classList.remove('dark');}}catch(e){}`;

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{ __html: themeScript }}
          suppressHydrationWarning
        />
      </head>
      <body className="bg-slate-50 dark:bg-[#091426] text-slate-900 dark:text-slate-100 transition-colors duration-200">
        <ThemeProvider>
          <LocaleProvider>
            <AuthProvider>
              <ClientShell>{children}</ClientShell>
            </AuthProvider>
          </LocaleProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
