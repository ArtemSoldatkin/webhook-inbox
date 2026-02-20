import * as z from 'zod';

const envSchema = z.object({
	VITE_ENV: z.enum(['dev', 'uat', 'prod']).default('dev')
});

export default envSchema.parse(import.meta.env);
