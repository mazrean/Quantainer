import adapter from '@sveltejs/adapter-auto';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess(),

	kit: {
		adapter: adapter(),
		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte',
		vite: {
			server: {
				port: 3000,
				proxy: {
					'/api': {
						target: 'http://localhost:3000',
						changeOrigin: true,
					}
				}
			}
		}
	}
};

export default config;
