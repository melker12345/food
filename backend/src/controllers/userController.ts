import type { Request, Response } from 'express';
import { UserModel } from '../models/User';
import { z } from 'zod';

const upsertUserSchema = z.object({
	authProviderId: z.string().min(1),
	name: z.string().min(1),
	email: z.string().email(),
	dietaryPreferences: z.array(z.string()).optional(),
	healthGoals: z.array(z.string()).optional(),
});

export async function upsertUser(req: Request, res: Response) {
	const parse = upsertUserSchema.safeParse(req.body);
	if (!parse.success) {
		return res.status(400).json({ error: 'Invalid request', details: parse.error.flatten() });
	}
	const { authProviderId, name, email, dietaryPreferences = [], healthGoals = [] } = parse.data;
  try {
		const user = await UserModel.findOneAndUpdate(
			{ authProviderId },
			{ $set: { name, email, dietaryPreferences, healthGoals } },
			{ new: true, upsert: true }
		);
		return res.json(user);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to upsert user' });
	}
}

export async function getMe(req: Request, res: Response) {
	const authProviderId = String(req.query.authProviderId || '');
	if (!authProviderId) {
		return res.status(400).json({ error: 'authProviderId is required' });
	}
	const user = await UserModel.findOne({ authProviderId });
	if (!user) return res.status(404).json({ error: 'User not found' });
	return res.json(user);
}


