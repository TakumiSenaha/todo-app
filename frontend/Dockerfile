# Builder stage
FROM node:18-alpine AS builder

# Install pnpm
RUN npm install -g pnpm

# Set working directory
WORKDIR /app

# Copy dependency files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy all source code
COPY . .

# Build the application
RUN pnpm run build

# Final stage
FROM node:18-alpine AS final

# Install pnpm
RUN npm install -g pnpm

# Set working directory
WORKDIR /app

# Copy built application from builder stage
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next

# Copy dependency files
COPY package.json pnpm-lock.yaml ./

# Install production dependencies only
RUN pnpm install --prod

# Expose port
EXPOSE 3000

# Run the application
CMD ["pnpm", "start"]
