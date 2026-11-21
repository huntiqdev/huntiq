import { MetadataRoute } from 'next';

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: [
          '/admin/',
          '/executive-console/',
          '/platform-console/',
          '/security-console/',
          '/operations-console/',
          '/moderation-console/',
          '/support-console/',
          '/staff/',
          '/api/',
        ],
      },
    ],
    sitemap: 'https://huntiq.com/sitemap.xml',
  };
}
