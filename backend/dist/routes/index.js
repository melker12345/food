"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const userController_1 = require("../controllers/userController");
const mealController_1 = require("../controllers/mealController");
const planController_1 = require("../controllers/planController");
const router = (0, express_1.Router)();
// Health
router.get('/health', (_req, res) => res.json({ status: 'ok' }));
// Users
router.post('/users/upsert', userController_1.upsertUser);
router.get('/users/me', userController_1.getMe);
// Meals
router.post('/meals', mealController_1.createMeal);
router.get('/meals/random', mealController_1.getRandomMeal);
router.post('/meals/like', mealController_1.likeMeal);
// Plans
router.post('/plans/generate', planController_1.generateWeeklyPlan);
router.get('/plans/weekly', planController_1.getWeeklyPlan);
exports.default = router;
//# sourceMappingURL=index.js.map