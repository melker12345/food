"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateWeeklyPlan = generateWeeklyPlan;
exports.getWeeklyPlan = getWeeklyPlan;
const Plan_1 = require("../models/Plan");
const User_1 = require("../models/User");
const Meal_1 = require("../models/Meal");
const zod_1 = require("zod");
function getMonday(date = new Date()) {
    const d = new Date(date);
    const day = d.getDay();
    const diff = (day === 0 ? -6 : 1) - day; // adjust when day is sunday
    d.setDate(d.getDate() + diff);
    d.setHours(0, 0, 0, 0);
    return d;
}
const generateSchema = zod_1.z.object({ authProviderId: zod_1.z.string() });
async function generateWeeklyPlan(req, res) {
    const parse = generateSchema.safeParse(req.body);
    if (!parse.success)
        return res.status(400).json({ error: 'Invalid payload', details: parse.error.flatten() });
    const { authProviderId } = parse.data;
    try {
        const user = await User_1.UserModel.findOne({ authProviderId });
        if (!user)
            return res.status(404).json({ error: 'User not found' });
        const weekStartDate = getMonday();
        const likedMeals = user.likedMeals.length
            ? await Meal_1.MealModel.find({ _id: { $in: user.likedMeals } })
            : await Meal_1.MealModel.aggregate([{ $sample: { size: 21 } }]);
        const days = Array.from({ length: 7 }).map((_, idx) => {
            const date = new Date(weekStartDate);
            date.setDate(weekStartDate.getDate() + idx);
            return { date, meals: [] };
        });
        const mealTypes = ['breakfast', 'lunch', 'dinner'];
        let i = 0;
        for (const day of days) {
            for (const mealType of mealTypes) {
                const meal = likedMeals[i % likedMeals.length];
                day.meals.push({ meal: meal._id, mealType });
                i += 1;
            }
        }
        const plan = await Plan_1.WeeklyPlanModel.findOneAndUpdate({ user: user._id, weekStartDate }, { $set: { days } }, { new: true, upsert: true });
        return res.json(plan);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to generate plan' });
    }
}
async function getWeeklyPlan(req, res) {
    const authProviderId = String(req.query.authProviderId || '');
    if (!authProviderId)
        return res.status(400).json({ error: 'authProviderId is required' });
    try {
        const user = await User_1.UserModel.findOne({ authProviderId });
        if (!user)
            return res.status(404).json({ error: 'User not found' });
        const weekStartDate = getMonday();
        const plan = await Plan_1.WeeklyPlanModel.findOne({ user: user._id, weekStartDate }).populate('days.meals.meal');
        if (!plan)
            return res.status(404).json({ error: 'Plan not found' });
        return res.json(plan);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to get plan' });
    }
}
//# sourceMappingURL=planController.js.map