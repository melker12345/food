"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
const mongoose_1 = __importDefault(require("mongoose"));
const routes_1 = __importDefault(require("./routes"));
const env_1 = require("./config/env");
async function startServer() {
    const app = (0, express_1.default)();
    const config = (0, env_1.getConfig)();
    app.use((0, cors_1.default)({ origin: config.corsOrigin }));
    app.use(express_1.default.json({ limit: '1mb' }));
    app.use('/api', routes_1.default);
    mongoose_1.default.set('strictQuery', true);
    await mongoose_1.default.connect(config.mongoUri);
    app.listen(config.port, () => {
        console.log(`API listening on http://localhost:${config.port}`);
    });
}
startServer().catch((error) => {
    console.error('Failed to start server', error);
    process.exit(1);
});
//# sourceMappingURL=index.js.map