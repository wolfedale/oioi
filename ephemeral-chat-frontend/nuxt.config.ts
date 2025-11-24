// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  devServer: {
    port: 3001
  },

  modules: [
    '@nuxtjs/tailwindcss'
  ],

  app: {
    head: {
      title: 'Ephemeral Chat',
      meta: [
        { name: 'description', content: 'A beautiful ephemeral chat experience for 2026' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' }
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap'
        }
      ]
    }
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || 'http://localhost:3000'
    }
  }
})
