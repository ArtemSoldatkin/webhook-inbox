import { config } from 'dotenv';

config({ path: '../.env' });

import * as z from 'zod';

const envSchema = z.object({
	ENV: z.enum(['dev', 'uat', 'prod']).default('dev'),
	API_PROTOCOL: z.enum(['http', 'https']).default('http'),
	API_HOST: z.string().default('localhost'),
	API_PORT: z.string().default('3000'),

	UI_PROTOCOL: z.enum(['http', 'https']).default('http'),
	UI_HOST: z.string().default('localhost'),
	UI_PORT: z.string().default('5173'),

	UI_HTTPS_KEY_PATH: z.string().default('./certs/dev.key'),
	UI_HTTPS_CERT_PATH: z.string().default('./certs/dev.crt')
});

export default envSchema.parse(process.env);
