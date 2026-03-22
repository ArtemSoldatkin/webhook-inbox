import * as z from 'zod';

/** Validates the public frontend environment variables. */
const envSchema = z.object({
	VITE_ENV: z.enum(['dev', 'uat', 'prod']).default('dev')
});

/** Parsed frontend environment values. */
export default envSchema.parse(import.meta.env);
