import { createApp, h } from 'vue';
import { createInertiaApp } from '@inertiajs/vue3';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import '../css/app.css';
import axios from 'axios';
import { initializeTheme } from './composables/useAppearance';


// Configure axios BEFORE Inertia uses it
axios.defaults.withCredentials = true;
axios.defaults.xsrfCookieName = 'XSRF-TOKEN';
axios.defaults.xsrfHeaderName = 'X-XSRF-TOKEN';
axios.defaults.withXSRFToken = true

const appName = import.meta.env.VITE_APP_NAME || 'Sheta App';

createInertiaApp({
  title: (title) => (title ? `${title} - ${appName}` : appName),
  // @ts-expect-error
  resolve: (name) => resolvePageComponent(`./Pages/${name}.vue`, import.meta.glob('./Pages/**/*.vue')),
  setup({ el, App, props, plugin }) {
    return createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el);
  },
  progress: {
    color: '#4B5563',
  },
});

initializeTheme();
