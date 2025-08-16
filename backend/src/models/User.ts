import { Schema, model, Types } from 'mongoose';

export type DietaryPreference =
	| 'vegetarian'
	| 'vegan'
	| 'gluten-free'
	| 'keto'
	| 'paleo'
	| 'halal'
	| 'kosher'
	| 'none';

export interface UserDocument {
	_id: Types.ObjectId;
	authProviderId: string; // e.g., Firebase UID or Auth0 sub
	name: string;
	email: string;
	dietaryPreferences: DietaryPreference[];
	healthGoals: string[];
	likedMeals: Types.ObjectId[]; // references Meal _id
	createdAt: Date;
	updatedAt: Date;
}

const userSchema = new Schema<UserDocument>(
	{
		authProviderId: { type: String, required: true, index: true, unique: true },
		name: { type: String, required: true },
		email: { type: String, required: true, index: true, unique: true },
		dietaryPreferences: { type: [String], default: [] },
		healthGoals: { type: [String], default: [] },
		likedMeals: [{ type: Schema.Types.ObjectId, ref: 'Meal', default: [] }],
	},
	{ timestamps: true }
);

export const UserModel = model<UserDocument>('User', userSchema);


