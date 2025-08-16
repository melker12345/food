"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createMeal = createMeal;
exports.getRandomMeal = getRandomMeal;
exports.likeMeal = likeMeal;
const Meal_1 = require("../models/Meal");
const User_1 = require("../models/User");
const zod_1 = require("zod");
const mealCreateSchema = zod_1.z.object({
    name: zod_1.z.string(),
    imageUrl: zod_1.z.string().url(),
    ingredients: zod_1.z.array(zod_1.z.object({ name: zod_1.z.string(), quantity: zod_1.z.string() })),
    instructions: zod_1.z.string(),
    nutrition: zod_1.z.object({ calories: zod_1.z.number(), protein: zod_1.z.number(), carbs: zod_1.z.number(), fat: zod_1.z.number() }),
    dietaryTags: zod_1.z.array(zod_1.z.string()).default([]),
});
async function createMeal(req, res) {
    const result = mealCreateSchema.safeParse(req.body);
    if (!result.success) {
        return res.status(400).json({ error: 'Invalid meal payload', details: result.error.flatten() });
    }
    try {
        const meal = await Meal_1.MealModel.create(result.data);
        return res.status(201).json(meal);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to create meal' });
    }
}
async function getRandomMeal(req, res) {
    try {
        const { dietary } = req.query;
        const tags = typeof dietary === 'string' && dietary.length > 0 ? dietary.split(',') : [];
        const match = tags.length ? { dietaryTags: { $all: tags } } : {};
        const [meal] = await Meal_1.MealModel.aggregate([{ $match: match }, { $sample: { size: 1 } }]);
        if (!meal)
            return res.status(404).json({ error: 'No meals found' });
        return res.json(meal);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to fetch random meal' });
    }
}
const likeSchema = zod_1.z.object({ authProviderId: zod_1.z.string(), mealId: zod_1.z.string() });
async function likeMeal(req, res) {
    const parse = likeSchema.safeParse(req.body);
    if (!parse.success)
        return res.status(400).json({ error: 'Invalid payload', details: parse.error.flatten() });
    const { authProviderId, mealId } = parse.data;
    try {
        const user = await User_1.UserModel.findOneAndUpdate({ authProviderId }, { $addToSet: { likedMeals: mealId } }, { new: true });
        return res.json(user);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to like meal' });
    }
}
//# sourceMappingURL=mealController.js.map