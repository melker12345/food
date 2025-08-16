import type { Request, Response } from 'express';
import { MealModel } from '../models/Meal';
import { UserModel } from '../models/User';
import { z } from 'zod';

const mealCreateSchema = z.object({
	name: z.string(),
	imageUrl: z.string().url(),
	ingredients: z.array(z.object({ name: z.string(), quantity: z.string() })),
	instructions: z.string(),
	nutrition: z.object({ calories: z.number(), protein: z.number(), carbs: z.number(), fat: z.number() }),
	dietaryTags: z.array(z.string()).default([]),
});

export async function createMeal(req: Request, res: Response) {
	const result = mealCreateSchema.safeParse(req.body);
	if (!result.success) {
		return res.status(400).json({ error: 'Invalid meal payload', details: result.error.flatten() });
	}
	try {
		const meal = await MealModel.create(result.data);
		return res.status(201).json(meal);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to create meal' });
	}
}

export async function getRandomMeal(req: Request, res: Response) {
	try {
		const { dietary } = req.query;
		const tags = typeof dietary === 'string' && dietary.length > 0 ? dietary.split(',') : [];
		const match = tags.length ? { dietaryTags: { $all: tags } } : {};
		const [meal] = await MealModel.aggregate([{ $match: match }, { $sample: { size: 1 } }]);
		if (!meal) return res.status(404).json({ error: 'No meals found' });
		return res.json(meal);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to fetch random meal' });
	}
}

const likeSchema = z.object({ authProviderId: z.string(), mealId: z.string() });

export async function likeMeal(req: Request, res: Response) {
	const parse = likeSchema.safeParse(req.body);
	if (!parse.success) return res.status(400).json({ error: 'Invalid payload', details: parse.error.flatten() });
	const { authProviderId, mealId } = parse.data;
	try {
		const user = await UserModel.findOneAndUpdate(
			{ authProviderId },
			{ $addToSet: { likedMeals: mealId } },
			{ new: true }
		);
		return res.json(user);
	} catch (error) {
		return res.status(500).json({ error: 'Failed to like meal' });
	}
}


