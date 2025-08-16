export type AppConfig = {
    port: number;
    mongoUri: string;
    corsOrigin: string | RegExp | (string | RegExp)[];
};
export declare function getConfig(): AppConfig;
//# sourceMappingURL=env.d.ts.map