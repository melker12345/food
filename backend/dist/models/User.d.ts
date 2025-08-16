import { Types } from 'mongoose';
export type DietaryPreference = 'vegetarian' | 'vegan' | 'gluten-free' | 'keto' | 'paleo' | 'halal' | 'kosher' | 'none';
export interface UserDocument {
    _id: Types.ObjectId;
    authProviderId: string;
    name: string;
    email: string;
    dietaryPreferences: DietaryPreference[];
    healthGoals: string[];
    likedMeals: Types.ObjectId[];
    createdAt: Date;
    updatedAt: Date;
}
export declare const UserModel: import("mongoose").Model<UserDocument, {}, {}, {}, import("mongoose").Document<unknown, {}, UserDocument, {}, {}> & UserDocument & Required<{
    _id: Types.ObjectId;
}> & {
    __v: number;
}, any>;
//# sourceMappingURL=User.d.ts.map