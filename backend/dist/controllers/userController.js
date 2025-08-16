"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.upsertUser = upsertUser;
exports.getMe = getMe;
const User_1 = require("../models/User");
const zod_1 = require("zod");
const upsertUserSchema = zod_1.z.object({
    authProviderId: zod_1.z.string().min(1),
    name: zod_1.z.string().min(1),
    email: zod_1.z.string().email(),
    dietaryPreferences: zod_1.z.array(zod_1.z.string()).optional(),
    healthGoals: zod_1.z.array(zod_1.z.string()).optional(),
});
async function upsertUser(req, res) {
    const parse = upsertUserSchema.safeParse(req.body);
    if (!parse.success) {
        return res.status(400).json({ error: 'Invalid request', details: parse.error.flatten() });
    }
    const { authProviderId, name, email, dietaryPreferences = [], healthGoals = [] } = parse.data;
    try {
        const user = await User_1.UserModel.findOneAndUpdate({ authProviderId }, { $set: { name, email, dietaryPreferences, healthGoals } }, { new: true, upsert: true });
        return res.json(user);
    }
    catch (error) {
        return res.status(500).json({ error: 'Failed to upsert user' });
    }
}
async function getMe(req, res) {
    const authProviderId = String(req.query.authProviderId || '');
    if (!authProviderId) {
        return res.status(400).json({ error: 'authProviderId is required' });
    }
    const user = await User_1.UserModel.findOne({ authProviderId });
    if (!user)
        return res.status(404).json({ error: 'User not found' });
    return res.json(user);
}
//# sourceMappingURL=userController.js.map