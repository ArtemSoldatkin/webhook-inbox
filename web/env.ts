import { config } from 'dotenv';

config({ path: '../.env' });

import * as z from 'zod';

const envSchema = z.object({
	ENV: z.enum(['dev', 'uat', 'prod']).default('dev'),
	API_PROTOCOL: z.enum(['http', 'https']).default('http'),
	API_HOST: z.string().default('localhost'),
	API_PORT: z.string().default('3000')
});

export default envSchema.parse(process.env);
