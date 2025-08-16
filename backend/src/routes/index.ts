import { Router } from 'express';
import { upsertUser, getMe } from '../controllers/userController';
import { createMeal, getRandomMeal, likeMeal } from '../controllers/mealController';
import { generateWeeklyPlan, getWeeklyPlan } from '../controllers/planController';

const router = Router();

// Health
router.get('/health', (_req, res) => res.json({ status: 'ok' }));

// Users
router.post('/users/upsert', upsertUser);
router.get('/users/me', getMe);

// Meals
router.post('/meals', createMeal);
router.get('/meals/random', getRandomMeal);
router.post('/meals/like', likeMeal);

// Plans
router.post('/plans/generate', generateWeeklyPlan);
router.get('/plans/weekly', getWeeklyPlan);

export default router;


