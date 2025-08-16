import { Types } from 'mongoose';
export interface NutritionalInfo {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
}
export interface IngredientItem {
    name: string;
    quantity: string;
}
export interface MealDocument {
    _id: Types.ObjectId;
    name: string;
    imageUrl: string;
    ingredients: IngredientItem[];
    instructions: string;
    nutrition: NutritionalInfo;
    dietaryTags: string[];
    createdAt: Date;
    updatedAt: Date;
}
export declare const MealModel: import("mongoose").Model<MealDocument, {}, {}, {}, import("mongoose").Document<unknown, {}, MealDocument, {}, {}> & MealDocument & Required<{
    _id: Types.ObjectId;
}> & {
    __v: number;
}, any>;
//# sourceMappingURL=Meal.d.ts.map