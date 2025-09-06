FROM node:18-alpine

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

# Expose port
EXPOSE 3000

# Run in development mode for hot reloading
CMD ["pnpm", "dev"]
