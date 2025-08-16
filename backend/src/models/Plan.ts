import { Schema, model, Types } from 'mongoose';

export interface DailyMealEntry {
	meal: Types.ObjectId; // ref Meal
	mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
}

export interface WeeklyPlanDocument {
	_id: Types.ObjectId;
	user: Types.ObjectId; // ref User
	weekStartDate: Date; // Monday of the week
	days: Array<{
		date: Date;
		meals: DailyMealEntry[];
	}>;
	createdAt: Date;
	updatedAt: Date;
}

const dailyMealSchema = new Schema<DailyMealEntry>({
	meal: { type: Schema.Types.ObjectId, ref: 'Meal', required: true },
	mealType: { type: String, enum: ['breakfast', 'lunch', 'dinner', 'snack'], required: true },
});

const weeklyPlanSchema = new Schema<WeeklyPlanDocument>(
	{
		user: { type: Schema.Types.ObjectId, ref: 'User', required: true, index: true },
		weekStartDate: { type: Date, required: true, index: true },
		days: [
			new Schema(
				{
					date: { type: Date, required: true },
					meals: { type: [dailyMealSchema], default: [] },
				},
				{ _id: false }
			),
		],
	},
	{ timestamps: true }
);

weeklyPlanSchema.index({ user: 1, weekStartDate: 1 }, { unique: true });

export const WeeklyPlanModel = model<WeeklyPlanDocument>('WeeklyPlan', weeklyPlanSchema);


