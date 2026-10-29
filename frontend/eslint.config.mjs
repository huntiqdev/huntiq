import nextConfig from 'eslint-config-next';

export default [
  ...nextConfig,
  {
    rules: {
      'react/no-unescaped-entities': 'warn',
      'react-hooks/set-state-in-effect': 'warn',
      'react-hooks/immutability': 'warn',
      'react-hooks/preserve-manual-memoization': 'warn',
      'react-hooks/purity': 'warn',
      'react-hooks/refs': 'warn',
      'react-hooks/static-components': 'warn',
      'react-hooks/exhaustive-deps': 'warn',
      '@next/next/no-location-assign-relative-destination': 'warn',
    },
  },
];
