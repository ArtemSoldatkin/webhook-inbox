import { config } from 'dotenv';

import * as z from 'zod';

config({ path: '../.env' });

const EnvSchema = z.object({
	API_PORT: z.string().default('3000')
});

export default EnvSchema.parse(process.env);
