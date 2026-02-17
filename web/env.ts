import { config } from 'dotenv';

config({ path: '../.env' });

import * as z from 'zod';

const envSchema = z.object({
	API_PORT: z.string().default('3000')
});

export default envSchema.parse(process.env);
