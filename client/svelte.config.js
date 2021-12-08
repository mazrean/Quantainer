//import adapter from '@sveltejs/adapter-auto';
import preprocess from 'svelte-preprocess';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess(),

	kit: {
		adapter: adapter({
			fallback: 'index.html',
		}),
		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte',
		prerender: {
			enabled: false,
		},
		ssr: false,
		vite: {
			server: {
				port: 3000,
				proxy: {
					'/api/v1': {
						target: 'http://localhost:3000',
						changeOrigin: true,
					}
				}
			}
		}
	}
};

export default config;
