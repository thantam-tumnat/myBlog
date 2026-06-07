# Use a base image with Node.js pre-installed
FROM node:20-alpine AS builder

# Set the working directory in the container
WORKDIR /app/react-app

# Install dependencies
COPY package.json .
RUN npm install

# Copy the project files into the container
COPY . .

# Build the Vite project
RUN npm run build

# Multi
FROM node:20-alpine AS production
WORKDIR /app/react-app

# Copy folder from builder
COPY --from=builder /app/react-app/dist/ /app/react-app/dist/

# Expose the port as an environment variable (default: 8003)
ARG APP_PORT=8003
ENV APP_PORT=${APP_PORT}

EXPOSE ${APP_PORT}

COPY package.json .
COPY package-lock.json .
COPY vite.config.js .

RUN npm install vite

# Define the command to run your app
CMD ["npm", "run", "preview"]
