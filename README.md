Food App Monorepo

Stack
- Frontend: Expo (React Native, TypeScript) in `frontend`
- Backend: Express + TypeScript + Mongoose in `backend`
- Database: MongoDB (Docker Compose)

Quick Start
1) Start Docker daemon (Linux):
```bash
sudo systemctl start docker
```

2) Start MongoDB:
```bash
docker compose up -d
```

3) Seed sample meals:
```bash
cd backend
npm run seed
```

4) Start backend API (http://localhost:4000):
```bash
npm run dev
```

5) Start frontend app:
```bash
cd ../frontend
npm run web
```

Note: On Android/iOS devices, replace `API_BASE` in `frontend/App.tsx` with your machine LAN IP (e.g., `http://192.168.1.10:4000`).

API
- `GET /api/health`
- `POST /api/users/upsert` { authProviderId, name, email, dietaryPreferences?, healthGoals? }
- `GET /api/users/me?authProviderId=...`
- `POST /api/meals` create meal
- `GET /api/meals/random` optional `?dietary=vegan,gluten-free`
- `POST /api/meals/like` { authProviderId, mealId }
- `POST /api/plans/generate` { authProviderId }
- `GET /api/plans/weekly?authProviderId=...`

Env Vars (backend)
Defaults are reasonable; create `backend/.env` if needed:
```
PORT=4000
MONGODB_URI=mongodb://localhost:27017/food
CORS_ORIGIN=*
```


