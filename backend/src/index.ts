import express from 'express';
import cors from 'cors';
import mongoose from 'mongoose';
import router from './routes';
import { getConfig } from './config/env';

async function startServer() {
	const app = express();
	const config = getConfig();

	app.use(cors({ origin: config.corsOrigin }));
	app.use(express.json({ limit: '1mb' }));
	app.use('/api', router);

	mongoose.set('strictQuery', true);
	await mongoose.connect(config.mongoUri);

	app.listen(config.port, () => {
		console.log(`API listening on http://localhost:${config.port}`);
	});
}

startServer().catch((error) => {
	console.error('Failed to start server', error);
	process.exit(1);
});


