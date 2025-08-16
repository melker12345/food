"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getConfig = getConfig;
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
function getConfig() {
    const port = Number(process.env.PORT || 4000);
    const mongoUri = process.env.MONGODB_URI || 'mongodb://localhost:27017/food';
    const corsOriginRaw = process.env.CORS_ORIGIN || '*';
    let corsOrigin = '*';
    if (corsOriginRaw.includes(',')) {
        corsOrigin = corsOriginRaw.split(',').map((s) => s.trim());
    }
    else {
        corsOrigin = corsOriginRaw.trim();
    }
    return { port, mongoUri, corsOrigin };
}
//# sourceMappingURL=env.js.map