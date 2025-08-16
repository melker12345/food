import dotenv from 'dotenv';

dotenv.config();

export type AppConfig = {
	port: number;
	mongoUri: string;
	corsOrigin: string | RegExp | (string | RegExp)[];
};

export function getConfig(): AppConfig {
	const port = Number(process.env.PORT || 4000);
	const mongoUri = process.env.MONGODB_URI || 'mongodb://localhost:27017/food';
	const corsOriginRaw = process.env.CORS_ORIGIN || '*';

	let corsOrigin: AppConfig['corsOrigin'] = '*';
	if (corsOriginRaw.includes(',')) {
		corsOrigin = corsOriginRaw.split(',').map((s) => s.trim());
	} else {
		corsOrigin = corsOriginRaw.trim();
	}

	return { port, mongoUri, corsOrigin };
}


