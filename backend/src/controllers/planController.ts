import type { Request, Response } from 'express';
import { WeeklyPlanModel } from '../models/Plan';
import { UserModel } from '../models/User';
import { MealModel } from '../models/Meal';
import { z } from 'zod';

function getMonday(date = new Date()): Date {
	const d = new Date(date);
	const day = d.getDay();
	const diff = (day === 0 ? -6 : 1) - day; // adjust when day is sunday
	d.setDate(d.getDate() + diff);
	d.setHours(0, 0, 0, 0);
	return d;
}

const generateSchema = z.object({ authProviderId: z.string() });

export async function generateWeeklyPlan(req: Request, res: Response) {
	const parse = generateSchema.safeParse(req.body);
	if (!parse.success) return res.status(400).json({ error: 'Invalid payload', details: parse.error.flatten() });
	const { authProviderId } = parse.data;
	try {
		const user = await UserModel.findOne({ authProviderId });
		if (!user) return res.status(404).json({ error: 'User not found' });

		const weekStartDate = getMonday();
		const likedMeals = user.likedMeals.length
			? await MealModel.find({ _id: { $in: user.likedMeals } })
			: await MealModel.aggregate([{ $sample: { size: 21 } }]);

		const days = Array.from({ length: 7 }).map((_, idx) => {
			const date = new Date(weekStartDate);
			date.setDate(weekStartDate.getDate() + idx);
			return { date, meals: [] as any[] };
		});

		const mealTypes: Array<'breakfast' | 'lunch' | 'dinner'> = ['breakfast', 'lunch', 'dinner'];
		let i = 0;
		for (const day of days) {
			for (const mealType of mealTypes) {
				const meal = likedMeals[i % likedMeals.length];
				day.meals.push({ meal: meal._id, mealType });
				i += 1;
			}
		}

		const plan = await WeeklyPlanModel.findOneAndUpdate(
			{ user: user._id, weekStartDate },
			{ $set: { days } },
			{ new: true, upsert: true }
		);
		return res.json(plan);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to generate plan' });
	}
}

export async function getWeeklyPlan(req: Request, res: Response) {
	const authProviderId = String(req.query.authProviderId || '');
	if (!authProviderId) return res.status(400).json({ error: 'authProviderId is required' });
	try {
		const user = await UserModel.findOne({ authProviderId });
		if (!user) return res.status(404).json({ error: 'User not found' });
		const weekStartDate = getMonday();
		const plan = await WeeklyPlanModel.findOne({ user: user._id, weekStartDate }).populate('days.meals.meal');
		if (!plan) return res.status(404).json({ error: 'Plan not found' });
		return res.json(plan);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to get plan' });
	}
}


