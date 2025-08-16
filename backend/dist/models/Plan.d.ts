import { Types } from 'mongoose';
export interface DailyMealEntry {
    meal: Types.ObjectId;
    mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
}
export interface WeeklyPlanDocument {
    _id: Types.ObjectId;
    user: Types.ObjectId;
    weekStartDate: Date;
    days: Array<{
        date: Date;
        meals: DailyMealEntry[];
    }>;
    createdAt: Date;
    updatedAt: Date;
}
export declare const WeeklyPlanModel: import("mongoose").Model<WeeklyPlanDocument, {}, {}, {}, import("mongoose").Document<unknown, {}, WeeklyPlanDocument, {}, {}> & WeeklyPlanDocument & Required<{
    _id: Types.ObjectId;
}> & {
    __v: number;
}, any>;
//# sourceMappingURL=Plan.d.ts.map