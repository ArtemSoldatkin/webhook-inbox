import * as z from 'zod';

/** Validates the public frontend environment variables. */
const envSchema = z.object({
	VITE_ENV: z.enum(['dev', 'uat', 'prod']).default('dev'),
	VITE_API_BASE_URL: z.string().default('http://localhost:3000/api/v1')
});

/** Parsed frontend environment values. */
export default envSchema.parse(import.meta.env);
