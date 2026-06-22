import { createSSRApp, h } from 'vue';
import { renderToString } from '@vue/server-renderer';
import { createInertiaApp } from '@inertiajs/vue3';
import createServer from '@inertiajs/vue3/server';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import axios from 'axios';
import { initializeTheme } from './composables/useAppearance';

// Configure axios BEFORE Inertia uses it
axios.defaults.withCredentials = true;
axios.defaults.xsrfCookieName = 'XSRF-TOKEN';
axios.defaults.xsrfHeaderName = 'X-XSRF-TOKEN';
axios.defaults.withXSRFToken = true

const appName = import.meta.env.VITE_APP_NAME || 'Sheta App';

createServer((page) =>
  createInertiaApp({
    title: (title) => (title ? `${title} - ${appName}` : appName),
    page,
    render: renderToString,
    // @ts-expect-error
    resolve: (name) => resolvePageComponent(`./Pages/${name}.vue`, import.meta.glob('./Pages/**/*.vue')),
    setup({ App, props, plugin }) {
      return createSSRApp({ render: () => h(App, props) })
        .use(plugin);
    },
  }),
  { cluster: true },
);

initializeTheme();

