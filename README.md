# The Block

## How to Run

### Automatic setup

- Install [mise](https://mise.jdx.dev/)
- Run `mise trust` and `mise install` to install dev tooling
- Run `npm run prepare` to install dependencies
- Run `npm run dev` to start dev server with backend and frontend services

### Manual setup

- Install [Node.js 24 or newer](https://nodejs.org/en/download/) and [go 1.26 or newer](https://go.dev/dl/)
- Run `npm run prepare` to install dependencies
- Run `npm run dev` to start dev server with backend and frontend services

## Time Spent

I spent approximately 8 hours on this project.

## Assumptions and Scope

This program is a mock of what the experience might be like for a user who is logged in. Since there is no authentication, all endpoints are public. This simulates a search page, vehicle info page, and views for placing a bid or buying a vehicle with a history on bids/purchases. This does not handle inventory management, but it does simulate bid increasing and purchase logic in the form of mock UIs and simple endpoints.

## Stack

- **Frontend:** Vue, shadcn-vue components, Pinia stores, openapi-typescript and openapi-fetch
- **Backend:** Go, Fuego Web Framework, GORM ORM
- **Database:** SQLite

## What I Built

I built a simple search and purchase experience that serves as a full stack prototype of a vehicle auction platform. This demonstrates the user flow for searching for vehicles, inspecting vehicle details, and making bids and purchases. The user can filter for details of vehicles, as well as sort by auction details like price and end time.

The goal was to make this project type safe across the whole stack. Since this was not written using a fullstack framework such as Next.js or Ruby on Rails, I used code generation from OpenAPI documents to generate types and interfaces. In the backend, Go generates the OpenAPI spec using the Fuego web framework. The frontend then uses openapi-typescript and openapi-fetch to generate types from the OpenAPI spec and have type safe fetch requests to the API backend. The backend also uses GORM, an ORM, to keep our SQL tables and statements typesafe.

The UI is built off of shadcn-vue primatives and Tailwind. Mobile and Desktop views were taken into consideration when building.

## Notable Decisions

I let the backend send all vehicles that matched filters to the frontend. I did not include pagination or virtual list scrolling, which would help improve client performance as the number of vehicles increase.

I did not make my filters dynamically show based on current filters. Some filters would have to be dynamic, while other would have to be static due to the nature of what customers might want to search for.

Vehicles are not removed from user view if the auction ends or a "purchase" is made. All vehicles are visible to users when no filters are applied.

## Testing

I added simple unit tests for bidding logic, since it wrote data to the DB and it could easily be verified in a test environment. While I could test search query logic and UI, I figured it was best to skip them due to time constraints and the ability to test them manually as I was working on features.

## What I'd Do With More Time

I would add more test to make sure that features and UI work, especially in more complex situations. I would spend more time refining the UI to make it more impactful, as while I do think the UI is clean it can feel generic at times.

Documentation of backend API endpoints could also be improved, as we could embed them in the OpenAPI spec using Fuego. This would help consumers understand the API better and what the expectations are besides input and output schemas.

Migrating to a fullstack framework like Nuxt could give some nice to have features such as SSR, a robust module/plugin system, and using Backend-for-frontend archecture to slim down API requests from the backend if we aren't the ones in control of the API we are consuming. These might introduce latency issues, which is not ideal if you are trying to have bidding be instant, so there are tradeoffs to consider.
