"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MealModel = void 0;
const mongoose_1 = require("mongoose");
const ingredientSchema = new mongoose_1.Schema({
    name: { type: String, required: true },
    quantity: { type: String, required: true },
});
const nutritionSchema = new mongoose_1.Schema({
    calories: { type: Number, required: true },
    protein: { type: Number, required: true },
    carbs: { type: Number, required: true },
    fat: { type: Number, required: true },
});
const mealSchema = new mongoose_1.Schema({
    name: { type: String, required: true, index: true },
    imageUrl: { type: String, required: true },
    ingredients: { type: [ingredientSchema], default: [] },
    instructions: { type: String, required: true },
    nutrition: { type: nutritionSchema, required: true },
    dietaryTags: { type: [String], default: [] },
}, { timestamps: true });
exports.MealModel = (0, mongoose_1.model)('Meal', mealSchema);
//# sourceMappingURL=Meal.js.map