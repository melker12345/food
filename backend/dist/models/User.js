"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserModel = void 0;
const mongoose_1 = require("mongoose");
const userSchema = new mongoose_1.Schema({
    authProviderId: { type: String, required: true, index: true, unique: true },
    name: { type: String, required: true },
    email: { type: String, required: true, index: true, unique: true },
    dietaryPreferences: { type: [String], default: [] },
    healthGoals: { type: [String], default: [] },
    likedMeals: [{ type: mongoose_1.Schema.Types.ObjectId, ref: 'Meal', default: [] }],
}, { timestamps: true });
exports.UserModel = (0, mongoose_1.model)('User', userSchema);
//# sourceMappingURL=User.js.map