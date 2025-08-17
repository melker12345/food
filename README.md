# Food Planning App

A comprehensive meal planning application that helps users discover, plan, and organize their meals with intelligent recommendations and automated shopping list generation.

## Features

### 🍽️ Core Features
- **Swipe-Based Meal Discovery**: Tinder-like interface for discovering new meals
- **Personalized Recommendations**: AI-powered suggestions based on dietary preferences
- **Weekly Meal Planning**: Single active meal plan automatically populated from liked meals
- **Dynamic Shopping Lists**: Single shopping list that updates automatically when meals change
- **Interactive Planning**: Visual weekly grid with drag-and-drop meal management
- **Smart Auto-Population**: Meals intelligently distributed by type (breakfast, lunch, dinner)
- **User Profiles & Preferences**: Customizable dietary restrictions and preferences

### 🎯 Advanced Features
- **Social Features**: Like, review, and share favorite meals
- **Nutrition Tracking**: Detailed nutritional information for all meals
- **Trending Meals**: Discover popular meals in the community
- **Dark Mode**: Complete dark/light theme support
- **Responsive Design**: Optimized for desktop and mobile devices

### 🔧 Technical Features
- **Real-time Updates**: Live synchronization across devices
- **Offline Support**: Core features work without internet
- **Secure Authentication**: JWT-based user authentication
- **RESTful API**: Clean, documented API endpoints
- **Database Optimization**: Efficient PostgreSQL database design

## Tech Stack

### Backend
- **Go 1.21**: High-performance backend API
- **Gin Framework**: Fast HTTP web framework
- **PostgreSQL**: Robust relational database
- **GORM**: Object-relational mapping
- **JWT Authentication**: Secure token-based auth

### Frontend
- **React 18**: Modern frontend framework
- **TypeScript**: Type-safe development
- **Tailwind CSS**: Utility-first styling
- **Redux Toolkit**: State management
- **React Router**: Client-side routing
- **Framer Motion**: Smooth animations

### Infrastructure
- **Docker**: Containerized deployment
- **Docker Compose**: Multi-service orchestration
- **PostgreSQL**: Production-ready database
- **Nginx**: (Production) Reverse proxy

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Option 1: Docker (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd food

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### Option 2: Local Development
```bash
# Start development environment
./scripts/dev-up.sh

# Stop development environment
./scripts/dev-down.sh
```

### Accessing the Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432

## API Documentation

### Authentication Endpoints
```
POST /api/v1/register    - User registration
POST /api/v1/login       - User login
GET  /api/v1/profile     - Get user profile
PUT  /api/v1/profile     - Update user profile
```

### Meal Endpoints
```
GET    /api/v1/meals                 - Get all meals
GET    /api/v1/meals/:id             - Get specific meal
GET    /api/v1/meals/personalized    - Get personalized recommendations
GET    /api/v1/meals/trending        - Get trending meals
POST   /api/v1/meals/:id/like        - Like a meal
POST   /api/v1/meals/:id/dislike     - Dislike a meal
```

### Meal Planning Endpoints
```
# Current Week Meal Plan (Primary Workflow)
GET  /api/v1/current-meal-plan                    - Get user's current weekly meal plan
POST /api/v1/current-meal-plan/populate-from-liked - Auto-populate from liked meals ✨
PUT  /api/v1/current-meal-plan/meals              - Update specific meal in plan
PUT  /api/v1/shopping-items/:item_id              - Toggle shopping item purchased status

# Legacy Multiple Plans (Alternative)
GET    /api/v1/meal-plans                - Get user's meal plans
POST   /api/v1/meal-plans                - Create new meal plan
POST   /api/v1/meal-plans/auto-generate  - Auto-generate meal plan from liked meals
GET    /api/v1/meal-plans/:id            - Get specific meal plan
PUT    /api/v1/meal-plans/:id            - Update meal plan
DELETE /api/v1/meal-plans/:id            - Delete meal plan
```

### Shopping List Endpoints
```
POST /api/v1/meal-plans/:id/shopping-list  - Generate shopping list
GET  /api/v1/shopping-lists                - Get all shopping lists
PUT  /api/v1/shopping-list-items/:id       - Update shopping list item
```

## Database Schema

### Key Tables
- **users**: User accounts and preferences
- **meals**: Recipe information and metadata
- **ingredients**: Food items and nutritional data
- **meal_plans**: Weekly meal planning data
- **shopping_lists**: Generated grocery lists
- **user_meal_interactions**: Likes/dislikes tracking

## Development

### Project Structure
```
food/
├── backend/           # Go API server
├── frontend/          # React application
├── database/          # Database initialization
├── scripts/           # Development scripts
├── docker-compose.yml # Container orchestration
└── README.md         # This file
```

### Adding New Features

1. **Backend**: Add handlers in `backend/handlers/`
2. **Frontend**: Add components in `frontend/src/components/`
3. **Database**: Update models in `backend/models/`
4. **API**: Update Redux slices in `frontend/src/store/slices/`

### Environment Variables

Backend:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=food_app
JWT_SECRET=your-secret-key
```

## Deployment

### Production Deployment
```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d

# View production logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Environment Setup
1. Update environment variables in docker-compose files
2. Configure SSL certificates for HTTPS
3. Set up domain name and DNS
4. Configure backup strategy for database

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Integration Tests
```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
npm run test:integration
```

## Troubleshooting

### Common Issues

**Database Connection Failed**
- Check if PostgreSQL container is running
- Verify database credentials
- Ensure port 5432 is not blocked

**Frontend Not Loading**
- Check if Node.js dependencies are installed
- Verify backend API is accessible
- Check browser console for errors

**Build Failures**
- Clear Docker cache: `docker system prune -a`
- Rebuild containers: `docker-compose build --no-cache`

### Getting Help
- Check the [Issues](link-to-issues) page
- Review application logs in `logs/` directory
- Contact the development team

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Recipe data provided by community contributors
- UI/UX inspired by modern food applications
- Icons from Heroicons
- Images from Unsplash