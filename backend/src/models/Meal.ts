import { Schema, model, Types } from 'mongoose';

export interface NutritionalInfo {
	calories: number;
	protein: number; // grams
	carbs: number; // grams
	fat: number; // grams
}

export interface IngredientItem {
	name: string;
	quantity: string; // human-readable, e.g., "2 cups", "300 g"
}

export interface MealDocument {
	_id: Types.ObjectId;
	name: string;
	imageUrl: string;
	ingredients: IngredientItem[];
	instructions: string;
	nutrition: NutritionalInfo;
	dietaryTags: string[]; // e.g., vegetarian, vegan, gluten-free
	createdAt: Date;
	updatedAt: Date;
}

const ingredientSchema = new Schema<IngredientItem>({
	name: { type: String, required: true },
	quantity: { type: String, required: true },
});

const nutritionSchema = new Schema<NutritionalInfo>({
	calories: { type: Number, required: true },
	protein: { type: Number, required: true },
	carbs: { type: Number, required: true },
	fat: { type: Number, required: true },
});

const mealSchema = new Schema<MealDocument>(
	{
		name: { type: String, required: true, index: true },
		imageUrl: { type: String, required: true },
		ingredients: { type: [ingredientSchema], default: [] },
		instructions: { type: String, required: true },
		nutrition: { type: nutritionSchema, required: true },
		dietaryTags: { type: [String], default: [] },
	},
	{ timestamps: true }
);

export const MealModel = model<MealDocument>('Meal', mealSchema);


