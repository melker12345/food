"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.WeeklyPlanModel = void 0;
const mongoose_1 = require("mongoose");
const dailyMealSchema = new mongoose_1.Schema({
    meal: { type: mongoose_1.Schema.Types.ObjectId, ref: 'Meal', required: true },
    mealType: { type: String, enum: ['breakfast', 'lunch', 'dinner', 'snack'], required: true },
});
const weeklyPlanSchema = new mongoose_1.Schema({
    user: { type: mongoose_1.Schema.Types.ObjectId, ref: 'User', required: true, index: true },
    weekStartDate: { type: Date, required: true, index: true },
    days: [
        new mongoose_1.Schema({
            date: { type: Date, required: true },
            meals: { type: [dailyMealSchema], default: [] },
        }, { _id: false }),
    ],
}, { timestamps: true });
weeklyPlanSchema.index({ user: 1, weekStartDate: 1 }, { unique: true });
exports.WeeklyPlanModel = (0, mongoose_1.model)('WeeklyPlan', weeklyPlanSchema);
//# sourceMappingURL=Plan.js.map